// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

// Package helm contains operations for working with helm charts.
package helm

import (
	"fmt"
	"strconv"

	"github.com/defenseunicorns/zarf/src/pkg/message"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"

	"helm.sh/helm/v3/pkg/chart/loader"
)

// loadChartFromTarball returns a helm chart from a tarball.
func (h *HelmCfg) loadChartFromTarball() (*chart.Chart, error) {
	// Get the path the temporary helm chart tarball
	sourceFile := StandardName(h.componentPaths.Charts, h.chart) + ".tgz"

	// Load the loadedChart tarball
	loadedChart, err := loader.Load(sourceFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load helm chart archive: %w", err)
	}

	if err = loadedChart.Validate(); err != nil {
		return nil, fmt.Errorf("unable to validate loaded helm chart: %w", err)
	}

	return loadedChart, nil
}

// parseChartValues reads the context of the chart values into an interface if it exists.
func (h *HelmCfg) parseChartValues() (map[string]any, error) {
	valueOpts := &values.Options{}

	for idx := range h.chart.ValuesFiles {
		path := StandardName(h.componentPaths.Values, h.chart) + "-" + strconv.Itoa(idx)
		valueOpts.ValueFiles = append(valueOpts.ValueFiles, path)
	}

	httpProvider := getter.Provider{
		Schemes: []string{"http", "https"},
		New:     getter.NewHTTPGetter,
	}

	providers := getter.Providers{httpProvider}
	return valueOpts.MergeValues(providers)
}

func (h *HelmCfg) createActionConfig(namespace string, spinner *message.Spinner) error {
	// Initialize helm SDK
	actionConfig := new(action.Configuration)
	// Set the setings for the helm SDK
	h.settings = cli.New()

	// Set the namespace for helm
	h.settings.SetNamespace(namespace)

	// Setup K8s connection
	err := actionConfig.Init(h.settings.RESTClientGetter(), namespace, "", spinner.Updatef)

	// Set the actionConfig is the received Helm pointer
	h.actionConfig = actionConfig

	return err
}
