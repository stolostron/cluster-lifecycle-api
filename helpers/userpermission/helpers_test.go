package userpermission

import (
	"testing"

	clusterviewv1alpha1 "github.com/stolostron/cluster-lifecycle-api/clusterview/v1alpha1"
	authzv1 "k8s.io/api/authorization/v1"
	rbacv1 "k8s.io/api/rbac/v1"
)

func TestEvaluateUserPermissionRule(t *testing.T) {
	tcs := []struct {
		name           string
		userPermission clusterviewv1alpha1.UserPermission
		interestedVerb []string
		expected       []*permissionRule
	}{
		{
			name: "single rule single binding",
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"get", "list"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
					},
				},
			},
			expected: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
		{
			name: "single rule multiple bindings",
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"get", "list"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
						{
							Cluster:    "cluster2",
							Namespaces: []string{"kube-system"},
						},
					},
				},
			},
			expected: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
				{
					cluster:    "cluster2",
					namespaces: []string{"kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
		{
			name: "multiple rules single binding",
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"get", "list"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
							{
								Verbs:     []string{"get", "list", "watch"},
								APIGroups: []string{"apps"},
								Resources: []string{"deployments"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
					},
				},
			},
			expected: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list", "watch"},
						APIGroups: []string{"apps"},
						Resources: []string{"deployments"},
					},
				},
			},
		},
		{
			name: "admin rule",
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"*"},
								APIGroups: []string{"*"},
								Resources: []string{"*"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"*"},
						},
					},
				},
			},
			expected: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"*"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"*"},
						APIGroups: []string{"*"},
						Resources: []string{"*"},
					},
				},
			},
		},
		{
			name: "empty rules",
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
					},
				},
			},
			expected: []*permissionRule{},
		},
		{
			name: "empty bindings",
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"get", "list"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{},
				},
			},
			expected: []*permissionRule{},
		},
		{
			name:           "filter by interested verbs - match",
			interestedVerb: []string{"get"},
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"get", "list"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
					},
				},
			},
			expected: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
		{
			name:           "filter by interested verbs - no match",
			interestedVerb: []string{"delete"},
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"get", "list"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
					},
				},
			},
			expected: []*permissionRule{},
		},
		{
			name:           "wildcard verbs with interested verbs",
			interestedVerb: []string{"get"},
			userPermission: clusterviewv1alpha1.UserPermission{
				Status: clusterviewv1alpha1.UserPermissionStatus{
					ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
						Rules: []rbacv1.PolicyRule{
							{
								Verbs:     []string{"*"},
								APIGroups: []string{""},
								Resources: []string{"pods"},
							},
						},
					},
					Bindings: []clusterviewv1alpha1.ClusterBinding{
						{
							Cluster:    "cluster1",
							Namespaces: []string{"default"},
						},
					},
				},
			},
			expected: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"*"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := evaluateUserPermissionRule(tc.userPermission, tc.interestedVerb...)
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d permission rules, got %d", len(tc.expected), len(result))
				return
			}
			for i, rule := range result {
				if rule.cluster != tc.expected[i].cluster {
					t.Errorf("expected cluster %s, got %s", tc.expected[i].cluster, rule.cluster)
				}
				if len(rule.namespaces) != len(tc.expected[i].namespaces) {
					t.Errorf("expected %d namespaces, got %d", len(tc.expected[i].namespaces), len(rule.namespaces))
				}
				for j, ns := range rule.namespaces {
					if ns != tc.expected[i].namespaces[j] {
						t.Errorf("expected namespace %s, got %s", tc.expected[i].namespaces[j], ns)
					}
				}
				if len(rule.Verbs) != len(tc.expected[i].Verbs) {
					t.Errorf("expected %d verbs, got %d", len(tc.expected[i].Verbs), len(rule.Verbs))
				}
				if len(rule.APIGroups) != len(tc.expected[i].APIGroups) {
					t.Errorf("expected %d api groups, got %d", len(tc.expected[i].APIGroups), len(rule.APIGroups))
				}
				if len(rule.Resources) != len(tc.expected[i].Resources) {
					t.Errorf("expected %d resources, got %d", len(tc.expected[i].Resources), len(rule.Resources))
				}
			}
		})
	}
}

func TestInternalPermissionCache_Add(t *testing.T) {
	tcs := []struct {
		name     string
		initial  map[string][]*permissionRule
		toAdd    []*permissionRule
		expected map[string][]*permissionRule
	}{
		{
			name:    "add to empty cache",
			initial: map[string][]*permissionRule{},
			toAdd: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
			expected: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get", "list"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
		},
		{
			name: "add to existing cluster",
			initial: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"list"},
						APIGroups: []string{"apps"},
						Resources: []string{"deployments"},
					},
				},
			},
			expected: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
					{
						cluster:    "cluster1",
						namespaces: []string{"kube-system"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"list"},
							APIGroups: []string{"apps"},
							Resources: []string{"deployments"},
						},
					},
				},
			},
		},
		{
			name: "add admin rule overrides existing",
			initial: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"*"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"*"},
						APIGroups: []string{"*"},
						Resources: []string{"*"},
					},
				},
			},
			expected: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"*"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"*"},
							APIGroups: []string{"*"},
							Resources: []string{"*"},
						},
					},
				},
			},
		},
		{
			name: "add to different clusters",
			initial: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []*permissionRule{
				{
					cluster:    "cluster2",
					namespaces: []string{"kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"list"},
						APIGroups: []string{"apps"},
						Resources: []string{"deployments"},
					},
				},
			},
			expected: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
				"cluster2": {
					{
						cluster:    "cluster2",
						namespaces: []string{"kube-system"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"list"},
							APIGroups: []string{"apps"},
							Resources: []string{"deployments"},
						},
					},
				},
			},
		},
		{
			name: "merge same resource rules with different namespaces",
			initial: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get", "list"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []*permissionRule{
				{
					cluster:    "cluster1",
					namespaces: []string{"kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
			expected: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default", "kube-system"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get", "list"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cache := &internalPermissionCache{
				rulesMap: tc.initial,
			}
			cache.add(tc.toAdd...)

			if len(cache.rulesMap) != len(tc.expected) {
				t.Errorf("expected %d clusters in cache, got %d", len(tc.expected), len(cache.rulesMap))
			}

			for cluster, expectedRules := range tc.expected {
				actualRules, ok := cache.rulesMap[cluster]
				if !ok {
					t.Errorf("expected cluster %s to be in cache, but it was not", cluster)
					continue
				}
				if len(actualRules) != len(expectedRules) {
					t.Errorf("for cluster %s, expected %d rules, got %d", cluster, len(expectedRules), len(actualRules))
				}
			}
		})
	}
}

func TestEvaluateUserPermissionRule_WithConsolidation(t *testing.T) {
	tcs := []struct {
		name            string
		userPermissions []clusterviewv1alpha1.UserPermission
		expected        []PermissionRule
	}{
		{
			name: "same resource rule across multiple clusters - should consolidate",
			userPermissions: []clusterviewv1alpha1.UserPermission{
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get", "list"},
									APIGroups: []string{""},
									Resources: []string{"pods"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"default"},
							},
							{
								Cluster:    "cluster2",
								Namespaces: []string{"default"},
							},
						},
					},
				},
			},
			expected: []PermissionRule{
				{
					Clusters:   []string{"cluster1", "cluster2"},
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
		{
			name: "multiple rules with different namespaces - separate consolidation",
			userPermissions: []clusterviewv1alpha1.UserPermission{
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get"},
									APIGroups: []string{"apps"},
									Resources: []string{"deployments"},
								},
								{
									Verbs:     []string{"list"},
									APIGroups: []string{""},
									Resources: []string{"services"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"default"},
							},
							{
								Cluster:    "cluster2",
								Namespaces: []string{"default"},
							},
						},
					},
				},
			},
			expected: []PermissionRule{
				{
					Clusters:   []string{"cluster1", "cluster2"},
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get"},
						APIGroups: []string{"apps"},
						Resources: []string{"deployments"},
					},
				},
				{
					Clusters:   []string{"cluster1", "cluster2"},
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"list"},
						APIGroups: []string{""},
						Resources: []string{"services"},
					},
				},
			},
		},
		{
			name: "multiple UserPermissions - consolidate across them",
			userPermissions: []clusterviewv1alpha1.UserPermission{
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get"},
									APIGroups: []string{""},
									Resources: []string{"pods"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"default"},
							},
						},
					},
				},
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get"},
									APIGroups: []string{""},
									Resources: []string{"pods"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster2",
								Namespaces: []string{"default"},
							},
						},
					},
				},
			},
			expected: []PermissionRule{
				{
					Clusters:   []string{"cluster1", "cluster2"},
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
		{
			name: "multiple UserPermissions with namespace merging within cluster",
			userPermissions: []clusterviewv1alpha1.UserPermission{
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get"},
									APIGroups: []string{""},
									Resources: []string{"pods"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"default"},
							},
						},
					},
				},
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get"},
									APIGroups: []string{""},
									Resources: []string{"pods"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"kube-system"},
							},
						},
					},
				},
			},
			expected: []PermissionRule{
				{
					Clusters:   []string{"cluster1"},
					Namespaces: []string{"default", "kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
		},
		{
			name: "multiple UserPermissions with admin override",
			userPermissions: []clusterviewv1alpha1.UserPermission{
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"get"},
									APIGroups: []string{""},
									Resources: []string{"pods"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"default"},
							},
						},
					},
				},
				{
					Status: clusterviewv1alpha1.UserPermissionStatus{
						ClusterRoleDefinition: clusterviewv1alpha1.ClusterRoleDefinition{
							Rules: []rbacv1.PolicyRule{
								{
									Verbs:     []string{"*"},
									APIGroups: []string{"*"},
									Resources: []string{"*"},
								},
							},
						},
						Bindings: []clusterviewv1alpha1.ClusterBinding{
							{
								Cluster:    "cluster1",
								Namespaces: []string{"*"},
							},
						},
					},
				},
			},
			expected: []PermissionRule{
				{
					Clusters:   []string{"cluster1"},
					Namespaces: []string{"*"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"*"},
						APIGroups: []string{"*"},
						Resources: []string{"*"},
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			// Use the cache to test consolidation
			cache := &internalPermissionCache{
				rulesMap: make(map[string][]*permissionRule),
			}
			// Process all UserPermissions
			for _, userPermission := range tc.userPermissions {
				permissionRules := evaluateUserPermissionRule(userPermission)
				cache.add(permissionRules...)
			}
			result := cache.consolidateList()

			if len(result) != len(tc.expected) {
				t.Errorf("expected %d consolidated rules, got %d", len(tc.expected), len(result))
				return
			}

			// Since map iteration order is not guaranteed, we need to match rules by their content
			for _, expectedRule := range tc.expected {
				found := false
				for _, actualRule := range result {
					if len(actualRule.Verbs) == len(expectedRule.Verbs) &&
						len(actualRule.APIGroups) == len(expectedRule.APIGroups) &&
						len(actualRule.Resources) == len(expectedRule.Resources) &&
						len(actualRule.Namespaces) == len(expectedRule.Namespaces) {
						// Check if clusters match (order may differ)
						if len(actualRule.Clusters) == len(expectedRule.Clusters) {
							clustersMatch := true
							for _, ec := range expectedRule.Clusters {
								clusterFound := false
								for _, ac := range actualRule.Clusters {
									if ec == ac {
										clusterFound = true
										break
									}
								}
								if !clusterFound {
									clustersMatch = false
									break
								}
							}
							if clustersMatch {
								found = true
								break
							}
						}
					}
				}
				if !found {
					t.Errorf("expected rule not found in results: %+v", expectedRule)
				}
			}
		})
	}
}

func TestInternalPermissionCache_ConsolidateList(t *testing.T) {
	tcs := []struct {
		name        string
		rulesMap    map[string][]*permissionRule
		expectedLen int
		validate    func(t *testing.T, result []PermissionRule)
	}{
		{
			name:        "empty cache",
			rulesMap:    map[string][]*permissionRule{},
			expectedLen: 0,
		},
		{
			name: "single cluster single rule",
			rulesMap: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			expectedLen: 1,
			validate: func(t *testing.T, result []PermissionRule) {
				if len(result[0].Clusters) != 1 || result[0].Clusters[0] != "cluster1" {
					t.Errorf("expected clusters [cluster1], got %v", result[0].Clusters)
				}
			},
		},
		{
			name: "single cluster multiple rules",
			rulesMap: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
					{
						cluster:    "cluster1",
						namespaces: []string{"kube-system"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"list"},
							APIGroups: []string{"apps"},
							Resources: []string{"deployments"},
						},
					},
				},
			},
			expectedLen: 2,
		},
		{
			name: "multiple clusters with same rule - should consolidate",
			rulesMap: map[string][]*permissionRule{
				"cluster1": {
					{
						cluster:    "cluster1",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
				"cluster2": {
					{
						cluster:    "cluster2",
						namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			expectedLen: 1,
			validate: func(t *testing.T, result []PermissionRule) {
				if len(result) == 0 {
					t.Fatal("expected at least one result")
				}
				if len(result[0].Clusters) != 2 {
					t.Errorf("expected 2 clusters to be consolidated, got %d: %v", len(result[0].Clusters), result[0].Clusters)
				}
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cache := &internalPermissionCache{
				rulesMap: tc.rulesMap,
			}
			result := cache.consolidateList()
			if len(result) != tc.expectedLen {
				t.Errorf("expected %d rules, got %d", tc.expectedLen, len(result))
			}
			if tc.validate != nil {
				tc.validate(t, result)
			}
		})
	}
}
