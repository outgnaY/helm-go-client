package helmclient

import (
	"github.com/outgnaY/helm-go-client/internal/monocular"
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const (
	searchHubDefaultSearchEndPoint = "https://hub.helm.sh"
)

type searchHubClient interface {
	SearchHub(args []string) ([]monocular.SearchResult, error)
}

type searchHubClientImpl struct {
	searchHubOpts *searchHubOptions
	env           *helmEnv
}

type SearchHubOption struct {
	f func(o *searchHubOptions)
}

type searchHubOptions struct {
	searchEndpoint string
}

func (o *searchHubOptions) apply(opts []SearchHubOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newSearchHubOptions(opts []SearchHubOption) *searchHubOptions {
	options := &searchHubOptions{
		searchEndpoint: searchHubDefaultSearchEndPoint,
	}
	options.apply(opts)
	return options
}

func SearchHubWithSearchEndpoint(searchEndpoint string) SearchHubOption {
	return SearchHubOption{f: func(o *searchHubOptions) {
		o.searchEndpoint = searchEndpoint
	}}
}

func (c *searchHubClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *searchHubClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newSearchHubClient(opts []SearchHubOption, env *helmEnv) (*searchHubClientImpl, error) {
	o := newSearchHubOptions(opts)
	return &searchHubClientImpl{
		searchHubOpts: o,
		env:           env,
	}, nil
}

func (o *searchHubOptions) run(args []string) ([]monocular.SearchResult, error) {
	c, err := monocular.New(o.searchEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to create connection to %q", o.searchEndpoint))
	}
	q := strings.Join(args, " ")
	return c.Search(q)
}

func (c *searchHubClientImpl) SearchHub(args []string) ([]monocular.SearchResult, error) {
	return c.searchHubOpts.run(args)
}
