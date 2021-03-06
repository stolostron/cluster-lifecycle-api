// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/stolostron/cluster-lifecycle-api/inventory/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// BareMetalAssetLister helps list BareMetalAssets.
// All objects returned here must be treated as read-only.
type BareMetalAssetLister interface {
	// List lists all BareMetalAssets in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.BareMetalAsset, err error)
	// BareMetalAssets returns an object that can list and get BareMetalAssets.
	BareMetalAssets(namespace string) BareMetalAssetNamespaceLister
	BareMetalAssetListerExpansion
}

// bareMetalAssetLister implements the BareMetalAssetLister interface.
type bareMetalAssetLister struct {
	indexer cache.Indexer
}

// NewBareMetalAssetLister returns a new BareMetalAssetLister.
func NewBareMetalAssetLister(indexer cache.Indexer) BareMetalAssetLister {
	return &bareMetalAssetLister{indexer: indexer}
}

// List lists all BareMetalAssets in the indexer.
func (s *bareMetalAssetLister) List(selector labels.Selector) (ret []*v1alpha1.BareMetalAsset, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.BareMetalAsset))
	})
	return ret, err
}

// BareMetalAssets returns an object that can list and get BareMetalAssets.
func (s *bareMetalAssetLister) BareMetalAssets(namespace string) BareMetalAssetNamespaceLister {
	return bareMetalAssetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// BareMetalAssetNamespaceLister helps list and get BareMetalAssets.
// All objects returned here must be treated as read-only.
type BareMetalAssetNamespaceLister interface {
	// List lists all BareMetalAssets in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.BareMetalAsset, err error)
	// Get retrieves the BareMetalAsset from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.BareMetalAsset, error)
	BareMetalAssetNamespaceListerExpansion
}

// bareMetalAssetNamespaceLister implements the BareMetalAssetNamespaceLister
// interface.
type bareMetalAssetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all BareMetalAssets in the indexer for a given namespace.
func (s bareMetalAssetNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.BareMetalAsset, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.BareMetalAsset))
	})
	return ret, err
}

// Get retrieves the BareMetalAsset from the indexer for a given namespace and name.
func (s bareMetalAssetNamespaceLister) Get(name string) (*v1alpha1.BareMetalAsset, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("baremetalasset"), name)
	}
	return obj.(*v1alpha1.BareMetalAsset), nil
}
