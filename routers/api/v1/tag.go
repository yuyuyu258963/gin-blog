package v1

import (
	"gin_example/models"
	"gin_example/pkg/app"
	"gin_example/pkg/e"
	"gin_example/pkg/export"
	"gin_example/pkg/logging"
	"gin_example/pkg/setting"
	"gin_example/pkg/util"
	"gin_example/service/tagService"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// 获取多个文章标签
// @Summary 获取多个文章标签
// @Tags Tag
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [Get]
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS
	// util.GetPage 保证了各接口处理page的逻辑是一致的
	data["list"], _ = models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"], _ = models.GetTagTotal(maps)

	appG.Response(http.StatusOK, code, data)
}

// TODO 涉及幂等性校验的部分需要用事务实现

// @Summary 新增文章标签
// @Tags Tag
// @Produce  json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	// 参数验证
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, make(map[string]interface{}))
		return
	}
	exists, err := models.ExistTagByName(name)
	if err != nil {
		code = e.ERROR_CHECK_TAG_EXISTS
		goto reply
	}
	if !exists {
		code = e.ERROR_EXIST_TAG
		goto reply
	}
	code = e.SUCCESS
	models.AddTag(name, state, createdBy)

reply:
	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

// @Summary 修改文章标签
// @Produce  json
// @Tags Tag
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许为0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, make(map[string]interface{}))
		return
	}

	data := make(map[string]interface{})
	tagService := tagService.Tag{ID: id}
	exists, err := tagService.ExistTagByID()
	if err != nil {
		code = e.ERROR_CHECK_TAG_EXISTS
		goto reply
	}
	if !exists {
		code = e.ERROR_NOT_EXITS_TAG
		goto reply
	}

	code = e.SUCCESS
	data["modified_by"] = modifiedBy
	if name != "" {
		data["name"] = name
	}
	if state != -1 {
		data["state"] = state
	}

	models.EditTag(id, data)

reply:
	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

// 删除文章标签
// @Summary 删除文章标签
// @Tags Tag
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [Delete]
func DeleteTag(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, make(map[string]interface{}))
	}
	tagService := tagService.Tag{ID: id}
	exists, err := tagService.ExistTagByID()
	if err != nil {
		code = e.ERROR_CHECK_TAG_EXISTS
		goto reply
	}
	if !exists {
		code = e.ERROR_NOT_EXITS_TAG
		goto reply
	}

	code = e.SUCCESS
	models.DeleteTag(id)

reply:
	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

// 导出标签
// @Summary 导出标签
// @Tags Tag
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /tags/export [post]
func ExportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := &tagService.Tag{
		Name:  name,
		State: state,
		// PageNum:  setting.PageSize,
		PageSize: setting.AppSetting.PageSize,
	}

	filename, err := tagService.Export()
	if err != nil {
		logging.InfoFiled(logging.LogFields{"error": err})
		appG.Response(http.StatusOK, e.ERROR_EXPORT_TAG_FAIL, err)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelFullPath() + filename,
	})
}

// 导入标签
// @Summary 导入标签
// @Tags Tag
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /tags/import [post]
func ImportTag(c *gin.Context) {
	appG := app.Gin{C: c}

	formFile, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	tagS := tagService.Tag{}
	err = tagS.Import(formFile)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
