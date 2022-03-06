package main

import (
	"fmt"
	"gortcenter/app/cmd"
	"gortcenter/pkg/config"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	// 应用的主入口，默认调用 cmd.CmdServe 命令
	var rootCmd = &cobra.Command{
		Use:   config.App.GetString("name"),
		Short: "A simple forum project",
		Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

		// rootCmd 的所有子命令都会执行以下代码
		PersistentPreRun: func(command *cobra.Command, args []string) {

		},
	}

	// 注册子命令
	rootCmd.AddCommand(
		cmd.CmdServe,
	)

	// 配置默认运行 Web 服务
	cmd.RegisterDefaultCmd(rootCmd, cmd.CmdServe)

	// 执行主命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Failed to run app with %v: %s", os.Args, err.Error())
	}
}
