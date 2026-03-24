// Copyright Contributors to the Open Cluster Management project

package tlsprofile

import (
	"context"
	"reflect"
	"time"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// StartTLSProfileWatcher watches the APIServer CR for TLS profile changes and triggers a graceful restart.
// The cancel function will be called when a TLS profile change is detected, triggering graceful shutdown.
// This allows Kubernetes to restart the pod with the new TLS settings.
func StartTLSProfileWatcher(ctx context.Context, kubeConfig *rest.Config, cancel context.CancelFunc) error {
	// Get initial TLS profile to compare against
	initialProfile, err := GetTLSSecurityProfile(kubeConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.V(4).Info("APIServer CR not found, not running on OpenShift - TLS profile watcher disabled")
			return nil
		}
		klog.V(4).Infof("Failed to get initial TLS profile, watcher will not start: %v", err)
		return nil // Don't fail the application if watcher can't start
	}

	klog.Infof("Starting TLS profile watcher with initial profile: type=%s", initialProfile.Type)

	// Create OpenShift config client
	client, err := configclient.NewForConfig(kubeConfig)
	if err != nil {
		klog.V(4).Infof("Failed to create config client for watcher: %v", err)
		return nil
	}

	// Create ListWatch for APIServer resource, filtered to only "cluster"
	listWatch := cache.NewListWatchFromClient(
		client.ConfigV1().RESTClient(),
		"apiservers",
		"",
		fields.OneTermEqualSelector("metadata.name", "cluster"),
	)

	// Create informer from ListWatch
	_, controller := cache.NewInformer(
		listWatch,
		&configv1.APIServer{},
		10*time.Minute,
		cache.ResourceEventHandlerFuncs{
			UpdateFunc: func(oldObj, newObj interface{}) {
				handleAPIServerUpdate(oldObj, newObj, cancel)
			},
		},
	)

	// Start the controller
	go controller.Run(ctx.Done())

	klog.Info("TLS profile watcher started successfully")
	return nil
}

// tlsProfileChanged checks if two TLS profiles are different.
func tlsProfileChanged(old, new *configv1.TLSSecurityProfile) bool {
	// If both are nil, no change
	if old == nil && new == nil {
		return false
	}

	// If one is nil and the other isn't, it's a change
	if (old == nil) != (new == nil) {
		return true
	}

	// Compare profile types
	if old.Type != new.Type {
		return true
	}

	// For custom profiles, deep compare the custom settings
	if old.Type == configv1.TLSProfileCustomType {
		return !reflect.DeepEqual(old.Custom, new.Custom)
	}

	// For predefined profiles, type comparison is sufficient
	return false
}

// handleAPIServerUpdate handles APIServer update events.
func handleAPIServerUpdate(oldObj, newObj interface{}, cancel context.CancelFunc) {
	oldAPI, ok1 := oldObj.(*configv1.APIServer)
	newAPI, ok2 := newObj.(*configv1.APIServer)
	if !ok1 || !ok2 {
		return
	}

	// Only care about the cluster APIServer
	if newAPI.Name != "cluster" {
		return
	}

	// Check if TLS profile changed
	if !tlsProfileChanged(oldAPI.Spec.TLSSecurityProfile, newAPI.Spec.TLSSecurityProfile) {
		return
	}

	klog.Infof("TLS profile changed from %v to %v - triggering graceful shutdown",
		getTLSProfileType(oldAPI.Spec.TLSSecurityProfile),
		getTLSProfileType(newAPI.Spec.TLSSecurityProfile))

	// Trigger graceful shutdown by cancelling the context
	// This allows HTTP servers and other components to shut down cleanly
	// Kubernetes will then automatically restart the pod
	klog.Info("Cancelling context for graceful shutdown due to TLS profile change...")
	cancel()
}

// getTLSProfileType returns a string representation of the profile type.
func getTLSProfileType(profile *configv1.TLSSecurityProfile) string {
	if profile == nil {
		return "nil (default Intermediate)"
	}
	return string(profile.Type)
}
