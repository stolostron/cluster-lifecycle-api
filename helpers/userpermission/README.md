# UserPermission Helper

This package provides helper functions for querying and consolidating user permissions across multiple clusters.

## GetSelfPermissionRules

`GetSelfPermissionRules` retrieves and consolidates permission rules for the authenticated user across multiple clusters.

### Function Signature

```go
func GetSelfPermissionRules(ctx context.Context, config *rest.Config, interestedVerb ...string) ([]PermissionRule, error)
```

### Parameters

- `ctx`: Context for the API request
- `config`: Kubernetes REST config for the authenticated user whose permissions should be evaluated. **This must be the user's config, not a service account or admin config.**
- `interestedVerb`: Optional list of verbs to filter by (e.g., "get", "list", "create"). If provided, only rules containing ALL specified verbs will be returned. Rules with wildcard verbs ("*") will always be included regardless of this filter.

### Returns

- `[]PermissionRule`: A consolidated list of permission rules with the following structure:
  ```go
  type PermissionRule struct {
      authzv1.ResourceRule           // Contains Verbs, APIGroups, Resources
      Clusters   []string `json:"clusters"`    // List of cluster names where this rule applies
      Namespaces []string `json:"namespaces"`  // List of namespaces where this rule applies
  }
  ```
- `error`: Any error encountered during the retrieval or processing

### How It Works

The function:
1. Queries all UserPermission resources accessible by the authenticated user
2. Consolidates permissions across clusters and namespaces:
   - Groups identical resource rules across multiple clusters into a single PermissionRule entry
   - Merges namespace lists when the same resource rule exists within a single cluster
   - Handles admin permissions (`Verbs:*, APIGroups:*, Resources:*`) specially by overriding other permissions for that cluster
3. Optionally filters rules by interested verbs

## Usage Examples

### Example 1: Service Configuration for Regular User

A service needs to check what resources a regular user can access across clusters.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/stolostron/cluster-lifecycle-api/helpers/userpermission"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    ctx := context.Background()

    // Load user's kubeconfig - this is the authenticated user's config
    // NOT a service account or admin config
    userConfig, err := clientcmd.BuildConfigFromFlags("", "/path/to/user/kubeconfig")
    if err != nil {
        log.Fatal(err)
    }

    // Get all permissions for the user
    permissions, err := userpermission.GetSelfPermissionRules(ctx, userConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Display the consolidated permissions
    for _, perm := range permissions {
        fmt.Printf("Verbs: %v\n", perm.Verbs)
        fmt.Printf("APIGroups: %v\n", perm.APIGroups)
        fmt.Printf("Resources: %v\n", perm.Resources)
        fmt.Printf("Clusters: %v\n", perm.Clusters)
        fmt.Printf("Namespaces: %v\n\n", perm.Namespaces)
    }
}
```

**Expected Output for a Developer User:**
```
Verbs: [get list watch]
APIGroups: []
Resources: [pods]
Clusters: [dev-cluster-1 dev-cluster-2]
Namespaces: [default my-app]

Verbs: [get list watch create update delete]
APIGroups: [apps]
Resources: [deployments]
Clusters: [dev-cluster-1]
Namespaces: [my-app]
```

### Example 2: Service Configuration for Admin User

When an admin user authenticates, the service should see admin-level permissions.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/stolostron/cluster-lifecycle-api/helpers/userpermission"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    ctx := context.Background()

    // Load admin user's kubeconfig
    adminConfig, err := clientcmd.BuildConfigFromFlags("", "/path/to/admin/kubeconfig")
    if err != nil {
        log.Fatal(err)
    }

    // Get all permissions for the admin user
    permissions, err := userpermission.GetSelfPermissionRules(ctx, adminConfig)
    if err != nil {
        log.Fatal(err)
    }

    // Admin users typically have wildcard permissions
    for _, perm := range permissions {
        fmt.Printf("Verbs: %v\n", perm.Verbs)
        fmt.Printf("APIGroups: %v\n", perm.APIGroups)
        fmt.Printf("Resources: %v\n", perm.Resources)
        fmt.Printf("Clusters: %v\n", perm.Clusters)
        fmt.Printf("Namespaces: %v\n\n", perm.Namespaces)
    }
}
```

**Expected Output for an Admin User:**
```
Verbs: [*]
APIGroups: [*]
Resources: [*]
Clusters: [prod-cluster-1 prod-cluster-2]
Namespaces: [*]
```

### Example 3: Filtering by Interested Verbs (Read-Only Access)

A service wants to check if a user has read-only access to specific resources.

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/stolostron/cluster-lifecycle-api/helpers/userpermission"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    ctx := context.Background()

    // Load user's kubeconfig
    userConfig, err := clientcmd.BuildConfigFromFlags("", "/path/to/user/kubeconfig")
    if err != nil {
        log.Fatal(err)
    }

    // Get only permissions that include "get" and "list" verbs
    permissions, err := userpermission.GetSelfPermissionRules(ctx, userConfig, "get", "list")
    if err != nil {
        log.Fatal(err)
    }

    // Only rules containing both "get" AND "list" verbs will be returned
    for _, perm := range permissions {
        fmt.Printf("Resource: %v.%v\n", perm.Resources, perm.APIGroups)
        fmt.Printf("Available Verbs: %v\n", perm.Verbs)
        fmt.Printf("Accessible Clusters: %v\n", perm.Clusters)
        fmt.Printf("Accessible Namespaces: %v\n\n", perm.Namespaces)
    }
}
```

**Expected Output:**
```
Resource: [pods].[]
Available Verbs: [get list]
Accessible Clusters: [dev-cluster-1 dev-cluster-2]
Accessible Namespaces: [default my-app]

Resource: [deployments].[apps]
Available Verbs: [get list]
Accessible Clusters: [dev-cluster-1]
Accessible Namespaces: [my-app]
```

### Example 4: Web Service with Multiple User Tokens

A web service receives user tokens and needs to determine what each user can access.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/stolostron/cluster-lifecycle-api/helpers/userpermission"
    "k8s.io/client-go/rest"
)

func getUserPermissions(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Extract the user's token from the request
    // This could come from an Authorization header, cookie, etc.
    userToken := r.Header.Get("Authorization")

    // Create a REST config using the user's token
    userConfig := &rest.Config{
        Host:        "https://your-kubernetes-api-server",
        BearerToken: userToken,
        TLSClientConfig: rest.TLSClientConfig{
            Insecure: false, // Set to true for development only
            CAFile:   "/path/to/ca.crt",
        },
    }

    // Get permissions for this specific user
    permissions, err := userpermission.GetSelfPermissionRules(ctx, userConfig)
    if err != nil {
        http.Error(w, "Failed to get user permissions", http.StatusInternalServerError)
        log.Printf("Error getting permissions: %v", err)
        return
    }

    // Use permissions to determine what the user can access
    // For example, filter available clusters based on permissions
    accessibleClusters := make(map[string]bool)
    for _, perm := range permissions {
        for _, cluster := range perm.Clusters {
            accessibleClusters[cluster] = true
        }
    }

    fmt.Fprintf(w, "User has access to clusters: %v\n", getKeys(accessibleClusters))
}

func getKeys(m map[string]bool) []string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

## Configuration Requirements

### For Service Providers

When integrating this function into your service, ensure:

1. **User Authentication**: The `rest.Config` must be configured with the actual user's credentials (token, certificate, etc.), not a service account or admin config.

2. **Token Sources**: Common ways to obtain user configs:
   - From a kubeconfig file: Use `clientcmd.BuildConfigFromFlags()`
   - From a bearer token: Create a `rest.Config` with the `BearerToken` field set
   - From client certificates: Set `TLSClientConfig.CertFile` and `KeyFile`

3. **Multiple Users**: If your service handles multiple users:
   - Create a separate `rest.Config` for each user
   - DO NOT reuse configs between users
   - Cache results per user session if needed for performance

### Security Considerations

- Never use service account tokens to query user permissions
- Always validate that the config belongs to the authenticated user
- Consider caching results per user session to reduce API calls
- Implement proper error handling for unauthorized or forbidden responses

## Permission Consolidation Behavior

### Consolidation Rules

1. **Across Clusters**: If a user has the same permission on multiple clusters, they are consolidated into a single `PermissionRule` with multiple clusters.

2. **Within Clusters**: If a user has the same permission in multiple namespaces within a cluster, the namespaces are merged.

3. **Admin Override**: If a user has admin permissions (`Verbs:*, APIGroups:*, Resources:*`) on a cluster, all other permissions for that cluster are ignored.

### Examples

**Before Consolidation** (internal representation):
```
- Cluster: cluster1, Namespace: default, Verbs: [get, list], Resources: [pods]
- Cluster: cluster2, Namespace: default, Verbs: [get, list], Resources: [pods]
- Cluster: cluster1, Namespace: kube-system, Verbs: [get, list], Resources: [pods]
```

**After Consolidation**:
```
- Clusters: [cluster1, cluster2], Namespaces: [default, kube-system], Verbs: [get, list], Resources: [pods]
```

## Error Handling

```go
permissions, err := userpermission.GetSelfPermissionRules(ctx, userConfig)
if err != nil {
    // Common errors:
    // - User is not authenticated (401)
    // - User does not have permission to list UserPermissions (403)
    // - Network or API server errors
    log.Printf("Failed to get permissions: %v", err)
    return
}
```

## Related Resources

- UserPermission CRD documentation
- Kubernetes RBAC documentation
- ClusterView API reference
