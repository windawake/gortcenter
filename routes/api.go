// Package routes 注册路由
package routes

import (
	"gortcenter/app/http/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	var v1 *gin.RouterGroup
	v1 = r.Group("/api")
	rt := new(controllers.ResetTransController)

	v1.POST("/resetTransaction/commit", rt.Commit)
	v1.POST("/resetTransaction/rollback", rt.Rollback)
}
