package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/astaxie/beego/validation"
	"github.com/Unknwon/com"
	"ggin/pkg/e"
	"net/http"
	//"ggin/pkg/util"
	//"ggin/pkg/setting"
	"ggin/pkg/app"
	"ggin/service/article_service"
	"ggin/pkg/setting"
	"ggin/pkg/util"
	"ggin/service/tag_service"
)

// 获取单个文章
func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	article, err := articleService.Get()

	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

// 获取多个文章
func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}
	var state = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	var tagId = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 0, "tag_id").Message("标签ID必须大于0")
	}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}
	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

// 新增文章
func AddArticle(c *gin.Context) {
	appG := app.Gin{c}
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	createdBy := c.PostForm("created_by")
	coverImageUrl := c.PostForm("cover_image_url")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Required(coverImageUrl, "cover_image_url").Message("封面地址不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{ID: tagId}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         tagId,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CoverImageUrl: coverImageUrl,
		State:         state,
		CreateBy:      createdBy,
	}

	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章
func EditArticle(c *gin.Context) {
	appG := app.Gin{c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.PostForm("tag_id")).MustInt()
	title := c.PostForm("title")
	desc := c.PostForm("desc")
	content := c.PostForm("content")
	coverImageUrl := c.PostForm("cover_image_url")
	modifiedBy := c.PostForm("modified_by")

	var state = - 1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(coverImageUrl, "cover_image_url").Message("封面地址不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(coverImageUrl, 255, "cover_image_url").Message("封面地址最长255字符")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{
		ID:            id,
		Title:         title,
		Desc:          desc,
		Content:       content,
		CoverImageUrl: coverImageUrl,
		ModifiedBy:    modifiedBy,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	tagService := tag_service.Tag{ID : tagId}
	exists, err = tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
}

// 删除文章
func DelelteArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	err = articleService.Delete()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
