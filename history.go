package helmclient

import (
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	helmtime "helm.sh/helm/v3/pkg/time"
)

const (
	historyDefaultMax = 256
)

type historyClient interface {
	History(name string) (ReleaseHistory, error)
}

type historyClientImpl struct {
	cli *action.History
	env *helmEnv
}

type HistoryOption struct {
	f func(o *historyOptions)
}

type historyOptions struct {
	max int
}

func (o *historyOptions) apply(opts []HistoryOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newHistoryOptions(opts []HistoryOption) *historyOptions {
	options := &historyOptions{max: historyDefaultMax}
	options.apply(opts)
	return options
}

func HistoryWithMax(max int) HistoryOption {
	return HistoryOption{f: func(o *historyOptions) {
		o.max = max
	}}
}

func (c *historyClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *historyClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewHistory(cfg)
	// copy args
	copyHistoryClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func copyHistoryClientOptions(oldCli *action.History, newCli *action.History) {
	newCli.Max = oldCli.Max
}

func newHistoryClient(opts []HistoryOption, env *helmEnv) (*historyClientImpl, error) {
	o := newHistoryOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewHistory(cfg)
	mergeHistoryOptions(o, client)
	return &historyClientImpl{
		cli: client,
		env: env,
	}, nil
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func formatChartname(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return fmt.Sprintf("%s-%s", c.Name(), c.Metadata.Version)
}

func formatAppVersion(c *chart.Chart) string {
	if c == nil || c.Metadata == nil {
		// This is an edge case that has happened in prod, though we don't
		// know how: https://github.com/helm/helm/issues/1347
		return "MISSING"
	}
	return c.AppVersion()
}

type ReleaseInfo struct {
	Revision    int           `json:"revision"`
	Updated     helmtime.Time `json:"updated"`
	Status      string        `json:"status"`
	Chart       string        `json:"chart"`
	AppVersion  string        `json:"app_version"`
	Description string        `json:"description"`
}

type ReleaseHistory []ReleaseInfo

func getHistory(client *action.History, name string) (ReleaseHistory, error) {
	hist, err := client.Run(name)
	if err != nil {
		return nil, err
	}

	releaseutil.Reverse(hist, releaseutil.SortByRevision)

	var rels []*release.Release
	for i := 0; i < min(len(hist), client.Max); i++ {
		rels = append(rels, hist[i])
	}

	if len(rels) == 0 {
		return ReleaseHistory{}, nil
	}

	releaseHistory := getReleaseHistory(rels)

	return releaseHistory, nil
}

func getReleaseHistory(rls []*release.Release) (history ReleaseHistory) {
	for i := len(rls) - 1; i >= 0; i-- {
		r := rls[i]
		c := formatChartname(r.Chart)
		s := r.Info.Status.String()
		v := r.Version
		d := r.Info.Description
		a := formatAppVersion(r.Chart)

		rInfo := ReleaseInfo{
			Revision:    v,
			Status:      s,
			Chart:       c,
			AppVersion:  a,
			Description: d,
		}
		if !r.Info.LastDeployed.IsZero() {
			rInfo.Updated = r.Info.LastDeployed

		}
		history = append(history, rInfo)
	}
	return history
}

func (c *historyClientImpl) History(name string) (ReleaseHistory, error) {
	return getHistory(c.cli, name)
}

func mergeHistoryOptions(o *historyOptions, cli *action.History) {
	cli.Max = o.max
}
