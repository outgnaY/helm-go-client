package helmclient

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type RESTClientGetter struct {
	kubeConfig string
	namespace  string
}

func newRESTClientGetterFromOldWithNamespace(old *RESTClientGetter, namespace string) *RESTClientGetter {
	return &RESTClientGetter{
		kubeConfig: old.kubeConfig,
		namespace:  namespace,
	}
}

func newRESTClientGetter(kubeConfig string, namespace string) *RESTClientGetter {
	return &RESTClientGetter{
		kubeConfig: kubeConfig,
		namespace:  namespace,
	}
}

func (c *RESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return clientcmd.RESTConfigFromKubeConfig([]byte(c.kubeConfig))
}

func (c *RESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}
	// for discovery
	config.Burst = 100
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, err
	}
	return memory.NewMemCacheClient(discoveryClient), nil
}

func (c *RESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := c.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func (c *RESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	overrides.Context.Namespace = c.namespace
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
}
