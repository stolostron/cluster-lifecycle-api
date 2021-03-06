// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/stolostron/cluster-lifecycle-api/action/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ManagedClusterActionLister helps list ManagedClusterActions.
// All objects returned here must be treated as read-only.
type ManagedClusterActionLister interface {
	// List lists all ManagedClusterActions in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ManagedClusterAction, err error)
	// ManagedClusterActions returns an object that can list and get ManagedClusterActions.
	ManagedClusterActions(namespace string) ManagedClusterActionNamespaceLister
	ManagedClusterActionListerExpansion
}

// managedClusterActionLister implements the ManagedClusterActionLister interface.
type managedClusterActionLister struct {
	indexer cache.Indexer
}

// NewManagedClusterActionLister returns a new ManagedClusterActionLister.
func NewManagedClusterActionLister(indexer cache.Indexer) ManagedClusterActionLister {
	return &managedClusterActionLister{indexer: indexer}
}

// List lists all ManagedClusterActions in the indexer.
func (s *managedClusterActionLister) List(selector labels.Selector) (ret []*v1beta1.ManagedClusterAction, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ManagedClusterAction))
	})
	return ret, err
}

// ManagedClusterActions returns an object that can list and get ManagedClusterActions.
func (s *managedClusterActionLister) ManagedClusterActions(namespace string) ManagedClusterActionNamespaceLister {
	return managedClusterActionNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ManagedClusterActionNamespaceLister helps list and get ManagedClusterActions.
// All objects returned here must be treated as read-only.
type ManagedClusterActionNamespaceLister interface {
	// List lists all ManagedClusterActions in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta1.ManagedClusterAction, err error)
	// Get retrieves the ManagedClusterAction from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta1.ManagedClusterAction, error)
	ManagedClusterActionNamespaceListerExpansion
}

// managedClusterActionNamespaceLister implements the ManagedClusterActionNamespaceLister
// interface.
type managedClusterActionNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ManagedClusterActions in the indexer for a given namespace.
func (s managedClusterActionNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.ManagedClusterAction, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.ManagedClusterAction))
	})
	return ret, err
}

// Get retrieves the ManagedClusterAction from the indexer for a given namespace and name.
func (s managedClusterActionNamespaceLister) Get(name string) (*v1beta1.ManagedClusterAction, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("managedclusteraction"), name)
	}
	return obj.(*v1beta1.ManagedClusterAction), nil
}
