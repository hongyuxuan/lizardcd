/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package apply

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hongyuxuan/lizardcd/cli/common"
	commontypes "github.com/hongyuxuan/lizardcd/common/types"
	"github.com/spf13/cobra"
)

var namespace string
var cluster string
var manifest string
var variables string

// applyCmd represents the apply command
var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "apply a go-template file to be applyed",
	Run: func(cmd *cobra.Command, args []string) {
		common.InitConfig()

		b, err := ioutil.ReadFile(manifest)
		if err != nil {
			common.PrintFatal("failed to read file %s: %v", manifest, err)
		}
		vars := make(map[string]interface{})
		if variables != "" {
			for _, v := range strings.Split(variables, ",") {
				arr := strings.Split(v, "=")
				vars[arr[0]] = arr[1]
			}
		}
		var res *commontypes.Response
		if err := common.LizardServer.Patch(fmt.Sprintf("/lizardcd/kubernetes/cluster/%s/namespace/%s/apply/yaml", cluster, namespace)).SetBody(map[string]interface{}{
			"content":   string(b),
			"variables": vars,
		}).SetResult(&res).Do(context.Background()).Err; err != nil {
			common.PrintFatal("failed to apply file \"%s\" to cluster=%s, namespace=%s: %v", manifest, cluster, namespace, err)
		} else {
			common.PrintSuccess(res.Message)
		}
	},
}

func init() {
	ApplyCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "namespace of a kubernetes cluster")
	ApplyCmd.Flags().StringVar(&cluster, "cluster", "", "kubernetes cluster name")
	ApplyCmd.Flags().StringVarP(&manifest, "file", "f", "", "go-template file to be applyed")
	ApplyCmd.Flags().StringVarP(&variables, "vars", "v", "", "go-template variables for template, format: Appname=test,Namespace=default,Port=3000")
	ApplyCmd.MarkFlagRequired("cluster")
	ApplyCmd.MarkFlagRequired("file")
}
