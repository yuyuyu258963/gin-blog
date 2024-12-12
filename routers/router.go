package routers

import (
	"gin_example/middleware"
	"gin_example/pkg/export"
	"gin_example/pkg/setting"
	"gin_example/pkg/upload"
	api "gin_example/routers/api"
	v1 "gin_example/routers/api/v1"
	"net/http"

	_ "gin_example/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 创建一个新的gin.Engine的实例
func InitRouter() *gin.Engine {
	r := gin.New()
	// 使用两个默认的中间件
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	// 文件服务器的开启，即可以让用户访问本地的文件
	r.StaticFS("/"+setting.AppSetting.ImageSavePath,
		http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/"+setting.AppSetting.ExportSavePath,
		http.Dir(export.GetExcelFullPath()))

	// url := ginSwagger.URL("/swagger/doc.json") // 指定 swagger json 的路径
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/api/auth", api.GetAuth)
	r.GET("/api/upload", api.UploadImage) // 文件上传接口

	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JWT()) // 使用中间件
	{
		// =====  TAG ====
		tagApi(apiv1)
		// =====  Article ====
		articleApi(apiv1)
	}

	return r
}

// 注册tag相关的处理接口
func tagApi(r *gin.RouterGroup) {
	// 获取标签列表
	r.GET("/tags", v1.GetTags)
	// 新增标签
	r.POST("/tags", v1.AddTag)
	// 更新指定标签
	r.PUT("/tags/:id", v1.EditTag)
	// 删除指定标签
	r.DELETE("/tags/:id", v1.DeleteTag)
	// 导出标签
	r.POST("/tags/export", v1.ExportTag)
	// 导入标签
	r.POST("/tags/import", v1.ImportTag)
}

// 注册article相关的处理接口
func articleApi(r *gin.RouterGroup) {
	// 获取指定id的文章
	r.GET("/article/:id", v1.GetArticle)
	// 获取文章列表
	r.GET("/articles", v1.GetArticles)
	// 新增文章
	r.POST("/article", v1.AddArticle)
	// 修改文章
	r.PUT("/article/:id", v1.EditArticle)
	// 删除文章
	r.DELETE("/article/:id", v1.DeleteArticle)
}
