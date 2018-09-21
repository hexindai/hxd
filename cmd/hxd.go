package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// hxd represents the base command when called without any subcommands
var hxd = &cobra.Command{
	Use:   "hxd",
	Short: "和信贷命令行",
	Long:  "\n1.导出数据字典的功能，当前只支持mysql\n2.动态实时查看基于http协议的日志文件",
}

// Execute run hxd
func Execute() {
	if err := hxd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
