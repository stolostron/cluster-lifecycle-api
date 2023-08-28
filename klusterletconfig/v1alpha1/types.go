package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	operatorv1 "open-cluster-management.io/api/operator/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=klusterletconfigs
// +kubebuilder:resource:scope=Cluster

// KlusterletConfig contains the configuration of a klusterlet including the upgrade strategy, config overrides, proxy configurations etc.
type KlusterletConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of KlusterletConfig
	// +optional
	Spec KlusterletConfigSpec `json:"spec,omitempty"`

	// Status defines the observed state of KlusterletConfig
	// +optional
	Status KlusterletConfigStatus `json:"status,omitempty"`
}

// KlusterletConfigSpec defines the desired state of KlusterletConfig, usually provided by the user.
type KlusterletConfigSpec struct {
	// Registries includes the mirror and source registries. The source registry will be replaced by the Mirror.
	// +optional
	Registries []Registries `json:"registries,omitempty"`

	// PullSecret is the name of image pull secret.
	// +optional
	PullSecret corev1.ObjectReference `json:"pullSecret,omitempty"`

	// NodePlacement enables explicit control over the scheduling of the agent components.
	// If the placement is nil, the placement is not specified, it will be omitted.
	// If the placement is an empty object, the placement will match all nodes and tolerate nothing.
	// +optional
	NodePlacement *operatorv1.NodePlacement `json:"nodePlacement,omitempty"`
}

// KlusterletConfigStatus defines the observed state of KlusterletConfig.
type KlusterletConfigStatus struct {
}

type Registries struct {
	// Mirror is the mirrored registry of the Source. Will be ignored if Mirror is empty.
	// +kubebuilder:validation:Required
	// +required
	Mirror string `json:"mirror"`

	// Source is the source registry. All image registries will be replaced by Mirror if Source is empty.
	// +optional
	Source string `json:"source"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KlusterletConfigList contains a list of KlusterletConfig.
type KlusterletConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KlusterletConfig `json:"items"`
}
