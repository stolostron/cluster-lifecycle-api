// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/stolostron/cluster-lifecycle-api/client/imageregistry/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// ManagedClusterImageRegistries returns a ManagedClusterImageRegistryInformer.
	ManagedClusterImageRegistries() ManagedClusterImageRegistryInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// ManagedClusterImageRegistries returns a ManagedClusterImageRegistryInformer.
func (v *version) ManagedClusterImageRegistries() ManagedClusterImageRegistryInformer {
	return &managedClusterImageRegistryInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
