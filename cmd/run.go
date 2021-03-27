/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/sealyun/cloud-kernel/pkg/logger"
	"github.com/sealyun/cloud-kernel/pkg/vars"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "执行打包离线包并发布到sealyun上",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("run called")
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if vars.AkId == "" {
			logger.Error("云厂商的akId为空,无法创建虚拟机")
			cmd.Help()
			os.Exit(-1)
		}
		if vars.AkSK == "" {
			logger.Error("云厂商的akSK为空,无法创建虚拟机")
			cmd.Help()
			os.Exit(0)
		}
		if vars.MarketCtlToken == "" {
			logger.Error("MarketCtl的Token为空无法上传离线包")
			cmd.Help()
			os.Exit(0)
		}
		if vars.DingDing == "" {
			logger.Warn("钉钉的Token为空,无法自动通知")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	runCmd.Flags().StringVar(&vars.AkId, "ak", "", "云厂商的 akId")
	runCmd.Flags().StringVar(&vars.AkSK, "sk", "", "云厂商的 akSK")
	runCmd.Flags().StringVar(&vars.DingDing, "dingding", "", "钉钉的Token")
	runCmd.Flags().StringVar(&vars.MarketCtlToken, "marketctl", "", "marketctl的token")
	runCmd.Flags().BoolVar(&vars.IsArm64, "amd64", false, "是否为amd64")
	runCmd.Flags().Float64Var(&vars.DefaultPrice, "price", 50, "离线包的价格")
	runCmd.Flags().Float64Var(&vars.DefaultZeroPrice, "price", 0.01, "离线包.0版本的价格")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
