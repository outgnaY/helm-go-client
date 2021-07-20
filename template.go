package helmclient

import (
	"bytes"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	templateDefaultOutputDir      = ""
	templateDefaultValidate       = false
	templateDefaultIncludeCRDs    = false
	templateDefaultSkipTests      = false
	templateDefaultIsUpgrade      = false
	templateDefaultKubeVersion    = ""
	templateDefaultUseReleaseName = false
)

var (
	templateDefaultShowFiles = []string{}
	templateDefaultExtraAPIs = []string{}
)

type templateClient interface {
	Template(args []string) error
}

type templateClientImpl struct {
	cli         *action.Install
	env         *helmEnv
	valueOpts   *valueOptions
	showFiles   []string
	skipTests   bool
	kubeVersion string
	extraAPIs   []string
}

type TemplateOption struct {
	f func(o *templateOptions)
}

type templateOptions struct {
	showFiles      []string
	outputDir      string
	validate       bool
	includeCRDs    bool
	skipTests      bool
	isUpgrade      bool
	kubeVersion    string
	extraAPIs      []string
	useReleaseName bool
}

func (o *templateOptions) apply(opts []TemplateOption) {
	for _, op := range opts {
		op.f(o)
	}
}

func newTemplateOptions(opts []TemplateOption) *templateOptions {
	options := &templateOptions{
		showFiles:      templateDefaultShowFiles,
		outputDir:      templateDefaultOutputDir,
		validate:       templateDefaultValidate,
		includeCRDs:    templateDefaultIncludeCRDs,
		skipTests:      templateDefaultSkipTests,
		isUpgrade:      templateDefaultIsUpgrade,
		kubeVersion:    templateDefaultKubeVersion,
		extraAPIs:      templateDefaultExtraAPIs,
		useReleaseName: templateDefaultUseReleaseName,
	}
	options.apply(opts)
	return options
}

func TemplateWithShowFiles(showFiles []string) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.showFiles = showFiles
	}}
}

func TemplateWithOutputDir(outputDir string) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.outputDir = outputDir
	}}
}

func TemplateWithValidate(validate bool) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.validate = validate
	}}
}

func TemplateWithIncludeCRDs(includeCRDs bool) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.includeCRDs = includeCRDs
	}}
}

func TemplateWithSkipTests(skipTests bool) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.skipTests = skipTests
	}}
}

func TemplateWithIsUpgrade(isUpgrade bool) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.isUpgrade = isUpgrade
	}}
}

func TemplateWithKubeVersion(kubeVersion string) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.kubeVersion = kubeVersion
	}}
}

func TemplateWithExtraAPIs(extraAPIs []string) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.extraAPIs = extraAPIs
	}}
}

func TemplateWithUseReleaseName(useReleaseName bool) TemplateOption {
	return TemplateOption{f: func(o *templateOptions) {
		o.useReleaseName = useReleaseName
	}}
}

func (c *templateClientImpl) OverrideGlobalOpts(globalOpts []GlobalOption) error {
	// use old namespace
	return c.OverrideGlobalOptsWithNamespace(globalOpts, c.env.Namespace())
}

func (c *templateClientImpl) OverrideGlobalOptsWithNamespace(globalOpts []GlobalOption, namespace string) error {
	env, cfg, err := rebuildEnvAndCfg(globalOpts, namespace, c.env.clientGetter)
	if err != nil {
		return err
	}
	client := action.NewInstall(cfg)
	// copy args
	copyInstallClientOptions(c.cli, client)
	c.cli = client
	c.env = env
	return nil
}

func newTemplateClient(opts []TemplateOption, valueOpts []ValueOption, env *helmEnv) (*templateClientImpl, error) {
	o := newTemplateOptions(opts)
	cfg := new(action.Configuration)
	err := cfg.Init(env.clientGetter, env.Namespace(), "", debug)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(cfg)
	mergeTemplateOptions(o, client)
	v := &valueOptions{
		ValueFiles:   []string{},
		StringValues: []string{},
		Values:       []string{},
		FileValues:   []string{},
	}
	addValueOptions(valueOpts, v)

	return &templateClientImpl{
		cli:         client,
		env:         env,
		valueOpts:   v,
		showFiles:   o.showFiles,
		skipTests:   o.skipTests,
		kubeVersion: o.kubeVersion,
		extraAPIs:   o.extraAPIs,
	}, nil
}

func (c *templateClientImpl) Template(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("template requires at least 1 argument")
	}
	if c.kubeVersion != "" {
		parsedKubeVersion, err := chartutil.ParseKubeVersion(c.kubeVersion)
		if err != nil {
			return fmt.Errorf("invalid kube version '%s': %s", c.kubeVersion, err)
		}
		c.cli.KubeVersion = parsedKubeVersion
	}
	c.cli.DryRun = true
	c.cli.ReleaseName = "RELEASE-NAME"
	c.cli.Replace = true // Skip the name check
	c.cli.APIVersions = chartutil.VersionSet(c.extraAPIs)
	rel, err := runInstall(args, c.cli, (*values.Options)(c.valueOpts), os.Stdout, c.env)

	if err != nil && !c.env.settings.Debug {
		if rel != nil {
			return fmt.Errorf("%w\n\nUse --debug flag to render out invalid YAML", err)
		}
		return err
	}

	// We ignore a potential error here because, when the --debug flag was specified,
	// we always want to print the YAML, even if it is not valid. The error is still returned afterwards.
	if rel != nil {
		var manifests bytes.Buffer
		fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))
		if !c.cli.DisableHooks {
			fileWritten := make(map[string]bool)
			for _, m := range rel.Hooks {
				if c.skipTests && isTestHook(m) {
					continue
				}
				if c.cli.OutputDir == "" {
					fmt.Fprintf(&manifests, "---\n# Source: %s\n%s\n", m.Path, m.Manifest)
				} else {
					newDir := c.cli.OutputDir
					if c.cli.UseReleaseName {
						newDir = filepath.Join(c.cli.OutputDir, c.cli.ReleaseName)
					}
					err = writeToFile(newDir, m.Path, m.Manifest, fileWritten[m.Path])
					if err != nil {
						return err
					}
					fileWritten[m.Path] = true
				}

			}
		}

		// if we have a list of files to render, then check that each of the
		// provided files exists in the chart.
		if len(c.showFiles) > 0 {
			// This is necessary to ensure consistent manifest ordering when using --show-only
			// with globs or directory names.
			splitManifests := releaseutil.SplitManifests(manifests.String())
			manifestsKeys := make([]string, 0, len(splitManifests))
			for k := range splitManifests {
				manifestsKeys = append(manifestsKeys, k)
			}
			sort.Sort(releaseutil.BySplitManifestsOrder(manifestsKeys))

			manifestNameRegex := regexp.MustCompile("# Source: [^/]+/(.+)")
			var manifestsToRender []string
			for _, f := range c.showFiles {
				missing := true
				// Use linux-style filepath separators to unify user's input path
				f = filepath.ToSlash(f)
				for _, manifestKey := range manifestsKeys {
					manifest := splitManifests[manifestKey]
					submatch := manifestNameRegex.FindStringSubmatch(manifest)
					if len(submatch) == 0 {
						continue
					}
					manifestName := submatch[1]
					// manifest.Name is rendered using linux-style filepath separators on Windows as
					// well as macOS/linux.
					manifestPathSplit := strings.Split(manifestName, "/")
					// manifest.Path is connected using linux-style filepath separators on Windows as
					// well as macOS/linux
					manifestPath := strings.Join(manifestPathSplit, "/")

					// if the filepath provided matches a manifest path in the
					// chart, render that manifest
					if matched, _ := filepath.Match(f, manifestPath); !matched {
						continue
					}
					manifestsToRender = append(manifestsToRender, manifest)
					missing = false
				}
				if missing {
					return fmt.Errorf("could not find template %s in chart", f)
				}
			}
			for _, m := range manifestsToRender {
				fmt.Fprintf(os.Stdout, "---\n%s\n", m)
			}
		} else {
			fmt.Fprintf(os.Stdout, "%s", manifests.String())
		}
	}

	return err
}

func isTestHook(h *release.Hook) bool {
	for _, e := range h.Events {
		if e == release.HookTest {
			return true
		}
	}
	return false
}

// The following functions (writeToFile, createOrOpenFile, and ensureDirectoryForFile)
// are copied from the actions package. This is part of a change to correct a
// bug introduced by #8156. As part of the todo to refactor renderResources
// this duplicate code should be removed. It is added here so that the API
// surface area is as minimally impacted as possible in fixing the issue.
func writeToFile(outputDir string, name string, data string, append bool) error {
	outfileName := strings.Join([]string{outputDir, name}, string(filepath.Separator))

	err := ensureDirectoryForFile(outfileName)
	if err != nil {
		return err
	}

	f, err := createOrOpenFile(outfileName, append)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("---\n# Source: %s\n%s\n", name, data))

	if err != nil {
		return err
	}

	fmt.Printf("wrote %s\n", outfileName)
	return nil
}

func createOrOpenFile(filename string, append bool) (*os.File, error) {
	if append {
		return os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
	}
	return os.Create(filename)
}

func ensureDirectoryForFile(file string) error {
	baseDir := path.Dir(file)
	_, err := os.Stat(baseDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return os.MkdirAll(baseDir, 0755)
}

func mergeTemplateOptions(o *templateOptions, cli *action.Install) {
	cli.OutputDir = o.outputDir
	cli.IsUpgrade = o.isUpgrade
	cli.UseReleaseName = o.useReleaseName
	cli.ClientOnly = !o.validate
	cli.IncludeCRDs = o.includeCRDs
}
