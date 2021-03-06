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
	"encoding/json"
	"github.com/labring/cloud-kernel/pkg/ecs"
	"github.com/labring/cloud-kernel/pkg/github"
	"github.com/labring/cloud-kernel/pkg/logger"
	"github.com/labring/cloud-kernel/pkg/marketctl"
	_package "github.com/labring/cloud-kernel/pkg/package"
	"github.com/labring/cloud-kernel/pkg/retry"
	"github.com/labring/cloud-kernel/pkg/vars"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var gFetch []string

type errorObject struct {
	Version string `json:"version"`
	Error   string `json:"error"`
}

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "执行打包离线包并发布到sealyun上",
	Run: func(cmd *cobra.Command, args []string) {
		if len(gFetch) == 0 {
			logger.Warn("当月无需要更新版本")
			os.Exit(0)
		} else {
			e := make([]errorObject, 0)
			for _, v := range gFetch {
				logger.Debug("当前更新版本: " + v)
				if err := _package.Package(strings.ReplaceAll(v, "v", "")); err != nil {
					logger.Error(err)
					logger.Warn("更新版本发生错误,跳过当前版本: " + v)
					e = append(e, errorObject{
						Version: v,
						Error:   err.Error(),
					})
				}
			}
			if len(e) != 0 {
				s, _ := json.Marshal(e)
				logger.Fatal("执行打包流程报错 %s", string(s))
				os.Exit(0)
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if vars.AkID == "" {
			logger.Fatal("云厂商的akId为空,无法创建虚拟机")
			cmd.Help()
			os.Exit(-1)
		}
		if vars.AkSK == "" {
			logger.Fatal("云厂商的akSK为空,无法创建虚拟机")
			cmd.Help()
			os.Exit(0)
		}
		cloud := ecs.NewCloud()
		if err := cloud.Healthy(); err != nil {
			logger.Fatal("云厂商的AKSK验证失败: " + err.Error())
			cmd.Help()
			os.Exit(0)
		}
		if vars.DingDing == "" {
			logger.Warn("钉钉的Token为空,无法自动通知")
		}
		if vars.Uploading {
			if vars.MarketCtlToken == "" {
				logger.Fatal("MarketCtl的Token为空无法上传离线包")
				cmd.Help()
				os.Exit(0)
			}
			if err := retry.Do(func() error {
				err := marketctl.Healthy()
				if err != nil {
					return err
				}
				gf, err := github.Fetch()
				if err != nil {
					return err
				}
				gFetch = gf
				return nil
			}, 25, 1*time.Second, false); err != nil {
				logger.Fatal("Sealyun的状态监测失败: " + err.Error())
				cmd.Help()
				os.Exit(0)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.
	runCmd.Flags().StringVar(&vars.AkID, "ak", "", "云厂商的 akId")
	runCmd.Flags().StringVar(&vars.AkSK, "sk", "", "云厂商的 akSK")
	runCmd.Flags().StringVar(&vars.DingDing, "dingding", "", "钉钉的Token")
	runCmd.Flags().StringVar(&vars.MarketCtlToken, "marketctl", "", "marketctl的token")
	runCmd.Flags().BoolVar(&vars.IsArm64, "arm64", false, "是否为arm64")
	runCmd.Flags().Float64Var(&vars.DefaultPrice, "price", 50, "离线包的价格")
	runCmd.Flags().Float64Var(&vars.DefaultZeroPrice, "zoro-price", 0.01, "离线包.0版本的价格")
	runCmd.Flags().BoolVar(&vars.Uploading, "uploading", true, "是否上传")
	runCmd.Flags().BoolVar(&vars.Testing, "testing", true, "是否测试")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
