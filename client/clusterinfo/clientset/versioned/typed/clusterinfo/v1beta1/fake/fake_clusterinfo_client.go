// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "github.com/stolostron/cluster-lifecycle-api/client/clusterinfo/clientset/versioned/typed/clusterinfo/v1beta1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeInternalV1beta1 struct {
	*testing.Fake
}

func (c *FakeInternalV1beta1) ManagedClusterInfos(namespace string) v1beta1.ManagedClusterInfoInterface {
	return newFakeManagedClusterInfos(c, namespace)
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeInternalV1beta1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
