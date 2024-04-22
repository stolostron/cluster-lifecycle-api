package klusterletconfig

import (
	"reflect"
	"testing"

	klusterletconfigv1alpha1 "github.com/stolostron/cluster-lifecycle-api/klusterletconfig/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

func TestCheckAllFiledHasMergeFunc(test *testing.T) {
	// TestCheckAllFiledHasMergeFunc tests if all fields in the KlusterletConfigSpec have a corresponding merge function in the merge.go file.
	// When adding a new field to the KlusterletConfigSpec, if the merge function is not provided, it will fail in the unit tests.

	kcSpec := &klusterletconfigv1alpha1.KlusterletConfigSpec{}
	v := reflect.ValueOf(kcSpec).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		if _, ok := klusterletConfigMergeFuncs[fieldName]; !ok {
			test.Errorf("merge function for field %s is not provided", fieldName)
		}
	}
}

func TestMergeKlusterletConfigs(test *testing.T) {
	testcases := []struct {
		name     string
		kcs      []*klusterletconfigv1alpha1.KlusterletConfig
		expected *klusterletconfigv1alpha1.KlusterletConfig
	}{
		{
			name: "override strategy: override the base value if next KlusterletConfig in the list has a non-zero value",
			kcs: []*klusterletconfigv1alpha1.KlusterletConfig{
				{
					Spec: klusterletconfigv1alpha1.KlusterletConfigSpec{
						Registries: []klusterletconfigv1alpha1.Registries{
							{
								Mirror: "mirror",
								Source: "source",
							},
						},
						PullSecret: corev1.ObjectReference{
							Name: "pull-secret",
						},
						NodePlacement: &operatorv1.NodePlacement{
							NodeSelector: map[string]string{
								"key": "value",
							},
						},
						HubKubeAPIServerProxyConfig: klusterletconfigv1alpha1.KubeAPIServerProxyConfig{
							HTTPProxy:  "http://proxy",
							HTTPSProxy: "https://proxy",
						},
						HubKubeAPIServerURL:      "https://hub",
						HubKubeAPIServerCABundle: []byte("ca-bundle"),
					},
				},
				{
					Spec: klusterletconfigv1alpha1.KlusterletConfigSpec{
						Registries: []klusterletconfigv1alpha1.Registries{
							{
								Mirror: "mirror2",
								Source: "source2",
							},
						},
						PullSecret: corev1.ObjectReference{
							Name: "pull-secret2",
						},
						NodePlacement: &operatorv1.NodePlacement{
							NodeSelector: map[string]string{
								"key2": "value2",
							},
						},
						HubKubeAPIServerProxyConfig: klusterletconfigv1alpha1.KubeAPIServerProxyConfig{
							HTTPProxy:  "http://proxy2",
							HTTPSProxy: "https://proxy2",
						},
						HubKubeAPIServerURL:      "https://hub2",
						HubKubeAPIServerCABundle: []byte("ca-bundle2"),
					},
				},
			},
			expected: &klusterletconfigv1alpha1.KlusterletConfig{
				Spec: klusterletconfigv1alpha1.KlusterletConfigSpec{
					Registries: []klusterletconfigv1alpha1.Registries{
						{
							Mirror: "mirror2",
							Source: "source2",
						},
					},
					PullSecret: corev1.ObjectReference{
						Name: "pull-secret2",
					},
					NodePlacement: &operatorv1.NodePlacement{
						NodeSelector: map[string]string{
							"key2": "value2",
						},
					},
					HubKubeAPIServerProxyConfig: klusterletconfigv1alpha1.KubeAPIServerProxyConfig{
						HTTPProxy:  "http://proxy2",
						HTTPSProxy: "https://proxy2",
					},
					HubKubeAPIServerURL:      "https://hub2",
					HubKubeAPIServerCABundle: []byte("ca-bundle2"),
				},
			},
		},
		{
			name: "override strategy: flow to the next if the first is zero",
			kcs: []*klusterletconfigv1alpha1.KlusterletConfig{
				{
					Spec: klusterletconfigv1alpha1.KlusterletConfigSpec{
						Registries:                  nil,
						PullSecret:                  corev1.ObjectReference{},
						NodePlacement:               nil,
						HubKubeAPIServerProxyConfig: klusterletconfigv1alpha1.KubeAPIServerProxyConfig{},
						HubKubeAPIServerURL:         "",
						HubKubeAPIServerCABundle:    nil,
					},
				},
				{
					Spec: klusterletconfigv1alpha1.KlusterletConfigSpec{
						Registries: []klusterletconfigv1alpha1.Registries{
							{
								Mirror: "mirror2",
								Source: "source2",
							},
						},
						PullSecret: corev1.ObjectReference{
							Name: "pull-secret2",
						},
						NodePlacement: &operatorv1.NodePlacement{
							NodeSelector: map[string]string{
								"key2": "value2",
							},
						},
						HubKubeAPIServerProxyConfig: klusterletconfigv1alpha1.KubeAPIServerProxyConfig{
							HTTPProxy:  "http://proxy2",
							HTTPSProxy: "https://proxy2",
						},
						HubKubeAPIServerURL:      "https://hub2",
						HubKubeAPIServerCABundle: []byte("ca-bundle2"),
					},
				},
			},
			expected: &klusterletconfigv1alpha1.KlusterletConfig{
				Spec: klusterletconfigv1alpha1.KlusterletConfigSpec{
					Registries: []klusterletconfigv1alpha1.Registries{
						{
							Mirror: "mirror2",
							Source: "source2",
						},
					},
					PullSecret: corev1.ObjectReference{
						Name: "pull-secret2",
					},
					NodePlacement: &operatorv1.NodePlacement{
						NodeSelector: map[string]string{
							"key2": "value2",
						},
					},
					HubKubeAPIServerProxyConfig: klusterletconfigv1alpha1.KubeAPIServerProxyConfig{
						HTTPProxy:  "http://proxy2",
						HTTPSProxy: "https://proxy2",
					},
					HubKubeAPIServerURL:      "https://hub2",
					HubKubeAPIServerCABundle: []byte("ca-bundle2"),
				},
			},
		},
	}

	for _, testcase := range testcases {
		test.Run(testcase.name, func(test *testing.T) {
			merged, err := MergeKlusterletConfigs(testcase.kcs...)
			if err != nil {
				test.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(merged, testcase.expected) {
				test.Errorf("expected: %v, got: %v", testcase.expected, merged)
			}
		})
	}
}
