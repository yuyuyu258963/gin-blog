package v1

import (
	"gin_example/models"
	"gin_example/pkg/app"
	"gin_example/pkg/e"
	"gin_example/pkg/setting"
	"gin_example/pkg/util"
	"gin_example/service/articleService"
	"gin_example/service/tagService"
	"net/http"
	"strings"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// TODO 1.涉及幂等性校验的部分需要用事务实现
// TODO 2.如果用到了Redis那就要考虑将所有的操作都迁移使用Redis否则将导致数据不一致

// 获取单个文章
// @Summary GetArticles List
// @Description Get aim articles
// @Tags Article
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Failure 10003 {string} json "{"code":10003,"data":{},"msg":"文章不存在"}"
// @Router /api/v1/article/{id} [get]
func GetArticle(c *gin.Context) {
	appG := &app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()

	code := e.INVALID_PARAMS
	var data interface{}
	var err error
	var exists bool
	var articleS articleService.Article
	// 参数提前校验
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID 必须大于0")

	// 参数校验失败
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}

	articleS = articleService.Article{ID: id}
	exists, err = articleS.ExistArticleByID()
	if err != nil {
		code = e.ERROR_CHECK_ARTICLE_EXISTS
		goto reply
	}
	if !exists {
		code = e.ERROR_NOT_EXIST_ARTICLE
		goto reply
	}

	data, err = articleS.Get()
	if err != nil {
		code = e.ERROR_GET_ARTICLE_FAIL
		goto reply
	}
	code = e.SUCCESS

reply:
	appG.Response(http.StatusOK, code, data)
}

// 获取多个文章
// 支持根据state 或 tag 查询
// @Summary GetArticles List
// @Description Get aim articles
// @Tags Article
// @Accept json
// @Produce json
// @Param state body int false "State"
// @Param tag_id body int false "TagId"
// @Success 200 {string} string "ok"
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	appG := &app.Gin{C: c}
	data := make(map[string]interface{})
	maps := make(map[string]interface{}) // 用于组织查询条件
	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state

		valid.Range(state, 0, 1, "state").Message("状态只允许为0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		maps["tag_id"] = tagId

		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	code := e.INVALID_PARAMS
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		goto reply
	}

	code = e.SUCCESS

	data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
	data["total"] = models.GetArticleTotal(maps)

reply:
	appG.Response(http.StatusOK, code, data)
}

// 新增文章
// @Summary 新增文章
// @Description 新增文章
// @Tags Article
// @Accept json
// @Produce json
// @Param tag_id query int false "TagId"
// @Param title query string true "Title"
// @Param content query string true "Content"
// @Param created_by query string true "CreatedBy"
// @Param state query int false "State"
// @Param cover_image_url query string false "coverImageUrl"
// @Success 200 {string} string "ok"
// @Failure 10002 {string} json "{"code":10002,"data":{},"msg":"Tag不存在"}"
// @Router /api/v1/article [put]
func AddArticle(c *gin.Context) {
	appG := &app.Gin{C: c}

	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(strings.Trim(c.Query("state"), " ")).MustInt()
	coverImageUrl := c.Query("cover_image_url")

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("文件路径最长为255字符")

	data := make(map[string]interface{})
	code := e.INVALID_PARAMS
	// 参数校验失败
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, data)
		return
	}
	// TODO 这里没有ID了是不是就查不了了，只能去查Tag是不是存在
	tagService := tagService.Tag{ID: tagId}
	exist, err := tagService.ExistTagByID()
	if err != nil {
		code = e.ERROR_CHECK_TAG_EXISTS
		goto reply
	}
	if !exist {
		code = e.ERROR_NOT_EXITS_TAG
		goto reply
	}

	data["tag_id"] = tagId
	data["title"] = title
	data["desc"] = desc
	data["content"] = content
	data["created_by"] = createdBy
	data["state"] = state
	data["cover_image_url"] = coverImageUrl
	models.AddArticle(data)
	code = e.SUCCESS

reply:
	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

// 修改文章
// @Summary 修改文章
// @Description 修改文章
// @Tags Article
// @Accept json
// @Produce json
// @Param tag_id query int false "TagId"
// @Param title query string false "Title"
// @Param content query string false "Content"
// @Param modified_by query string true "ModifiedBy"
// @Param cover_image_url query string false "coverImageUrl"
// @Param state query int false "State"
// @Success 200 {string} string "ok"
// @Failure 10003 {string} json "{"code":10003,"data":{},"msg":"文章不存在"}"
// @Router /api/v1/article/{id} [post]
func EditArticle(c *gin.Context) {
	appG := &app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	coverImageUrl := c.Query("cover_image_url")

	// 参数校验
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()

		valid.Range(state, 0, 1, "state").Message("状态只允许为0或1")
	}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题不能超过100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("文件路径最长为255")

	code := e.INVALID_PARAMS

	var exists bool
	var err error
	var articleS articleService.Article
	data := make(map[string]interface{})

	// 参数校验不通过
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		goto reply
	}

	articleS = articleService.Article{ID: id}
	exists, err = articleS.ExistArticleByID()

	if err != nil {
		code = e.ERROR_CHECK_ARTICLE_EXISTS
		goto reply
	}
	if !exists {
		code = e.ERROR_NOT_EXIST_ARTICLE
		goto reply
	}

	if tagId > 0 {
		data["tag_id"] = tagId
	}
	if title != "" {
		data["title"] = title
	}
	if desc != "" {
		data["desc"] = desc
	}
	if content != "" {
		data["content"] = content
	}
	if coverImageUrl != "" {
		data["cover_image_url"] = coverImageUrl
	}
	data["modified_by"] = modifiedBy

	models.EditArticle(id, data)
	code = e.SUCCESS

reply:
	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}

// 删除文章
// @Summary 删除文章
// @Description 删除文章
// @Tags Article
// @Produce  json
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Failure 10003 {string} json "{"code":10003,"data":{},"msg":"文章不存在"}"
// @Router /api/v1/article/{id} [Delete]
func DeleteArticle(c *gin.Context) {
	appG := &app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}

	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, code, make(map[string]interface{}))
		return
	}

	articleS := articleService.Article{ID: id}
	exists, err := articleS.ExistArticleByID()
	if err != nil {
		code = e.ERROR_CHECK_ARTICLE_EXISTS
		goto reply
	}
	if !exists {
		code = e.ERROR_NOT_EXIST_ARTICLE
		goto reply
	}

	models.DeleteArticle(id)
	code = e.SUCCESS

reply:
	appG.Response(http.StatusOK, code, make(map[string]interface{}))
}
