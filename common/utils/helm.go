package utils

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hongyuxuan/lizardcd/common/errorx"
	"github.com/hongyuxuan/lizardcd/common/types"
	"github.com/zeromicro/go-zero/core/logx"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

type HelmUtil struct {
	logx.Logger
	ctx      context.Context
	settings *cli.EnvSettings
}

func NewHelmUtil(ctx context.Context) *HelmUtil {
	return &HelmUtil{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		settings: cli.New(),
	}
}

func (h *HelmUtil) Update(repoList []*repo.Entry) error {
	repoFile := h.settings.RepositoryConfig
	repositories, err := repo.LoadFile(h.settings.RepositoryConfig)
	if err != nil {
		repositories = repo.NewFile()
	}
	repositories.Repositories = repoList
	if err := repositories.WriteFile(repoFile, 0644); err != nil {
		return errorx.NewDefaultError("Cannot save repository config file to %s: %s", repoFile, err)
	}
	for _, repoEntry := range repositories.Repositories {
		chartRepository, err := repo.NewChartRepository(repoEntry, getter.All(h.settings))
		chartRepository.CachePath = h.settings.RepositoryCache
		if err != nil {
			return errorx.NewDefaultError("Cannot new repository charts: %s", err)
		}
		// update repo index
		cache, err := chartRepository.DownloadIndexFile()
		if err != nil {
			return errorx.NewDefaultError("Cannot download repository index: %s", err)
		}
		h.Logger.Infof("Successfully got an update from the \"%s\" repository to cache: %s", repoEntry.Name, cache)
	}
	return nil
}

func (h *HelmUtil) SearchChart(repoUrl, repoName, chartName string) ([]*types.ChartListResponse, error) {
	path := fmt.Sprintf("%s/%s-index.yaml", h.settings.RepositoryCache, repoName)
	indexFile, err := repo.LoadIndexFile(path)
	if err != nil {
		return nil, errorx.NewDefaultError("Repository \"%s\" index %s not exists", repoName, path)
	}
	var chartList []*types.ChartListResponse
	for _, entry := range indexFile.Entries {
		if strings.Contains(entry[0].Name, chartName) {
			chart := &types.ChartListResponse{
				RepoUrl:      repoUrl,
				ChartName:    entry[0].Name,
				ChartVersion: entry[0].Version,
				AppVersion:   entry[0].AppVersion,
				Icon:         entry[0].Icon,
				Description:  entry[0].Description,
			}
			chartList = append(chartList, chart)
		}
	}

	return chartList, nil
}

func (h *HelmUtil) SearchChartVersions(repoName, chartName string) ([]*types.ChartListResponse, error) {
	path := fmt.Sprintf("%s/%s-index.yaml", h.settings.RepositoryCache, repoName)
	indexFile, err := repo.LoadIndexFile(path)
	if err != nil {
		return nil, errorx.NewDefaultError("Repository \"%s\" index not exists", repoName)
	}

	var chartList []*types.ChartListResponse
	for _, entry := range indexFile.Entries[chartName] {
		chart := &types.ChartListResponse{
			ChartName:    entry.Name,
			ChartVersion: entry.Version,
			AppVersion:   entry.AppVersion,
			Description:  entry.Description,
		}
		chartList = append(chartList, chart)
	}

	return chartList, nil
}

func (h *HelmUtil) InstallChart(namespace, kubeconfig, repoUrl, chartName, chartVersion, releaseName string, valuesYaml []byte, wait bool, timeout time.Duration) (err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return errorx.NewDefaultError("Init install action failed: %v", err)
	}
	install := action.NewInstall(actionConfig)
	install.RepoURL = repoUrl
	install.Version = chartVersion
	install.Timeout = timeout
	install.CreateNamespace = false
	install.Wait = wait
	install.Namespace = namespace
	install.ReleaseName = releaseName

	chartRequested, err := install.ChartPathOptions.LocateChart(chartName, h.settings)
	if err != nil {
		return errorx.NewDefaultError("Download chart \"%s\" for install failed: %v", chartName, err)
	}

	chart, err := loader.Load(chartRequested)
	if err != nil {
		return errorx.NewDefaultError("Load chart \"%s\" for install failed: %v", chartName, err)
	}

	var mergedValues map[string]interface{}
	if valuesYaml != nil {
		os.WriteFile("./values.yaml", valuesYaml, 0644)
		providers := getter.All(h.settings)
		options := values.Options{
			ValueFiles: []string{"./values.yaml"},
		}
		if mergedValues, err = options.MergeValues(providers); err != nil {
			return
		}
	}

	if _, err = install.Run(chart, mergedValues); err != nil {
		h.Logger.Error(err)
		return
	}
	h.Logger.Infof("Successfully installed chart \"%s\"", chartName)
	os.Remove("./values.yaml")
	return nil
}

func (h *HelmUtil) UninstallChart(namespace, kubeconfig, releaseName string) (err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return errorx.NewDefaultError("Init install action failed: %v", err)
	}

	uninstall := action.NewUninstall(actionConfig)
	uninstall.Timeout = 60 * time.Second // 设置超时时间60秒
	uninstall.Wait = true
	uninstall.KeepHistory = false

	if _, err = uninstall.Run(releaseName); err != nil {
		h.Logger.Error(err)
		return
	}

	h.Logger.Infof("Successfully uninstalled chart \"%s\"", releaseName)
	return nil
}

func (h *HelmUtil) UpgradeChart(namespace, kubeconfig, repoUrl, releaseName, chartName, chartVersion string, revision int, valuesYaml []byte, wait bool, timeout time.Duration) (err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return errorx.NewDefaultError("Init install action failed: %v", err)
	}

	// helm upgrade
	upgrade := action.NewUpgrade(actionConfig)
	upgrade.RepoURL = repoUrl
	upgrade.Version = chartVersion
	upgrade.Namespace = namespace
	upgrade.Wait = wait
	upgrade.Timeout = timeout

	chartRequested, err := upgrade.ChartPathOptions.LocateChart(chartName, h.settings)
	if err != nil {
		return errorx.NewDefaultError("Download chart \"%s\" for install failed: %v", chartName, err)
	}

	chart, err := loader.Load(chartRequested)
	if err != nil {
		return errorx.NewDefaultError("Load chart \"%s\" for install failed: %v", chartName, err)
	}

	var mergedValues map[string]interface{}
	if valuesYaml != nil {
		os.WriteFile("./values.yaml", valuesYaml, 0644)
		providers := getter.All(h.settings)
		options := values.Options{
			ValueFiles: []string{"./values.yaml"},
		}
		if mergedValues, err = options.MergeValues(providers); err != nil {
			return
		}
	}

	if _, err = upgrade.Run(releaseName, chart, mergedValues); err != nil {
		h.Logger.Error(err)
		return
	}
	h.Logger.Infof("Successfully upgrade chart \"%s\"", chartName)
	os.Remove("./values.yaml")
	return nil
}

func (h *HelmUtil) ListRelease(namespace, kubeconfig, releaseName string) (elements []types.ReleaseElement, err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err = actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return nil, errorx.NewDefaultError("Init list action failed: %v", err)
	}
	list := action.NewList(actionConfig)
	list.All = true // fetch all status release
	list.SetStateMask()
	var releases []*release.Release
	if releases, err = list.Run(); err != nil {
		return
	}
	elements = make([]types.ReleaseElement, 0, len(releases))
	for _, r := range releases {
		if strings.Contains(r.Name, releaseName) {
			element := types.ReleaseElement{
				Name:         r.Name,
				Namespace:    r.Namespace,
				Revision:     strconv.Itoa(r.Version),
				Status:       r.Info.Status.String(),
				Chart:        formatChartname(r.Chart),
				ChartName:    r.Chart.Name(),
				ChartVersion: r.Chart.Metadata.Version,
				AppVersion:   formatAppVersion(r.Chart),
			}
			t := "-"
			if tspb := r.Info.LastDeployed; !tspb.IsZero() {
				t = tspb.Format("2006-01-02 15:04:05Z0700")
			}
			element.Updated = t
			elements = append(elements, element)
		}
	}
	return
}

func (h *HelmUtil) ShowDefaultValues(repoUrl, chartName, chartVersion string) (output string, err error) {
	actionConfig := new(action.Configuration)

	show := action.NewShowWithConfig(action.ShowValues, actionConfig)
	show.RepoURL = repoUrl
	show.Version = chartVersion
	var cp string
	if cp, err = show.ChartPathOptions.LocateChart(chartName, h.settings); err != nil {
		return "", err
	}
	return show.Run(cp)
}

func (h *HelmUtil) ShowReadme(repoUrl, chartName, chartVersion string) (output string, err error) {
	actionConfig := new(action.Configuration)

	show := action.NewShowWithConfig(action.ShowReadme, actionConfig)
	show.RepoURL = repoUrl
	show.Version = chartVersion
	var cp string
	if cp, err = show.ChartPathOptions.LocateChart(chartName, h.settings); err != nil {
		return "", err
	}
	return show.Run(cp)
}

func (h *HelmUtil) GetValues(namespace, kubeconfig, releaseName string, revision int) (output []byte, err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return nil, errorx.NewDefaultError("Init get values action failed: %v", err)
	}
	getvalue := action.NewGetValues(actionConfig)
	getvalue.Version = revision
	getvalue.AllValues = true
	var res map[string]interface{}
	if res, err = getvalue.Run(releaseName); err != nil {
		h.Logger.Error(err)
		return
	}
	b, _ := yaml.Marshal(res)
	return b, nil
}

func (h *HelmUtil) ReleaseHistory(namespace, kubeconfig, releaseName string) (res []types.ReleaseHistoryInfo, err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return nil, errorx.NewDefaultError("Init history action failed: %v", err)
	}
	history := action.NewHistory(actionConfig)
	var releases []*release.Release
	if releases, err = history.Run(releaseName); err != nil {
		h.Logger.Error(err)
		return
	}
	for _, r := range releases {
		rInfo := types.ReleaseHistoryInfo{
			Revision:    r.Version,
			Status:      r.Info.Status.String(),
			Chart:       formatChartname(r.Chart),
			AppVersion:  formatAppVersion(r.Chart),
			Description: r.Info.Description,
		}
		if !r.Info.LastDeployed.IsZero() {
			rInfo.Updated = r.Info.LastDeployed.Time
		}
		res = append(res, rInfo)
	}
	return
}

func (h *HelmUtil) Rollback(namespace, kubeconfig, releaseName string, revision int, wait bool, timeout time.Duration) (err error) {
	actionConfig := new(action.Configuration)
	if kubeconfig != "" {
		h.settings.KubeConfig = kubeconfig
	}
	if err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), h.Logger.Debugf); err != nil {
		return errorx.NewDefaultError("Init rollback action failed: %v", err)
	}
	rollback := action.NewRollback(actionConfig)
	rollback.Version = revision
	rollback.Wait = wait
	rollback.Timeout = timeout
	if err = rollback.Run(releaseName); err != nil {
		h.Logger.Error(err)
		return
	}
	return nil
}

func (h *HelmUtil) Pull(repoUrl, chartName, chartVersion, destDir string) (output string, err error) {
	actionConfig := new(action.Configuration)

	pull := action.NewPullWithOpts(action.WithConfig(actionConfig))
	pull.RepoURL = repoUrl
	pull.Version = chartVersion
	pull.DestDir = destDir
	pull.Settings = h.settings

	return pull.Run(chartName)
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
