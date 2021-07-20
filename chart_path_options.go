package helmclient

import (
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

type ChartPathOption struct {
	f func(o *chartPathOptions)
}

type chartPathOptions action.ChartPathOptions

func (o *chartPathOptions) apply(opts []ChartPathOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func addChartPathOptions(opts []ChartPathOption, c *chartPathOptions) {
	c.apply(opts)
}

func WithCaFile(caFile string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.CaFile = caFile
	}}
}

func WithCertFile(certFile string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.CertFile = certFile
	}}
}

func WithKeyFile(keyFile string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.KeyFile = keyFile
	}}
}

func WithInsecureSkipTLSverify(insecureSkipTLSVerify bool) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.InsecureSkipTLSverify = insecureSkipTLSVerify
	}}
}

func WithKeyring(keyring string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.Keyring = keyring
	}}
}

func WithPassword(password string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.Password = password
	}}
}

func WithPassCredentialsAll(passCredentialsAll bool) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.PassCredentialsAll = passCredentialsAll
	}}
}

func WithRepoURL(repoURL string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.RepoURL = repoURL
	}}
}

func WithUsername(username string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.Username = username
	}}
}

func WithVerify(verify bool) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.Verify = verify
	}}
}

func WithVersion(version string) ChartPathOption {
	return ChartPathOption{f: func(o *chartPathOptions) {
		o.Version = version
	}}
}

func defaultKeyring() string {
	if v, ok := os.LookupEnv("GNUPGHOME"); ok {
		return filepath.Join(v, "pubring.gpg")
	}
	return filepath.Join(homedir.HomeDir(), ".gnupg", "pubring.gpg")
}

func mergeChartPathOptions(src *chartPathOptions, dst *action.ChartPathOptions) {
	src.CaFile = dst.CaFile
	src.CertFile = dst.CertFile
	src.KeyFile = dst.KeyFile
	src.InsecureSkipTLSverify = dst.InsecureSkipTLSverify
	src.Keyring = dst.Keyring
	src.Password = dst.Password
	src.PassCredentialsAll = dst.PassCredentialsAll
	src.RepoURL = dst.RepoURL
	src.Username = dst.Username
	src.Verify = dst.Verify
	src.Version = dst.Version
}
