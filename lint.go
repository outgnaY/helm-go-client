package helmclient

import (
	"errors"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"os"
	"path/filepath"
	"strings"
)

const (
	lintDefaultStrict        = false
	lintDefaultWithSubcharts = false
)

type lintClient interface {
	Lint(args []string) error
}

type lintClientImpl struct {
	cli       *action.Lint
	env       *helmEnv
	valueOpts *valueOptions
}

type LintOption struct {
	f func(o *lintOptions)
}

type lintOptions struct {
	strict        bool
	withSubcharts bool
}

func (o *lintOptions) apply(opts []LintOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newLintOptions(opts []LintOption) *lintOptions {
	options := &lintOptions{
		strict:        lintDefaultStrict,
		withSubcharts: lintDefaultWithSubcharts,
	}
	options.apply(opts)
	return options
}

func LintWithStrict(strict bool) LintOption {
	return LintOption{f: func(o *lintOptions) {
		o.strict = strict
	}}
}

func LintWithWithSubcharts(withSubcharts bool) LintOption {
	return LintOption{f: func(o *lintOptions) {
		o.withSubcharts = withSubcharts
	}}
}

func (c *lintClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *lintClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	client := action.NewLint()
	// copy args
	copyLintClientOptions(c.cli, client)
	client.Namespace = env.Namespace()
	c.cli = client
	c.env = env
	return nil
}

func copyLintClientOptions(oldCli *action.Lint, newCli *action.Lint) {
	newCli.Strict = oldCli.Strict
	newCli.WithSubcharts = oldCli.WithSubcharts
}

func newLintClient(opts []LintOption, valueOpts []ValueOption, env *helmEnv) (*lintClientImpl, error) {
	o := newLintOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewLint()
	mergeLintOptions(o, client)
	v := &valueOptions{
		ValueFiles:   []string{},
		StringValues: []string{},
		Values:       []string{},
		FileValues:   []string{},
	}
	addValueOptions(valueOpts, v)
	return &lintClientImpl{
		cli:       client,
		env:       env,
		valueOpts: v,
	}, nil
}

func (c *lintClientImpl) Lint(args []string) error {
	paths := []string{"."}
	if len(args) > 0 {
		paths = args
	}
	if c.cli.WithSubcharts {
		for _, p := range paths {
			filepath.Walk(filepath.Join(p, "charts"), func(path string, info os.FileInfo, err error) error {
				if info != nil {
					if info.Name() == "Chart.yaml" {
						paths = append(paths, filepath.Dir(path))
					} else if strings.HasSuffix(path, ".tgz") || strings.HasSuffix(path, ".tar.gz") {
						paths = append(paths, path)
					}
				}
				return nil
			})
		}
	}
	c.cli.Namespace = c.env.Namespace()
	vals, err := ((*values.Options)(c.valueOpts)).MergeValues(getter.All(c.env.settings))
	if err != nil {
		return err
	}
	var message strings.Builder
	failed := 0
	for _, path := range paths {
		fmt.Fprintf(&message, "==> Linting %s\n", path)

		result := c.cli.Run([]string{path}, vals)

		// All the Errors that are generated by a chart
		// that failed a lint will be included in the
		// results.Messages so we only need to print
		// the Errors if there are no Messages.
		if len(result.Messages) == 0 {
			for _, err := range result.Errors {
				fmt.Fprintf(&message, "Error %s\n", err)
			}
		}

		for _, msg := range result.Messages {
			fmt.Fprintf(&message, "%s\n", msg)
		}

		if len(result.Errors) != 0 {
			failed++
		}

		// Adding extra new line here to break up the
		// results, stops this from being a big wall of
		// text and makes it easier to follow.
		fmt.Fprint(&message, "\n")
	}
	fmt.Fprint(os.Stdout, message.String())

	summary := fmt.Sprintf("%d chart(s) linted, %d chart(s) failed", len(paths), failed)
	if failed > 0 {
		return errors.New(summary)
	}
	fmt.Fprintln(os.Stdout, summary)
	return nil
}

func mergeLintOptions(o *lintOptions, cli *action.Lint) {
	cli.Strict = o.strict
	cli.WithSubcharts = o.withSubcharts
}
