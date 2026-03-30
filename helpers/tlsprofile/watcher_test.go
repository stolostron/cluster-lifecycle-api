// Copyright Contributors to the Open Cluster Management project

package tlsprofile

import (
	"context"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func TestTLSProfileChanged(t *testing.T) {
	tests := []struct {
		name     string
		old      *configv1.TLSSecurityProfile
		new      *configv1.TLSSecurityProfile
		expected bool
	}{
		{
			name:     "Both nil - no change",
			old:      nil,
			new:      nil,
			expected: false,
		},
		{
			name:     "Old nil, new set - changed",
			old:      nil,
			new:      &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
			expected: true,
		},
		{
			name:     "Old set, new nil - changed",
			old:      &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
			new:      nil,
			expected: true,
		},
		{
			name:     "Same profile type - no change",
			old:      &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
			new:      &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
			expected: false,
		},
		{
			name:     "Different profile types - changed",
			old:      &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
			new:      &configv1.TLSSecurityProfile{Type: configv1.TLSProfileModernType},
			expected: true,
		},
		{
			name: "Same custom profile - no change",
			old: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
					},
				},
			},
			new: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
					},
				},
			},
			expected: false,
		},
		{
			name: "Custom profile with different MinTLSVersion - changed",
			old: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
					},
				},
			},
			new: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS13,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
					},
				},
			},
			expected: true,
		},
		{
			name: "Custom profile with different ciphers - changed",
			old: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
					},
				},
			},
			new: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"},
					},
				},
			},
			expected: true,
		},
		{
			name: "Custom profile with additional cipher - changed",
			old: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers:       []string{"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"},
					},
				},
			},
			new: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
						Ciphers: []string{
							"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
							"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
						},
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tlsProfileChanged(tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("tlsProfileChanged() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestGetTLSProfileType(t *testing.T) {
	tests := []struct {
		name     string
		profile  *configv1.TLSSecurityProfile
		expected string
	}{
		{
			name:     "Nil profile",
			profile:  nil,
			expected: "nil (default Intermediate)",
		},
		{
			name:     "Old profile type",
			profile:  &configv1.TLSSecurityProfile{Type: configv1.TLSProfileOldType},
			expected: string(configv1.TLSProfileOldType),
		},
		{
			name:     "Intermediate profile type",
			profile:  &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
			expected: string(configv1.TLSProfileIntermediateType),
		},
		{
			name:     "Modern profile type",
			profile:  &configv1.TLSSecurityProfile{Type: configv1.TLSProfileModernType},
			expected: string(configv1.TLSProfileModernType),
		},
		{
			name:     "Custom profile type",
			profile:  &configv1.TLSSecurityProfile{Type: configv1.TLSProfileCustomType},
			expected: string(configv1.TLSProfileCustomType),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTLSProfileType(tt.profile)
			if result != tt.expected {
				t.Errorf("getTLSProfileType() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestStartTLSProfileWatcher_InvalidHost(t *testing.T) {
	// Test with unreachable host - watcher should return error
	ctx := context.Background()
	cancel := func() {}

	// Create a config with invalid host that will timeout
	invalidConfig := &rest.Config{
		Host:    "https://invalid-host-that-does-not-exist.local:6443",
		Timeout: 1, // 1 nanosecond to fail fast
	}

	// This will fail to get initial profile and should return an error
	err := StartTLSProfileWatcher(ctx, invalidConfig, cancel)
	if err == nil {
		t.Error("Expected error for invalid host, got nil")
	}
}

func TestStartTLSProfileWatcher_CancelledContext(t *testing.T) {
	// Test with pre-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	invalidConfig := &rest.Config{
		Host:    "https://invalid-host.local:6443",
		Timeout: 1,
	}

	// Should return error for invalid config
	err := StartTLSProfileWatcher(ctx, invalidConfig, cancel)
	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestHandleAPIServerUpdate(t *testing.T) {
	tests := []struct {
		name         string
		oldObj       interface{}
		newObj       interface{}
		expectCancel bool
	}{
		{
			name: "TLS profile changed - should cancel",
			oldObj: &configv1.APIServer{
				ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
				Spec: configv1.APIServerSpec{
					TLSSecurityProfile: &configv1.TLSSecurityProfile{
						Type: configv1.TLSProfileIntermediateType,
					},
				},
			},
			newObj: &configv1.APIServer{
				ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
				Spec: configv1.APIServerSpec{
					TLSSecurityProfile: &configv1.TLSSecurityProfile{
						Type: configv1.TLSProfileModernType,
					},
				},
			},
			expectCancel: true,
		},
		{
			name: "TLS profile unchanged - should not cancel",
			oldObj: &configv1.APIServer{
				ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
				Spec: configv1.APIServerSpec{
					TLSSecurityProfile: &configv1.TLSSecurityProfile{
						Type: configv1.TLSProfileIntermediateType,
					},
				},
			},
			newObj: &configv1.APIServer{
				ObjectMeta: metav1.ObjectMeta{Name: "cluster"},
				Spec: configv1.APIServerSpec{
					TLSSecurityProfile: &configv1.TLSSecurityProfile{
						Type: configv1.TLSProfileIntermediateType,
					},
				},
			},
			expectCancel: false,
		},
		{
			name: "Wrong APIServer name - should not cancel",
			oldObj: &configv1.APIServer{
				ObjectMeta: metav1.ObjectMeta{Name: "other"},
				Spec: configv1.APIServerSpec{
					TLSSecurityProfile: &configv1.TLSSecurityProfile{
						Type: configv1.TLSProfileIntermediateType,
					},
				},
			},
			newObj: &configv1.APIServer{
				ObjectMeta: metav1.ObjectMeta{Name: "other"},
				Spec: configv1.APIServerSpec{
					TLSSecurityProfile: &configv1.TLSSecurityProfile{
						Type: configv1.TLSProfileModernType,
					},
				},
			},
			expectCancel: false,
		},
		{
			name:         "Invalid old object type - should not cancel",
			oldObj:       "not an APIServer",
			newObj:       &configv1.APIServer{ObjectMeta: metav1.ObjectMeta{Name: "cluster"}},
			expectCancel: false,
		},
		{
			name:         "Invalid new object type - should not cancel",
			oldObj:       &configv1.APIServer{ObjectMeta: metav1.ObjectMeta{Name: "cluster"}},
			newObj:       "not an APIServer",
			expectCancel: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cancelled := false
			cancel := func() {
				cancelled = true
			}

			handleAPIServerUpdate(tt.oldObj, tt.newObj, cancel)

			if cancelled != tt.expectCancel {
				t.Errorf("Expected cancel=%v, got cancel=%v", tt.expectCancel, cancelled)
			}
		})
	}
}

func TestTLSProfileChanged_EdgeCases(t *testing.T) {
	// Additional edge cases to improve coverage
	tests := []struct {
		name     string
		old      *configv1.TLSSecurityProfile
		new      *configv1.TLSSecurityProfile
		expected bool
	}{
		{
			name: "Custom profile: old nil custom, new with custom",
			old: &configv1.TLSSecurityProfile{
				Type:   configv1.TLSProfileCustomType,
				Custom: nil,
			},
			new: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
					},
				},
			},
			expected: true,
		},
		{
			name: "Custom profile: both nil custom",
			old: &configv1.TLSSecurityProfile{
				Type:   configv1.TLSProfileCustomType,
				Custom: nil,
			},
			new: &configv1.TLSSecurityProfile{
				Type:   configv1.TLSProfileCustomType,
				Custom: nil,
			},
			expected: false,
		},
		{
			name: "Non-custom profile: same type, different custom (should not compare custom)",
			old: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileIntermediateType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS12,
					},
				},
			},
			new: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileIntermediateType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS13,
					},
				},
			},
			expected: false, // For non-custom profiles, custom field is ignored
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tlsProfileChanged(tt.old, tt.new)
			if result != tt.expected {
				t.Errorf("tlsProfileChanged() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
