// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/stolostron/cluster-lifecycle-api/action/v1beta1"
	scheme "github.com/stolostron/cluster-lifecycle-api/client/action/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ManagedClusterActionsGetter has a method to return a ManagedClusterActionInterface.
// A group's client should implement this interface.
type ManagedClusterActionsGetter interface {
	ManagedClusterActions(namespace string) ManagedClusterActionInterface
}

// ManagedClusterActionInterface has methods to work with ManagedClusterAction resources.
type ManagedClusterActionInterface interface {
	Create(ctx context.Context, managedClusterAction *v1beta1.ManagedClusterAction, opts v1.CreateOptions) (*v1beta1.ManagedClusterAction, error)
	Update(ctx context.Context, managedClusterAction *v1beta1.ManagedClusterAction, opts v1.UpdateOptions) (*v1beta1.ManagedClusterAction, error)
	UpdateStatus(ctx context.Context, managedClusterAction *v1beta1.ManagedClusterAction, opts v1.UpdateOptions) (*v1beta1.ManagedClusterAction, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.ManagedClusterAction, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.ManagedClusterActionList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ManagedClusterAction, err error)
	ManagedClusterActionExpansion
}

// managedClusterActions implements ManagedClusterActionInterface
type managedClusterActions struct {
	client rest.Interface
	ns     string
}

// newManagedClusterActions returns a ManagedClusterActions
func newManagedClusterActions(c *ActionV1beta1Client, namespace string) *managedClusterActions {
	return &managedClusterActions{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the managedClusterAction, and returns the corresponding managedClusterAction object, and an error if there is any.
func (c *managedClusterActions) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.ManagedClusterAction, err error) {
	result = &v1beta1.ManagedClusterAction{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("managedclusteractions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ManagedClusterActions that match those selectors.
func (c *managedClusterActions) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.ManagedClusterActionList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.ManagedClusterActionList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("managedclusteractions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested managedClusterActions.
func (c *managedClusterActions) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("managedclusteractions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a managedClusterAction and creates it.  Returns the server's representation of the managedClusterAction, and an error, if there is any.
func (c *managedClusterActions) Create(ctx context.Context, managedClusterAction *v1beta1.ManagedClusterAction, opts v1.CreateOptions) (result *v1beta1.ManagedClusterAction, err error) {
	result = &v1beta1.ManagedClusterAction{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("managedclusteractions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(managedClusterAction).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a managedClusterAction and updates it. Returns the server's representation of the managedClusterAction, and an error, if there is any.
func (c *managedClusterActions) Update(ctx context.Context, managedClusterAction *v1beta1.ManagedClusterAction, opts v1.UpdateOptions) (result *v1beta1.ManagedClusterAction, err error) {
	result = &v1beta1.ManagedClusterAction{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("managedclusteractions").
		Name(managedClusterAction.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(managedClusterAction).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *managedClusterActions) UpdateStatus(ctx context.Context, managedClusterAction *v1beta1.ManagedClusterAction, opts v1.UpdateOptions) (result *v1beta1.ManagedClusterAction, err error) {
	result = &v1beta1.ManagedClusterAction{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("managedclusteractions").
		Name(managedClusterAction.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(managedClusterAction).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the managedClusterAction and deletes it. Returns an error if one occurs.
func (c *managedClusterActions) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("managedclusteractions").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *managedClusterActions) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("managedclusteractions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched managedClusterAction.
func (c *managedClusterActions) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.ManagedClusterAction, err error) {
	result = &v1beta1.ManagedClusterAction{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("managedclusteractions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
