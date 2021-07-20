package helmclient

import (
	"bytes"
	"github.com/outgnaY/helm-go-client/internal/version"
	"fmt"
	"text/template"
)

type versionClient interface {
	Version() (string, error)
}

type versionClientImpl struct {
	versionOpts *versionOptions
	env         *helmEnv
}

type VersionOption struct {
	f func(o *versionOptions)
}

type versionOptions struct {
	short    bool
	template string
}

func (o *versionOptions) apply(opts []VersionOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newVersionOptions(opts []VersionOption) *versionOptions {
	options := &versionOptions{
		short:    false,
		template: "",
	}
	options.apply(opts)
	return options
}

func VersionWithShort(short bool) VersionOption {
	return VersionOption{f: func(o *versionOptions) {
		o.short = short
	}}
}

func VersionWithTemplate(template string) VersionOption {
	return VersionOption{f: func(o *versionOptions) {
		o.template = template
	}}
}

func (c *versionClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *versionClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newVersionClient(opts []VersionOption, env *helmEnv) (*versionClientImpl, error) {
	o := newVersionOptions(opts)
	return &versionClientImpl{
		versionOpts: o,
		env:         env,
	}, nil
}

func (o *versionOptions) run() (string, error) {
	if o.template != "" {
		buf := new(bytes.Buffer)
		tt, err := template.New("_").Parse(o.template)
		if err != nil {
			return "", err
		}
		err = tt.Execute(buf, version.Get())
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	}
	return formatVersion(o.short), nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		if len(v.GitCommit) >= 7 {
			return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
		}
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", v)
}

func (c *versionClientImpl) Version() (string, error) {
	return c.versionOpts.run()
}
