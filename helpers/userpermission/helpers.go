package userpermission

import (
	"context"
	clusterviewclientset "github.com/stolostron/cluster-lifecycle-api/client/clusterview/clientset/versioned"
	clusterviewv1alpha1 "github.com/stolostron/cluster-lifecycle-api/clusterview/v1alpha1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	authzv1 "k8s.io/api/authorization/v1"
	"k8s.io/client-go/rest"
)

var adminResourceRule = authzv1.ResourceRule{
	Verbs:     []string{"*"},
	APIGroups: []string{"*"},
	Resources: []string{"*"},
}

type PermissionRule struct {
	authzv1.ResourceRule
	Cluster    string   `json:"cluster"`
	Namespaces []string `json:"namespaces"`
}

// We use this to process a logic if a user is admin, all other
// permission rules are ignored.
type internalPermissionCache struct {
	// keeps a map of rules, the key is clusterName
	rulesMap map[string][]PermissionRule
}

func (i *internalPermissionCache) add(rules ...PermissionRule) {
	for _, rule := range rules {
		existing, ok := i.rulesMap[rule.Cluster]
		if !ok {
			i.rulesMap[rule.Cluster] = []PermissionRule{rule}
			continue
		}

		// if existing rule is admin, skip any update since it is highest permission
		if equality.Semantic.DeepEqual(existing, adminResourceRule) {
			continue
		}

		// if newly added rule is admin, override all the existing
		if equality.Semantic.DeepEqual(rule.ResourceRule, adminResourceRule) {
			i.rulesMap[rule.Cluster] = []PermissionRule{rule}
			continue
		}
		i.rulesMap[rule.Cluster] = append(existing, rule)
	}
}

func (i *internalPermissionCache) list() []PermissionRule {
	var rules []PermissionRule
	for _, clusterRules := range i.rulesMap {
		for _, rule := range clusterRules {
			rules = append(rules, rule)
		}
	}
	return rules
}

// config needs to be the user's config to evaluate
func GetSelfPermissionRules(ctx context.Context, config *rest.Config) ([]PermissionRule, error) {
	clusterViewClient, err := clusterviewclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	userPermissions, err := clusterViewClient.ClusterviewV1alpha1().UserPermissions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	cache := &internalPermissionCache{
		rulesMap: make(map[string][]PermissionRule),
	}
	for _, userPermission := range userPermissions.Items {
		permissionRules := evaluateUserPermissionRule(userPermission)
		cache.add(permissionRules...)
	}
	return cache.list(), nil
}

func evaluateUserPermissionRule(userPermission clusterviewv1alpha1.UserPermission) []PermissionRule {
	var permissionRules []PermissionRule

	for _, rule := range userPermission.Status.ClusterRoleDefinition.Rules {
		resourceRule := authzv1.ResourceRule{
			Verbs:     rule.Verbs,
			APIGroups: rule.APIGroups,
			Resources: rule.Resources,
		}
		for _, binding := range userPermission.Status.Bindings {
			permissionRule := PermissionRule{
				Cluster:      binding.Cluster,
				Namespaces:   binding.Namespaces,
				ResourceRule: resourceRule,
			}
			permissionRules = append(permissionRules, permissionRule)
		}
	}
	return permissionRules
}
