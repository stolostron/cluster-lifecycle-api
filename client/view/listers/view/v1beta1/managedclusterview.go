// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/stolostron/cluster-lifecycle-api/view/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ManagedClusterViewLister helps list ManagedClusterViews.
// All objects returned here must be treated as read-only.
type ManagedClusterViewLister interface {
	// List lists all ManagedClusterViews in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ManagedClusterView, err error)
	// ManagedClusterViews returns an object that can list and get ManagedClusterViews.
	ManagedClusterViews(namespace string) ManagedClusterViewNamespaceLister
	ManagedClusterViewListerExpansion
}

// managedClusterViewLister implements the ManagedClusterViewLister interface.
type managedClusterViewLister struct {
	indexer cache.Indexer
}

// NewManagedClusterViewLister returns a new ManagedClusterViewLister.
func NewManagedClusterViewLister(indexer cache.Indexer) ManagedClusterViewLister {
	return &managedClusterViewLister{indexer: indexer}
}

// List lists all ManagedClusterViews in the indexer.
func (s *managedClusterViewLister) List(selector labels.Selector) (ret []*v1beta1.ManagedClusterView, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ManagedClusterView))
	})
	return ret, err
}

// ManagedClusterViews returns an object that can list and get ManagedClusterViews.
func (s *managedClusterViewLister) ManagedClusterViews(namespace string) ManagedClusterViewNamespaceLister {
	return managedClusterViewNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ManagedClusterViewNamespaceLister helps list and get ManagedClusterViews.
// All objects returned here must be treated as read-only.
type ManagedClusterViewNamespaceLister interface {
	// List lists all ManagedClusterViews in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ManagedClusterView, err error)
	// Get retrieves the ManagedClusterView from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.ManagedClusterView, error)
	ManagedClusterViewNamespaceListerExpansion
}

// managedClusterViewNamespaceLister implements the ManagedClusterViewNamespaceLister
// interface.
type managedClusterViewNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ManagedClusterViews in the indexer for a given namespace.
func (s managedClusterViewNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.ManagedClusterView, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ManagedClusterView))
	})
	return ret, err
}

// Get retrieves the ManagedClusterView from the indexer for a given namespace and name.
func (s managedClusterViewNamespaceLister) Get(name string) (*v1beta1.ManagedClusterView, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("managedclusterview"), name)
	}
	return obj.(*v1beta1.ManagedClusterView), nil
}
