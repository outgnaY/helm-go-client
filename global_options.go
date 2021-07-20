package helmclient

import "helm.sh/helm/v3/pkg/cli"

type GlobalOption struct {
	f func(o *globalOptions)
}

type globalOptions cli.EnvSettings

func (o *globalOptions) apply(opts []GlobalOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func addGlobalOptions(opts []GlobalOption, g *globalOptions) {
	g.apply(opts)
}

func WithKubeConfig(kubeConfig string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeConfig = kubeConfig
	}}
}

func WithKubeContext(kubeContext string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeContext = kubeContext
	}}
}

func WithKubeToken(kubeToken string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeToken = kubeToken
	}}
}

func WithKubeAsUser(kubeAsUser string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeAsUser = kubeAsUser
	}}
}

func WithKubeAsGroups(kubeAsGroups []string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeAsGroups = kubeAsGroups
	}}
}

func WithKubeAPIServer(kubeAPIServer string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeAPIServer = kubeAPIServer
	}}
}

func WithKubeCaFile(kubeCaFile string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.KubeCaFile = kubeCaFile
	}}
}

func WithDebug(debug bool) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.Debug = debug
	}}
}

func WithRegistryConfig(registryConfig string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.RegistryConfig = registryConfig
	}}
}

func WithRepositoryConfig(repositoryConfig string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.RepositoryConfig = repositoryConfig
	}}
}

func WithRepositoryCache(repositoryCache string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.RepositoryCache = repositoryCache
	}}
}

func WithPluginsDirectory(pluginsDirectory string) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.PluginsDirectory = pluginsDirectory
	}}
}

func WithMaxHistory(maxHistory int) GlobalOption {
	return GlobalOption{f: func(o *globalOptions) {
		o.MaxHistory = maxHistory
	}}
}
