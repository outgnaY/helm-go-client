package helmclient

import (
	"helm.sh/helm/v3/pkg/cli/values"
)

type ValueOption struct {
	f func(o *valueOptions)
}

type valueOptions values.Options

func (o *valueOptions) apply(opts []ValueOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func addValueOptions(opts []ValueOption, v *valueOptions) {
	v.apply(opts)
}

func WithValueFiles(valueFiles []string) ValueOption {
	return ValueOption{f: func(o *valueOptions) {
		o.ValueFiles = valueFiles
	}}
}

func WithStringValues(stringValues []string) ValueOption {
	return ValueOption{f: func(o *valueOptions) {
		o.StringValues = stringValues
	}}
}

func WithValues(values []string) ValueOption {
	return ValueOption{f: func(o *valueOptions) {
		o.Values = values
	}}
}

func WithFileValues(fileValues []string) ValueOption {
	return ValueOption{f: func(o *valueOptions) {
		o.FileValues = fileValues
	}}
}
