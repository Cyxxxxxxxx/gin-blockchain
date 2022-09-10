package routes

import (
	"blockchain/controller"
	logger "blockchain/logs"
	"blockchain/middlewares"

	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"

	// "github.com/swaggo/gin-swagger/swaggerFiles"
	swaggerFiles "github.com/swaggo/files"
)

func Setup(mode string) *gin.Engine {
	//如果设置mode为release则设置gin为该模式
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middlewares.Cors())
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	r.StaticFile("/WW_verify_CSaj5WtuBq3AemhM.txt", "./WW_verify_CSaj5WtuBq3AemhM.txt")

	/* midgroups := r.Group("/api").Use(middleware.Cors())
	{
		groups.Get("/.."),controller.XXX()
	} */

	groups := r.Group("/api")
	{
		groups.GET("/GetBlockchain", controller.GetBlockchain)
		groups.POST("/WriteBlock", controller.WriteBlock)
	}

	return r
}
