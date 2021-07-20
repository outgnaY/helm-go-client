package helmclient

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"golang.org/x/term"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
	"time"
)

const (
	repoAddDefaultUsername              = ""
	repoAddDefaultPassword              = ""
	repoAddDefaultPassCredentialsAll    = false
	repoAddDefaultForceUpdate           = false
	repoAddDefaultAllowDeprecatedRepos  = false
	repoAddDefaultCertFile              = ""
	repoAddDefaultKeyFile               = ""
	repoAddDefaultCaFile                = ""
	repoAddDefaultInsecureSkipTLSverify = false
	repoAddDefaultDeprecateNoUpdate     = false
)

// Repositories that have been permanently deleted and no longer work
var deprecatedRepos = map[string]string{
	"//kubernetes-charts.storage.googleapis.com":           "https://charts.helm.sh/stable",
	"//kubernetes-charts-incubator.storage.googleapis.com": "https://charts.helm.sh/incubator",
}

type repoAddClient interface {
	RepoAdd(args []string) error
}

type repoAddClientImpl struct {
	repoAddOpts *repoAddOptions
	env         *helmEnv
}

type RepoAddOption struct {
	f func(o *repoAddOptions)
}

type repoAddOptions struct {
	name                 string
	url                  string
	username             string
	password             string
	passCredentialsAll   bool
	forceUpdate          bool
	allowDeprecatedRepos bool

	certFile              string
	keyFile               string
	caFile                string
	insecureSkipTLSverify bool

	repoFile  string
	repoCache string

	// Deprecated, but cannot be removed until Helm 4
	deprecatedNoUpdate bool
}

func (o *repoAddOptions) apply(opts []RepoAddOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newRepoAddOptions(opts []RepoAddOption) *repoAddOptions {
	options := &repoAddOptions{
		username:              repoAddDefaultUsername,
		password:              repoAddDefaultPassword,
		passCredentialsAll:    repoAddDefaultPassCredentialsAll,
		forceUpdate:           repoAddDefaultForceUpdate,
		allowDeprecatedRepos:  repoAddDefaultAllowDeprecatedRepos,
		certFile:              repoAddDefaultCertFile,
		keyFile:               repoAddDefaultKeyFile,
		caFile:                repoAddDefaultCaFile,
		insecureSkipTLSverify: repoAddDefaultInsecureSkipTLSverify,
		deprecatedNoUpdate:    repoAddDefaultDeprecateNoUpdate,
	}
	options.apply(opts)
	return options
}

func RepoAddWithUsername(username string) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.username = username
	}}
}

func RepoAddWithPassword(password string) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.password = password
	}}
}

func RepoAddWithPassCredentialsAll(passCredentialsAll bool) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.passCredentialsAll = passCredentialsAll
	}}
}

func RepoAddWithForceUpdate(forceUpdate bool) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.forceUpdate = forceUpdate
	}}
}

func RepoAddWithAllowDeprecatedRepos(allowDeprecatedRepos bool) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.allowDeprecatedRepos = allowDeprecatedRepos
	}}
}

func RepoAddWithCertFile(certFile string) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.certFile = certFile
	}}
}

func RepoAddWithKeyFile(keyFile string) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.keyFile = keyFile
	}}
}

func RepoAddWithDeprecatedNoUpdate(deprecatedNoUpdate bool) RepoAddOption {
	return RepoAddOption{f: func(o *repoAddOptions) {
		o.deprecatedNoUpdate = deprecatedNoUpdate
	}}
}

func (c *repoAddClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *repoAddClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env := rebuildEnv(globalOpts, namespace, c.env.clientGetter)
	c.env = env
	return nil
}

func newRepoAddClient(opts []RepoAddOption, env *helmEnv) (*repoAddClientImpl, error) {
	o := newRepoAddOptions(opts)
	return &repoAddClientImpl{
		repoAddOpts: o,
		env:         env,
	}, nil
}

func (o *repoAddOptions) run(out io.Writer, settings *cli.EnvSettings) error {
	// Block deprecated repos
	if !o.allowDeprecatedRepos {
		for oldURL, newURL := range deprecatedRepos {
			if strings.Contains(o.url, oldURL) {
				return fmt.Errorf("repo %q is no longer available; try %q instead", o.url, newURL)
			}
		}
	}

	// Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(o.repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	repoFileExt := filepath.Ext(o.repoFile)
	var lockPath string
	if len(repoFileExt) > 0 && len(repoFileExt) < len(o.repoFile) {
		lockPath = strings.Replace(o.repoFile, repoFileExt, ".lock", 1)
	} else {
		lockPath = o.repoFile + ".lock"
	}
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	b, err := ioutil.ReadFile(o.repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	if o.username != "" && o.password == "" {
		fd := int(os.Stdin.Fd())
		fmt.Fprint(out, "Password: ")
		password, err := term.ReadPassword(fd)
		fmt.Fprintln(out)
		if err != nil {
			return err
		}
		o.password = string(password)
	}

	c := repo.Entry{
		Name:                  o.name,
		URL:                   o.url,
		Username:              o.username,
		Password:              o.password,
		PassCredentialsAll:    o.passCredentialsAll,
		CertFile:              o.certFile,
		KeyFile:               o.keyFile,
		CAFile:                o.caFile,
		InsecureSkipTLSverify: o.insecureSkipTLSverify,
	}

	// If the repo exists do one of two things:
	// 1. If the configuration for the name is the same continue without error
	// 2. When the config is different require --force-update
	if !o.forceUpdate && f.Has(o.name) {
		existing := f.Get(o.name)
		if c != *existing {

			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return errors.Errorf("repository name (%s) already exists, please specify a different name", o.name)
		}

		// The add is idempotent so do nothing
		fmt.Fprintf(out, "%q already exists with the same configuration, skipping\n", o.name)
		return nil
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		return err
	}

	if o.repoCache != "" {
		r.CachePath = o.repoCache
	}
	if _, err := r.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", o.url)
	}

	f.Update(&c)

	if err := f.WriteFile(o.repoFile, 0644); err != nil {
		return err
	}
	fmt.Fprintf(out, "%q has been added to your repositories\n", o.name)
	return nil
}

func (c *repoAddClientImpl) RepoAdd(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("repo add requires 2 arguments exactly")
	}
	c.repoAddOpts.name = args[0]
	c.repoAddOpts.url = args[1]
	c.repoAddOpts.repoFile = c.env.settings.RepositoryConfig
	c.repoAddOpts.repoCache = c.env.settings.RepositoryCache
	return c.repoAddOpts.run(os.Stdout, c.env.settings)
}
