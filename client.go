package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"log"
	"os"
)

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	log.Output(2, fmt.Sprintf(format, v...))
}

func warning(format string, v ...interface{}) {
	format = fmt.Sprintf("WARNING: %s\n", format)
	fmt.Fprintf(os.Stderr, format, v...)
}

type HelmClient interface {
	List(opts []ListOption) (listClient, error)
	Uninstall(opts []UninstallOption) (uninstallClient, error)
	Rollback(opts []RollbackOption) (rollbackClient, error)
	GetValues(opts []GetValuesOption) (getValuesClient, error)
	History(opts []HistoryOption) (historyClient, error)
	GetAll(opts []GetAllOption) (getAllClient, error)
	GetManifest(opts []GetManifestOption) (getManifestClient, error)
	GetNotes(opts []GetNotesOption) (getNotesClient, error)
	GetHooks(opts []GetHooksOption) (getHooksClient, error)
	Lint(opts []LintOption, valueOpts []ValueOption) (lintClient, error)
	Install(opts []InstallOption, valueOpts []ValueOption, chartPathOpts []ChartPathOption) (installClient, error)
	Package(opts []PackageOption) (packageClient, error)
	Upgrade(opts []UpgradeOption, valueOpts []ValueOption, chartPathOpts []ChartPathOption) (upgradeClient, error)
	Create(opts []CreateOption) (createClient, error)
	Version(opts []VersionOption) (versionClient, error)
	ChartPull() (chartPullClient, error)
	ChartPush() (chartPushClient, error)
	ChartRemove() (chartRemoveClient, error)
	ChartSave() (chartSaveClient, error)
	ChartExport(opts []ChartExportOption) (chartExportClient, error)
	ChartList() (chartListClient, error)
	RegistryLogin(opts []RegistryLoginOption) (registryLoginClient, error)
	RegistryLogout() (registryLogoutClient, error)
	SearchHub(opts []SearchHubOption) (searchHubClient, error)
	SearchRepo(opts []SearchRepoOption) (searchRepoClient, error)
	RepoUpdate() (repoUpdateClient, error)
	RepoRemove() (repoRemoveClient, error)
	RepoList() (repoListClient, error)
	RepoAdd(opts []RepoAddOption) (repoAddClient, error)
	RepoIndex(opts []RepoIndexOption) (repoIndexClient, error)
	Pull(opts []PullOption, chartPathOpts []ChartPathOption) (pullClient, error)
	Template(opts []TemplateOption, valueOpts []ValueOption) (templateClient, error)
}

type helmEnv struct {
	settings     *cli.EnvSettings
	clientGetter *RESTClientGetter
}

func (env *helmEnv) Namespace() string {
	if ns, _, err := env.clientGetter.ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return "default"
}

type helmClientImpl struct {
	env *helmEnv
}

func (c *helmClientImpl) List(opts []ListOption) (listClient, error) {
	return newListClient(opts, c.env)
}

func (c *helmClientImpl) Uninstall(opts []UninstallOption) (uninstallClient, error) {
	return newUninstallClient(opts, c.env)
}

func (c *helmClientImpl) Rollback(opts []RollbackOption) (rollbackClient, error) {
	return newRollbackClient(opts, c.env)
}

func (c *helmClientImpl) GetValues(opts []GetValuesOption) (getValuesClient, error) {
	return newGetValuesClient(opts, c.env)
}

func (c *helmClientImpl) History(opts []HistoryOption) (historyClient, error) {
	return newHistoryClient(opts, c.env)
}

func (c *helmClientImpl) GetAll(opts []GetAllOption) (getAllClient, error) {
	return newGetAllClient(opts, c.env)
}

func (c *helmClientImpl) GetManifest(opts []GetManifestOption) (getManifestClient, error) {
	return newGetManifestClient(opts, c.env)
}

func (c *helmClientImpl) GetNotes(opts []GetNotesOption) (getNotesClient, error) {
	return newGetNotesClient(opts, c.env)
}

func (c *helmClientImpl) GetHooks(opts []GetHooksOption) (getHooksClient, error) {
	return newGetHooksClient(opts, c.env)
}

func (c *helmClientImpl) Lint(opts []LintOption, valueOpts []ValueOption) (lintClient, error) {
	return newLintClient(opts, valueOpts, c.env)
}

func (c *helmClientImpl) Install(opts []InstallOption, valueOpts []ValueOption, chartPathOpts []ChartPathOption) (installClient, error) {
	return newInstallClient(opts, valueOpts, chartPathOpts, c.env)
}

func (c *helmClientImpl) Package(opts []PackageOption) (packageClient, error) {
	return newPackageClient(opts, c.env)
}

func (c *helmClientImpl) Upgrade(opts []UpgradeOption, valueOpts []ValueOption, chartPathOpts []ChartPathOption) (upgradeClient, error) {
	return newUpgradeClient(opts, valueOpts, chartPathOpts, c.env)
}

func (c *helmClientImpl) Create(opts []CreateOption) (createClient, error) {
	return newCreateClient(opts, c.env)
}

func (c *helmClientImpl) Version(opts []VersionOption) (versionClient, error) {
	return newVersionClient(opts, c.env)
}

func (c *helmClientImpl) ChartPull() (chartPullClient, error) {
	return newChartPullClient(c.env)
}

func (c *helmClientImpl) ChartPush() (chartPushClient, error) {
	return newChartPushClient(c.env)
}

func (c *helmClientImpl) ChartRemove() (chartRemoveClient, error) {
	return newChartRemoveClient(c.env)
}

func (c *helmClientImpl) ChartSave() (chartSaveClient, error) {
	return newChartSaveClient(c.env)
}

func (c *helmClientImpl) ChartExport(opts []ChartExportOption) (chartExportClient, error) {
	return newChartExportClient(opts, c.env)
}

func (c *helmClientImpl) ChartList() (chartListClient, error) {
	return newChartListClient(c.env)
}

func (c *helmClientImpl) RegistryLogin(opts []RegistryLoginOption) (registryLoginClient, error) {
	return newRegistryLoginClient(opts, c.env)
}

func (c *helmClientImpl) RegistryLogout() (registryLogoutClient, error) {
	return newRegistryLogoutClient(c.env)
}

func (c *helmClientImpl) SearchHub(opts []SearchHubOption) (searchHubClient, error) {
	return newSearchHubClient(opts, c.env)
}

func (c *helmClientImpl) SearchRepo(opts []SearchRepoOption) (searchRepoClient, error) {
	return newSearchRepoClient(opts, c.env)
}

func (c *helmClientImpl) RepoUpdate() (repoUpdateClient, error) {
	return newRepoUpdateClient(c.env)
}

func (c *helmClientImpl) RepoRemove() (repoRemoveClient, error) {
	return newRepoRemoveClient(c.env)
}

func (c *helmClientImpl) RepoList() (repoListClient, error) {
	return newRepoListClient(c.env)
}

func (c *helmClientImpl) RepoAdd(opts []RepoAddOption) (repoAddClient, error) {
	return newRepoAddClient(opts, c.env)
}

func (c *helmClientImpl) RepoIndex(opts []RepoIndexOption) (repoIndexClient, error) {
	return newRepoIndexClient(opts, c.env)
}

func (c *helmClientImpl) Pull(opts []PullOption, chartPathOpts []ChartPathOption) (pullClient, error) {
	return newPullClient(opts, chartPathOpts, c.env)
}

func (c *helmClientImpl) Template(opts []TemplateOption, valueOpts []ValueOption) (templateClient, error) {
	return newTemplateClient(opts, valueOpts, c.env)
}

func NewHelmClient(kubeConfig string, namespace string) HelmClient {
	settings := cli.New()
	clientGetter := newRESTClientGetter(kubeConfig, namespace)
	env := &helmEnv{
		settings:     settings,
		clientGetter: clientGetter,
	}
	return &helmClientImpl{
		env: env,
	}
}

func NewHelmClientWithGlobalOpts(kubeConfig string, namespace string, globalOpts []GlobalOption) HelmClient {
	settings := cli.New()
	clientGetter := newRESTClientGetter(kubeConfig, namespace)
	addGlobalOptions(globalOpts, (*globalOptions)(settings))
	env := &helmEnv{
		settings:     settings,
		clientGetter: clientGetter,
	}
	return &helmClientImpl{
		env: env,
	}
}

func rebuildEnv(globalOpts []GlobalOption, namespace string, getter *RESTClientGetter) *helmEnv {
	settings := cli.New()
	addGlobalOptions(globalOpts, (*globalOptions)(settings))
	env := &helmEnv{
		settings:     settings,
		clientGetter: newRESTClientGetterFromOldWithNamespace(getter, namespace),
	}
	return env
}

func rebuildEnvAndCfg(globalOpts []GlobalOption, namespace string, getter *RESTClientGetter) (*helmEnv, *action.Configuration, error) {
	env := rebuildEnv(globalOpts, namespace, getter)
	cfg := new(action.Configuration)
	// must pass namespace explicitly cause cli.EnvSettings.namespace is private
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	return env, cfg, err
}
