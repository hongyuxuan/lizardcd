/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tkn

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"text/template"

	"github.com/hongyuxuan/lizardcd/cli/common"
	"github.com/hongyuxuan/lizardcd/cli/types"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/spf13/cobra"
)

var name string
var params []string
var serviceaccount string
var timeout string
var workspace string
var cluster string
var namespace string
var templates string

type Param struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start a tekton pipelinerun with params",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()
		// get template
		var templateRes types.TemplateRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/db/yaml_template?page=1&size=1&search=name==%s", templates)).
			SetSuccessResult(&templateRes).
			Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get template: %s: %v", templates, err)
		}
		if templateRes.Data.Total == 0 {
			common.PrintFatal("cannot find template: %s", templates)
		}
		content := templateRes.Data.Results[0].Content
		// get piepeline info
		var res types.TektonPipelineRes
		if err := common.LizardServer.Get(fmt.Sprintf("/lizardcd/tekton/cluster/%s/namespace/%s/pipelines/%s", cluster, namespace, name)).
			SetSuccessResult(&res).
			Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to get pipeline: %s: %v", name, err)
		}
		if res.Code != http.StatusOK {
			common.PrintFatal("failed to get pipeline: %s: %v", name, res.Message)
		}
		// generate pipelineRun yaml
		var varParams []Param
		for _, param := range params {
			arr := strings.Split(param, "=")
			varParams = append(varParams, Param{
				Name:  arr[0],
				Value: arr[1],
			})
		}
		variables := map[string]interface{}{
			"Annotations":    res.Data.Annotations,
			"Labels":         res.Data.Labels,
			"Namespace":      namespace,
			"Pipeline":       name,
			"Workspace":      workspace,
			"ServiceAccount": serviceaccount,
			"Timeout":        timeout,
			"Params":         varParams,
		}
		tmpl, err := template.New("yamlTemplates").Parse(content)
		if err != nil {
			common.PrintFatal(err.Error())
		}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, variables); err != nil {
			common.PrintFatal(err.Error())
		}
		// create pipelineRun
		var createRes commontypes.Response
		if err := common.LizardServer.Post(fmt.Sprintf("/lizardcd/tekton/cluster/%s/namespace/%s/apply?kind=PipelineRun", cluster, namespace)).
			SetBody(map[string]interface{}{
				"content":   content,
				"variables": variables,
			}).SetSuccessResult(&createRes).
			Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to create pipelineRun: %v", err)
		}
		if createRes.Code != http.StatusOK {
			common.PrintFatal("failed to start pipelineRun: %v", createRes.Message)
		}
		common.PrintSuccess("successfully start pipeline: %s in namespace %s", name, namespace)
	},
}

func init() {
	startCmd.Flags().StringVar(&name, "name", "", "Pipeline name")
	startCmd.Flags().StringVar(&namespace, "namespace", "default", "Namespace of pipeline")
	startCmd.Flags().StringVar(&cluster, "cluster", "tektonk8s", "K8s cluster of pipeline")
	startCmd.Flags().StringArrayVarP(&params, "params", "p", []string{}, "Pipelinerun params, can be specified repeatedly")
	startCmd.Flags().StringVar(&serviceaccount, "sa", "default", "Serviceaccount for pipelinerun")
	startCmd.Flags().StringVar(&timeout, "timeout", "1h0m0s", "Timeout of pipelinerun")
	startCmd.Flags().StringVar(&workspace, "workspace", "shared-workspace", "Workspace for pipelinerun")
	startCmd.Flags().StringVar(&templates, "template", "tekton_template_common_pipelinerun", "PipelineRun template")
	startCmd.MarkFlagRequired("name")
}
