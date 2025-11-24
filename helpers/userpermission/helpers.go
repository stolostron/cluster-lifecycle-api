package userpermission

import (
	"context"
	"encoding/json"
	"fmt"
	clusterviewclientset "github.com/stolostron/cluster-lifecycle-api/client/clusterview/clientset/versioned"
	clusterviewv1alpha1 "github.com/stolostron/cluster-lifecycle-api/clusterview/v1alpha1"
	"hash/fnv"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"slices"

	authzv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/rest"
)

var adminResourceRule = authzv1.ResourceRule{
	Verbs:     []string{"*"},
	APIGroups: []string{"*"},
	Resources: []string{"*"},
}

// internal cache to keep permissionRule
type permissionRule struct {
	authzv1.ResourceRule
	cluster    string   `json:"clusters"`
	namespaces []string `json:"namespaces"`
}

func newPermissionRule(cluster string, namespaces []string, resourceRule authzv1.ResourceRule) *permissionRule {
	return &permissionRule{
		ResourceRule: resourceRule,
		cluster:      cluster,
		namespaces:   namespaces,
	}
}

func (p *permissionRule) key() string {
	shallow := *p
	shallow.cluster = ""
	data, _ := json.Marshal(shallow)
	h := fnv.New128a()
	h.Write(data)
	var sum [16]byte
	return fmt.Sprintf("%x", h.Sum(sum[:0]))
}

type PermissionRule struct {
	authzv1.ResourceRule
	Clusters   []string `json:"clusters"`
	Namespaces []string `json:"namespaces"`
}

// We use this to process a logic if a user is admin, all other
// permission rules are ignored.
type internalPermissionCache struct {
	// keeps a map of rules, the key is clusterName
	rulesMap map[string][]*permissionRule
}

func (i *internalPermissionCache) add(rules ...*permissionRule) {
	for _, rule := range rules {
		existing, ok := i.rulesMap[rule.cluster]
		if !ok {
			i.rulesMap[rule.cluster] = []*permissionRule{rule}
			continue
		}

		// if newly added rule is admin, override all the existing
		if equality.Semantic.DeepEqual(rule.ResourceRule, adminResourceRule) {
			i.rulesMap[rule.cluster] = []*permissionRule{rule}
			continue
		}

		var skipAdd bool
		for idx, existingRule := range existing {
			// if existing rule is admin, skip any update since it is the highest permission
			if equality.Semantic.DeepEqual(existingRule.ResourceRule, adminResourceRule) {
				skipAdd = true
				continue
			}

			// if resourceRule are the same, merge the namespace
			if equality.Semantic.DeepEqual(existingRule.ResourceRule, rule.ResourceRule) {
				namespaces := append(existingRule.namespaces, rule.namespaces...)
				i.rulesMap[rule.cluster][idx].namespaces = slices.Compact(namespaces)
				skipAdd = true
			}
		}

		if !skipAdd {
			i.rulesMap[rule.cluster] = append(existing, rule)
		}
	}
}

// consolidate all rules and return
func (i *internalPermissionCache) consolidateList() []PermissionRule {
	conslidatedMap := map[string]PermissionRule{}

	for _, clusterRules := range i.rulesMap {
		for _, rule := range clusterRules {
			// if two rules has the same namespaces and resourceRules, merge them to one.
			existing, ok := conslidatedMap[rule.key()]
			if !ok {
				conslidatedMap[rule.key()] = PermissionRule{
					ResourceRule: rule.ResourceRule,
					Clusters:     []string{rule.cluster},
					Namespaces:   rule.namespaces,
				}
			} else {
				existing.Clusters = append(existing.Clusters, rule.cluster)
				conslidatedMap[rule.key()] = existing
			}
		}
	}

	conslidatedList := []PermissionRule{}
	for _, rule := range conslidatedMap {
		conslidatedList = append(conslidatedList, rule)
	}
	return conslidatedList
}

// config needs to be the user's config to evaluate
func GetSelfPermissionRules(ctx context.Context, config *rest.Config, interestedVerb ...string) ([]PermissionRule, error) {
	clusterViewClient, err := clusterviewclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	userPermissions, err := clusterViewClient.ClusterviewV1alpha1().UserPermissions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	cache := &internalPermissionCache{
		rulesMap: make(map[string][]*permissionRule),
	}
	for _, userPermission := range userPermissions.Items {
		permissionRules := evaluateUserPermissionRule(userPermission, interestedVerb...)
		cache.add(permissionRules...)
	}
	return cache.consolidateList(), nil
}

func evaluateUserPermissionRule(userPermission clusterviewv1alpha1.UserPermission, interestedVerb ...string) []*permissionRule {
	var permissionRules []*permissionRule

	for _, rule := range userPermission.Status.ClusterRoleDefinition.Rules {
		// we conly care about * and filtered verbs
		haveSet := sets.New[string](rule.Verbs...)

		// If verbs in rule include "*", the output verb is "*"
		// If there are interested verbs, only return rules that has interested verbs included
		// return all verbs in the rule otherwise.
		var verbs []string
		if haveSet.HasAll("*") {
			verbs = rule.Verbs
		} else if len(interestedVerb) > 0 {
			if !haveSet.HasAll(interestedVerb...) {
				continue
			}
			verbs = interestedVerb
		} else {
			verbs = rule.Verbs
		}

		resourceRule := authzv1.ResourceRule{
			Verbs:     verbs,
			APIGroups: rule.APIGroups,
			Resources: rule.Resources,
		}
		for _, binding := range userPermission.Status.Bindings {
			p := newPermissionRule(binding.Cluster, binding.Namespaces, resourceRule)
			permissionRules = append(permissionRules, p)
		}
	}
	return permissionRules
}
