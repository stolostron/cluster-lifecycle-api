// Copyright Contributors to the Open Cluster Management project

package tlsprofile

import (
	"crypto/tls"
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/client-go/rest"
)

func TestConvertTLSProfileToConfig(t *testing.T) {
	tests := []struct {
		name           string
		profile        *configv1.TLSSecurityProfile
		expectedMinTLS uint16
		expectCiphers  bool
	}{
		{
			name:           "Nil profile defaults to TLS 1.2",
			profile:        nil,
			expectedMinTLS: tls.VersionTLS12,
			expectCiphers:  false,
		},
		{
			name: "Old profile uses TLS 1.0",
			profile: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileOldType,
			},
			expectedMinTLS: tls.VersionTLS10,
			expectCiphers:  true,
		},
		{
			name: "Intermediate profile uses TLS 1.2",
			profile: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileIntermediateType,
			},
			expectedMinTLS: tls.VersionTLS12,
			expectCiphers:  true,
		},
		{
			name: "Modern profile uses TLS 1.3",
			profile: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileModernType,
			},
			expectedMinTLS: tls.VersionTLS13,
			expectCiphers:  false, // TLS 1.3 doesn't use configurable ciphers
		},
		{
			name: "Custom profile with TLS 1.2",
			profile: &configv1.TLSSecurityProfile{
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
			expectedMinTLS: tls.VersionTLS12,
			expectCiphers:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ConvertTLSProfileToConfig(tt.profile)

			if config == nil {
				t.Fatal("Expected non-nil TLS config")
			}

			if config.MinVersion != tt.expectedMinTLS {
				t.Errorf("Expected MinVersion %v, got %v", tt.expectedMinTLS, config.MinVersion)
			}

			if tt.expectCiphers && len(config.CipherSuites) == 0 {
				t.Error("Expected cipher suites to be configured")
			}

			if !tt.expectCiphers && len(config.CipherSuites) > 0 {
				t.Errorf("Expected no cipher suites for this profile, got %d", len(config.CipherSuites))
			}
		})
	}
}

func TestGetTLSVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected uint16
	}{
		{"TLS 1.0", string(configv1.VersionTLS10), tls.VersionTLS10},
		{"TLS 1.1", string(configv1.VersionTLS11), tls.VersionTLS11},
		{"TLS 1.2", string(configv1.VersionTLS12), tls.VersionTLS12},
		{"TLS 1.3", string(configv1.VersionTLS13), tls.VersionTLS13},
		{"Unknown defaults to TLS 1.2", "unknown", tls.VersionTLS12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTLSVersion(tt.version)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetCipherSuites(t *testing.T) {
	cipherNames := []string{
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
		"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
		"UNKNOWN_CIPHER", // Should be ignored with a warning
	}

	ciphers := getCipherSuites(cipherNames)

	// Should have 2 valid ciphers (the unknown one is skipped)
	if len(ciphers) != 2 {
		t.Errorf("Expected 2 cipher suites, got %d", len(ciphers))
	}

	// Verify the ciphers match
	expectedCiphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	}

	for i, expected := range expectedCiphers {
		if ciphers[i] != expected {
			t.Errorf("Expected cipher %v at index %d, got %v", expected, i, ciphers[i])
		}
	}
}

func TestGetCipherSuitesWithOpenShiftFormat(t *testing.T) {
	// Test OpenShift hyphenated format
	cipherNames := []string{
		"ECDHE-RSA-AES128-GCM-SHA256",
		"ECDHE-ECDSA-CHACHA20-POLY1305",
		"AES128-GCM-SHA256",
	}

	ciphers := getCipherSuites(cipherNames)

	if len(ciphers) != 3 {
		t.Errorf("Expected 3 cipher suites, got %d", len(ciphers))
	}

	expectedCiphers := []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	}

	for i, expected := range expectedCiphers {
		if ciphers[i] != expected {
			t.Errorf("Expected cipher %v at index %d, got %v", expected, i, ciphers[i])
		}
	}
}

func TestGetTLSVersionName(t *testing.T) {
	tests := []struct {
		name     string
		version  uint16
		expected string
	}{
		{"TLS 1.0", tls.VersionTLS10, "TLS 1.0"},
		{"TLS 1.1", tls.VersionTLS11, "TLS 1.1"},
		{"TLS 1.2", tls.VersionTLS12, "TLS 1.2"},
		{"TLS 1.3", tls.VersionTLS13, "TLS 1.3"},
		{"Unknown version", 0x9999, "Unknown (0x9999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTLSVersionName(tt.version)
			if result != tt.expected {
				t.Errorf("getTLSVersionName(%v) = %v, expected %v", tt.version, result, tt.expected)
			}
		})
	}
}

func TestConvertTLSProfileToConfigWithCustomNil(t *testing.T) {
	// Test custom profile with nil Custom field
	profile := &configv1.TLSSecurityProfile{
		Type:   configv1.TLSProfileCustomType,
		Custom: nil,
	}

	config := ConvertTLSProfileToConfig(profile)

	if config == nil {
		t.Fatal("Expected non-nil TLS config")
	}

	if config.MinVersion != tls.VersionTLS12 {
		t.Errorf("Expected MinVersion TLS 1.2 for custom profile with nil Custom, got %v", config.MinVersion)
	}
}

func TestConvertTLSProfileToConfigUnknownType(t *testing.T) {
	// Test unknown profile type
	profile := &configv1.TLSSecurityProfile{
		Type: "UnknownType",
	}

	config := ConvertTLSProfileToConfig(profile)

	if config == nil {
		t.Fatal("Expected non-nil TLS config")
	}

	if config.MinVersion != tls.VersionTLS12 {
		t.Errorf("Expected MinVersion TLS 1.2 for unknown profile type, got %v", config.MinVersion)
	}
}

func TestGetCipherSuiteMap(t *testing.T) {
	// Test that cipher suite map contains expected entries
	cipherMap := getCipherSuiteMap()

	// Test a few key mappings
	expectedMappings := map[string]uint16{
		"TLS_AES_128_GCM_SHA256":                tls.TLS_AES_128_GCM_SHA256,
		"ECDHE-RSA-AES128-GCM-SHA256":           tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		"DES-CBC3-SHA":                          tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	}

	for name, expected := range expectedMappings {
		if cipher, ok := cipherMap[name]; !ok {
			t.Errorf("Expected cipher %s to be in map", name)
		} else if cipher != expected {
			t.Errorf("Expected cipher %s to map to %v, got %v", name, expected, cipher)
		}
	}

	// Verify map is not empty
	if len(cipherMap) == 0 {
		t.Error("Cipher suite map should not be empty")
	}
}

func TestGetTLSConfig_InvalidConfig(t *testing.T) {
	// Test with invalid config that will fail to create client
	invalidConfig := &rest.Config{
		Host:    "https://invalid-host.local:6443",
		Timeout: 1, // Very short timeout
	}

	// Should return error for connection failures
	_, err := GetTLSConfig(invalidConfig)
	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestGetTLSSecurityProfile_InvalidConfig(t *testing.T) {
	// Test error handling with invalid config
	invalidConfig := &rest.Config{
		Host:    "https://nonexistent.local:6443",
		Timeout: 1,
	}

	_, err := GetTLSSecurityProfile(invalidConfig)
	// Should return an error (not found or connection error)
	if err == nil {
		t.Error("Expected error for invalid config, got nil")
	}
}

func TestConvertTLSProfileToConfigAllProfiles(t *testing.T) {
	// Test all profile types to ensure full coverage
	profiles := []struct {
		name    string
		profile *configv1.TLSSecurityProfile
	}{
		{
			name:    "Old profile",
			profile: &configv1.TLSSecurityProfile{Type: configv1.TLSProfileOldType},
		},
		{
			name:    "Intermediate profile",
			profile: &configv1.TLSSecurityProfile{Type: configv1.TLSProfileIntermediateType},
		},
		{
			name:    "Modern profile",
			profile: &configv1.TLSSecurityProfile{Type: configv1.TLSProfileModernType},
		},
		{
			name: "Custom profile with all fields",
			profile: &configv1.TLSSecurityProfile{
				Type: configv1.TLSProfileCustomType,
				Custom: &configv1.CustomTLSProfile{
					TLSProfileSpec: configv1.TLSProfileSpec{
						MinTLSVersion: configv1.VersionTLS13,
						Ciphers: []string{
							"TLS_AES_128_GCM_SHA256",
							"TLS_AES_256_GCM_SHA384",
						},
					},
				},
			},
		},
	}

	for _, tt := range profiles {
		t.Run(tt.name, func(t *testing.T) {
			config := ConvertTLSProfileToConfig(tt.profile)
			if config == nil {
				t.Fatal("Expected non-nil config")
			}
			if config.MinVersion == 0 {
				t.Error("MinVersion should be set")
			}
		})
	}
}

func TestGetCipherSuitesEmptyList(t *testing.T) {
	// Test with empty cipher list
	ciphers := getCipherSuites([]string{})
	if len(ciphers) != 0 {
		t.Errorf("Expected empty slice for empty input, got %d ciphers", len(ciphers))
	}
}

func TestGetCipherSuitesAllInvalid(t *testing.T) {
	// Test with all invalid cipher names
	ciphers := getCipherSuites([]string{"INVALID1", "INVALID2", "INVALID3"})
	if len(ciphers) != 0 {
		t.Errorf("Expected 0 ciphers for all invalid names, got %d", len(ciphers))
	}
}

func TestGetTLSVersionAllVersions(t *testing.T) {
	// Test all specific TLS versions for complete coverage
	versions := map[string]uint16{
		string(configv1.VersionTLS10): tls.VersionTLS10,
		string(configv1.VersionTLS11): tls.VersionTLS11,
		string(configv1.VersionTLS12): tls.VersionTLS12,
		string(configv1.VersionTLS13): tls.VersionTLS13,
		"":                            tls.VersionTLS12, // empty defaults to 1.2
		"VersionTLS99":                tls.VersionTLS12, // unknown defaults to 1.2
	}

	for input, expected := range versions {
		result := getTLSVersion(input)
		if result != expected {
			t.Errorf("getTLSVersion(%q) = %v, expected %v", input, result, expected)
		}
	}
}

func TestGetTLSVersionNameAllVersions(t *testing.T) {
	// Test all version names for coverage
	versions := map[uint16]string{
		tls.VersionTLS10: "TLS 1.0",
		tls.VersionTLS11: "TLS 1.1",
		tls.VersionTLS12: "TLS 1.2",
		tls.VersionTLS13: "TLS 1.3",
		0x9999:           "Unknown (0x9999)",
		0x0000:           "Unknown (0x0000)",
	}

	for version, expected := range versions {
		result := getTLSVersionName(version)
		if result != expected {
			t.Errorf("getTLSVersionName(0x%04x) = %q, expected %q", version, result, expected)
		}
	}
}
