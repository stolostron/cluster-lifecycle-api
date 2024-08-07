package v1alpha1

// This file contains a collection of methods that can be used from go-restful to
// generate Swagger API documentation for its models. Please read this PR for more
// information on the implementation: https://github.com/emicklei/go-restful/pull/215
//
// TODOs are ignored from the parser (e.g. TODO(andronat):... || TODO:...) if and only if
// they are on one line! For multiple line or blocks that you want to ignore use ---.
// Any context after a --- is ignored.
//
// Those methods can be generated by using hack/update-swagger-docs.sh

// AUTO-GENERATED FUNCTIONS START HERE
var map_InstallMode = map[string]string{
	"type":       "InstallModeType is the type of install mode.",
	"noOperator": "NoOperator is the setting of klusterlet installation when install type is noOperator.",
}

func (InstallMode) SwaggerDoc() map[string]string {
	return map_InstallMode
}

var map_KlusterletConfig = map[string]string{
	"":       "KlusterletConfig contains the configuration of a klusterlet including the upgrade strategy, config overrides, proxy configurations etc.",
	"spec":   "Spec defines the desired state of KlusterletConfig",
	"status": "Status defines the observed state of KlusterletConfig",
}

func (KlusterletConfig) SwaggerDoc() map[string]string {
	return map_KlusterletConfig
}

var map_KlusterletConfigList = map[string]string{
	"": "KlusterletConfigList contains a list of KlusterletConfig.",
}

func (KlusterletConfigList) SwaggerDoc() map[string]string {
	return map_KlusterletConfigList
}

var map_KlusterletConfigSpec = map[string]string{
	"":                                       "KlusterletConfigSpec defines the desired state of KlusterletConfig, usually provided by the user.",
	"registries":                             "Registries includes the mirror and source registries. The source registry will be replaced by the Mirror.",
	"pullSecret":                             "PullSecret is the name of image pull secret.",
	"nodePlacement":                          "NodePlacement enables explicit control over the scheduling of the agent components. If the placement is nil, the placement is not specified, it will be omitted. If the placement is an empty object, the placement will match all nodes and tolerate nothing.",
	"hubKubeAPIServerProxyConfig":            "HubKubeAPIServerProxyConfig holds proxy settings for connections between klusterlet/add-on agents on the managed cluster and the kube-apiserver on the hub cluster. Empty means no proxy settings is available.",
	"hubKubeAPIServerURL":                    "HubKubeAPIServerURL is the URL of the hub Kube API server. If not present, the .status.apiServerURL of Infrastructure/cluster will be used as the default value. e.g. `oc get infrastructure cluster -o jsonpath='{.status.apiServerURL}'`",
	"hubKubeAPIServerCABundle":               "HubKubeAPIServerCABundle is the CA bundle to verify the server certificate of the hub kube API against. If not present, CA bundle will be determined with the logic below: 1). Use the certificate of the named certificate configured in APIServer/cluster if FQDN matches; 2). Otherwise use the CA certificates from kube-root-ca.crt ConfigMap in the cluster namespace;",
	"appliedManifestWorkEvictionGracePeriod": "AppliedManifestWorkEvictionGracePeriod is the eviction grace period the work agent will wait before evicting the AppliedManifestWorks, whose corresponding ManifestWorks are missing on the hub cluster, from the managed cluster. If not present, the default value of the work agent will be used. If its value is set to \"INFINITE\", it means the AppliedManifestWorks will never been evicted from the managed cluster.",
	"installMode":                            "InstallMode is the mode to install the klusterlet",
	"bootstrapKubeConfigs":                   "BootstrapKubeConfigSecrets is the list of secrets that reflects the Klusterlet.Spec.RegistrationConfiguration.BootstrapKubeConfigs.",
}

func (KlusterletConfigSpec) SwaggerDoc() map[string]string {
	return map_KlusterletConfigSpec
}

var map_KlusterletConfigStatus = map[string]string{
	"": "KlusterletConfigStatus defines the observed state of KlusterletConfig.",
}

func (KlusterletConfigStatus) SwaggerDoc() map[string]string {
	return map_KlusterletConfigStatus
}

var map_KubeAPIServerProxyConfig = map[string]string{
	"":           "KubeAPIServerProxyConfig describes the proxy settings for the connections to a kube-apiserver",
	"httpProxy":  "HTTPProxy is the URL of the proxy for HTTP requests",
	"httpsProxy": "HTTPSProxy is the URL of the proxy for HTTPS requests HTTPSProxy will be chosen if both HTTPProxy and HTTPSProxy are set.",
	"caBundle":   "CABundle is a CA certificate bundle to verify the proxy server. It will be ignored if only HTTPProxy is set; And it is required when HTTPSProxy is set and self signed CA certificate is used by the proxy server.",
}

func (KubeAPIServerProxyConfig) SwaggerDoc() map[string]string {
	return map_KubeAPIServerProxyConfig
}

var map_NoOperator = map[string]string{
	"postfix": "Postfix is the postfix of the klusterlet name. The name of the klusterlet is \"klusterlet\" if it is not set, and \"klusterlet-{Postfix}\". The install namespace is \"open-cluster-management-agent\" if it is not set, and \"open-cluster-management-{Postfix}\".",
}

func (NoOperator) SwaggerDoc() map[string]string {
	return map_NoOperator
}

var map_Registries = map[string]string{
	"mirror": "Mirror is the mirrored registry of the Source. Will be ignored if Mirror is empty.",
	"source": "Source is the source registry. All image registries will be replaced by Mirror if Source is empty.",
}

func (Registries) SwaggerDoc() map[string]string {
	return map_Registries
}

// AUTO-GENERATED FUNCTIONS END HERE
