// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/stolostron/cluster-lifecycle-api/client/view/clientset/versioned/typed/view/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeViewV1beta1 struct {
	*testing.Fake
}

func (c *FakeViewV1beta1) ManagedClusterViews(namespace string) v1beta1.ManagedClusterViewInterface {
	return newFakeManagedClusterViews(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeViewV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
