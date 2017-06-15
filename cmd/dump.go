// Copyright © 2017 Runrioter Wung <runrioter@qq.com>
// Licensed under the MIT license.

package cmd

import (
	"fmt"

	"bufio"
	"os"

	"syscall"

	"github.com/hexindai/hxd/dump"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var schema string
var table string
var host string
var port string

var username string
var password string

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "导出数据库的数据字典并生成excel文件",
	Long: `该命令为和信贷数据字典导出命令。例子如下：
hxd dump --host=[主机ip] --port=[端口号] --username=[用户名] --password=[密码] --schema=[库名] --table=[表名]`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(schema) > 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for len(username) == 0 {
				fmt.Print("请键入数据库用户名:")
				if scanner.Scan() {
					username = scanner.Text()
				}
			}
			for len(password) == 0 {
				fmt.Print("请键入数据库密码:")
				bytes, err := terminal.ReadPassword(int(syscall.Stdin))
				if err != nil {
					continue
				}
				password = string(bytes)
			}
			dump.GenerateExcel(host, port, username, password, schema, table)
		} else {
			fmt.Println("hxd dump --host=[主机ip] --port=[端口号] --username=[用户名] --password=[密码] --schema=[库名] --table=[表名]")
		}
	},
}

func init() {
	RootCmd.AddCommand(dumpCmd)

	// Here you will define your flags and configuration settings.
	dumpCmd.Flags().StringVarP(&schema, "schema", "s", "", "数据库名")
	dumpCmd.Flags().StringVarP(&table, "table", "t", "", "数据表名")
	dumpCmd.Flags().StringVarP(&username, "username", "u", "", "数据库用户名")
	dumpCmd.Flags().StringVarP(&password, "password", "p", "", "数据库密码")
	dumpCmd.Flags().StringVar(&host, "host", "127.0.0.1", "主机名")
	dumpCmd.Flags().StringVar(&port, "port", "3306", "端口号")
}
