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
		expected       []PermissionRule
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
			expected: []PermissionRule{
				{
					Cluster:    "cluster1",
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
			expected: []PermissionRule{
				{
					Cluster:    "cluster1",
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
				{
					Cluster:    "cluster2",
					Namespaces: []string{"kube-system"},
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
			expected: []PermissionRule{
				{
					Cluster:    "cluster1",
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
				{
					Cluster:    "cluster1",
					Namespaces: []string{"default"},
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
			expected: []PermissionRule{
				{
					Cluster:    "cluster1",
					Namespaces: []string{"*"},
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
			expected: []PermissionRule{},
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
			expected: []PermissionRule{},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			result := evaluateUserPermissionRule(tc.userPermission)
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d permission rules, got %d", len(tc.expected), len(result))
				return
			}
			for i, rule := range result {
				if rule.Cluster != tc.expected[i].Cluster {
					t.Errorf("expected cluster %s, got %s", tc.expected[i].Cluster, rule.Cluster)
				}
				if len(rule.Namespaces) != len(tc.expected[i].Namespaces) {
					t.Errorf("expected %d namespaces, got %d", len(tc.expected[i].Namespaces), len(rule.Namespaces))
				}
				for j, ns := range rule.Namespaces {
					if ns != tc.expected[i].Namespaces[j] {
						t.Errorf("expected namespace %s, got %s", tc.expected[i].Namespaces[j], ns)
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
		initial  map[string][]PermissionRule
		toAdd    []PermissionRule
		expected map[string][]PermissionRule
	}{
		{
			name:    "add to empty cache",
			initial: map[string][]PermissionRule{},
			toAdd: []PermissionRule{
				{
					Cluster:    "cluster1",
					Namespaces: []string{"default"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"get", "list"},
						APIGroups: []string{""},
						Resources: []string{"pods"},
					},
				},
			},
			expected: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
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
			initial: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []PermissionRule{
				{
					Cluster:    "cluster1",
					Namespaces: []string{"kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"list"},
						APIGroups: []string{"apps"},
						Resources: []string{"deployments"},
					},
				},
			},
			expected: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
					{
						Cluster:    "cluster1",
						Namespaces: []string{"kube-system"},
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
			initial: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []PermissionRule{
				{
					Cluster:    "cluster1",
					Namespaces: []string{"*"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"*"},
						APIGroups: []string{"*"},
						Resources: []string{"*"},
					},
				},
			},
			expected: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"*"},
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
			initial: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			toAdd: []PermissionRule{
				{
					Cluster:    "cluster2",
					Namespaces: []string{"kube-system"},
					ResourceRule: authzv1.ResourceRule{
						Verbs:     []string{"list"},
						APIGroups: []string{"apps"},
						Resources: []string{"deployments"},
					},
				},
			},
			expected: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
				"cluster2": {
					{
						Cluster:    "cluster2",
						Namespaces: []string{"kube-system"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"list"},
							APIGroups: []string{"apps"},
							Resources: []string{"deployments"},
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

func TestInternalPermissionCache_List(t *testing.T) {
	tcs := []struct {
		name        string
		rulesMap    map[string][]PermissionRule
		expectedLen int
	}{
		{
			name:        "empty cache",
			rulesMap:    map[string][]PermissionRule{},
			expectedLen: 0,
		},
		{
			name: "single cluster single rule",
			rulesMap: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			expectedLen: 1,
		},
		{
			name: "single cluster multiple rules",
			rulesMap: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
					{
						Cluster:    "cluster1",
						Namespaces: []string{"kube-system"},
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
			name: "multiple clusters",
			rulesMap: map[string][]PermissionRule{
				"cluster1": {
					{
						Cluster:    "cluster1",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
				"cluster2": {
					{
						Cluster:    "cluster2",
						Namespaces: []string{"default"},
						ResourceRule: authzv1.ResourceRule{
							Verbs:     []string{"get"},
							APIGroups: []string{""},
							Resources: []string{"pods"},
						},
					},
				},
			},
			expectedLen: 2,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			cache := &internalPermissionCache{
				rulesMap: tc.rulesMap,
			}
			result := cache.list()
			if len(result) != tc.expectedLen {
				t.Errorf("expected %d rules, got %d", tc.expectedLen, len(result))
			}
		})
	}
}
