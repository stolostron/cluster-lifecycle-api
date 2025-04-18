// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	context "context"

	actionv1beta1 "github.com/stolostron/cluster-lifecycle-api/action/v1beta1"
	scheme "github.com/stolostron/cluster-lifecycle-api/client/action/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ManagedClusterActionsGetter has a method to return a ManagedClusterActionInterface.
// A group's client should implement this interface.
type ManagedClusterActionsGetter interface {
	ManagedClusterActions(namespace string) ManagedClusterActionInterface
}

// ManagedClusterActionInterface has methods to work with ManagedClusterAction resources.
type ManagedClusterActionInterface interface {
	Create(ctx context.Context, managedClusterAction *actionv1beta1.ManagedClusterAction, opts v1.CreateOptions) (*actionv1beta1.ManagedClusterAction, error)
	Update(ctx context.Context, managedClusterAction *actionv1beta1.ManagedClusterAction, opts v1.UpdateOptions) (*actionv1beta1.ManagedClusterAction, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, managedClusterAction *actionv1beta1.ManagedClusterAction, opts v1.UpdateOptions) (*actionv1beta1.ManagedClusterAction, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*actionv1beta1.ManagedClusterAction, error)
	List(ctx context.Context, opts v1.ListOptions) (*actionv1beta1.ManagedClusterActionList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *actionv1beta1.ManagedClusterAction, err error)
	ManagedClusterActionExpansion
}

// managedClusterActions implements ManagedClusterActionInterface
type managedClusterActions struct {
	*gentype.ClientWithList[*actionv1beta1.ManagedClusterAction, *actionv1beta1.ManagedClusterActionList]
}

// newManagedClusterActions returns a ManagedClusterActions
func newManagedClusterActions(c *ActionV1beta1Client, namespace string) *managedClusterActions {
	return &managedClusterActions{
		gentype.NewClientWithList[*actionv1beta1.ManagedClusterAction, *actionv1beta1.ManagedClusterActionList](
			"managedclusteractions",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *actionv1beta1.ManagedClusterAction { return &actionv1beta1.ManagedClusterAction{} },
			func() *actionv1beta1.ManagedClusterActionList { return &actionv1beta1.ManagedClusterActionList{} },
		),
	}
}
