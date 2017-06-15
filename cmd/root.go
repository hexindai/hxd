// Copyright © 2017 Runrioter Wung <runrioter@qq.com>
// Licensed under the MIT license.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "hxd",
	Short: "和信贷数据库相关命令行",
	Long: `
	
	该命令当前处于初期阶段，当前只支持mysql, 只有导出数据字典的功能。
	`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

}
