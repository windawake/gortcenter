package cmd

import (
	"fmt"
	"gortcenter/pkg/config"
	"gortcenter/routes"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

// CmdServe represents the available web sub-command.
var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs,
}

func runWeb(cmd *cobra.Command, args []string) {

	// 设置 gin 的运行模式，支持 debug, release, test
	// release 会屏蔽调试信息，官方建议生产环境中使用
	// 非 release 模式 gin 终端打印太多信息，干扰到我们程序中的 Log
	// 故此设置为 release，有特殊情况手动改为 debug 即可

	mode := gin.ReleaseMode
	if config.App.GetBool("debug") {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)

	// gin 实例
	router := gin.New()

	//  注册 API 路由
	routes.RegisterAPIRoutes(router)

	// 运行服务器
	err := router.Run(":" + config.App.GetString("port"))

	if err != nil {
		fmt.Println(err)
	}
}
