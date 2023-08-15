package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=managedclusterinfos

// KlusterletVersion is the configuration for the klusterlet version and upgrade strategy of selected managedClusters.
// This is where parameters related to automatic updates can be set.
type KlusterletVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines rolling upgrade strategy of the Klusterlets for the managed clusters.
	// +optional
	Spec KlusterletVersionSpec `json:"spec,omitempty"`

	// Status represents the version status of the Klusterlets.
	// +optional
	Status KlusterletVersionStatus `json:"status,omitempty"`
}

type RollingType string

const (
	// Rolling is the type to rolling
	Rolling       RollingType = "Rolling"
	RollingCanary RollingType = "RollingCanary"
)

type KlusterletVersionSpec struct {
	// There are 2 types for klusterlet upgrade, Rolling and RollingCanary.
	// The default is type is Rolling.
	Type RollingType `json:"type,omitempty"`
	// RollingConfig is the config for Rolling upgrade.
	RollingConfig RollingConfig `json:"rollingConfig,omitempty"`
	// RollingCanaryConfig is the config for RollingCanary upgrade.
	RollingCanaryConfig RollingCanaryConfig `json:"rollingCanaryConfig,omitempty"`
}

type RollingConfig struct {
	// +optional
	MaxConcurrentlyUpdating *intstr.IntOrString `json:"maxConcurrentlyUpdating,omitempty"`
}
type RollingCanaryConfig struct {
	PlacementRef `json:",inline"`
	// +optional
	MaxConcurrentlyUpdating *intstr.IntOrString `json:"maxConcurrentlyUpdating,omitempty"`
}

type PlacementRef struct {
	// Name of the placement
	// +kubebuilder:validation:Required
	// +required
	Name string `json:"name"`

	// Namespace of the placement
	// +kubebuilder:validation:Required
	// +required
	Namespace string `json:"namespace"`
}

type KlusterletVersionStatus struct {
	// DesiredVersion is the MCE version.
	DesiredVersion string `json:"desiredVersion,omitempty"`

	// conditions describe the state of klusterlet upgrade.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KlusterletVersionList is a list of KlusterletVersionList objects.
type KlusterletVersionList struct {
	metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of KlusterletUpgradeStrategy objects.
	Items []KlusterletVersion `json:"items"`
}
