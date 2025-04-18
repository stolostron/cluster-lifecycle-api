// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	viewv1beta1 "github.com/stolostron/cluster-lifecycle-api/view/v1beta1"
	labels "k8s.io/apimachinery/pkg/labels"
	listers "k8s.io/client-go/listers"
	cache "k8s.io/client-go/tools/cache"
)

// ManagedClusterViewLister helps list ManagedClusterViews.
// All objects returned here must be treated as read-only.
type ManagedClusterViewLister interface {
	// List lists all ManagedClusterViews in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*viewv1beta1.ManagedClusterView, err error)
	// ManagedClusterViews returns an object that can list and get ManagedClusterViews.
	ManagedClusterViews(namespace string) ManagedClusterViewNamespaceLister
	ManagedClusterViewListerExpansion
}

// managedClusterViewLister implements the ManagedClusterViewLister interface.
type managedClusterViewLister struct {
	listers.ResourceIndexer[*viewv1beta1.ManagedClusterView]
}

// NewManagedClusterViewLister returns a new ManagedClusterViewLister.
func NewManagedClusterViewLister(indexer cache.Indexer) ManagedClusterViewLister {
	return &managedClusterViewLister{listers.New[*viewv1beta1.ManagedClusterView](indexer, viewv1beta1.Resource("managedclusterview"))}
}

// ManagedClusterViews returns an object that can list and get ManagedClusterViews.
func (s *managedClusterViewLister) ManagedClusterViews(namespace string) ManagedClusterViewNamespaceLister {
	return managedClusterViewNamespaceLister{listers.NewNamespaced[*viewv1beta1.ManagedClusterView](s.ResourceIndexer, namespace)}
}

// ManagedClusterViewNamespaceLister helps list and get ManagedClusterViews.
// All objects returned here must be treated as read-only.
type ManagedClusterViewNamespaceLister interface {
	// List lists all ManagedClusterViews in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*viewv1beta1.ManagedClusterView, err error)
	// Get retrieves the ManagedClusterView from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*viewv1beta1.ManagedClusterView, error)
	ManagedClusterViewNamespaceListerExpansion
}

// managedClusterViewNamespaceLister implements the ManagedClusterViewNamespaceLister
// interface.
type managedClusterViewNamespaceLister struct {
	listers.ResourceIndexer[*viewv1beta1.ManagedClusterView]
}
