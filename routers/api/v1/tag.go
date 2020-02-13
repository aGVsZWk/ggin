package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/Unknwon/com"
	"ggin/pkg/e"
	"ggin/pkg/util"
	"ggin/pkg/setting"
	"net/http"
	"github.com/astaxie/beego/validation"
	"ggin/pkg/app"
	"ggin/service/tag_service"
)

// 获取多个文章标签
func GetTags(c *gin.Context) {
	appG := app.Gin{c}
	name := c.Query("name")
	var state = -1
	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}
	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

// @Summary 新增文章标签
// @Produce json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query string true "CreatedBy"
// @Success 200 {string} string "{"code": 200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	appG := app.Gin{c}
	name := c.PostForm("name")
	state := com.StrTo(c.DefaultPostForm("state", "0")).MustInt()
	createdBy := c.PostForm("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}
	err = tagService.Add()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 修改文章标签
func EditTag(c *gin.Context) {
	appG := app.Gin{c}

	id := com.StrTo(c.Param("id")).MustInt()
	name := c.PostForm("name")
	modifiedBy := c.PostForm("modified_by")

	valid := validation.Validation{}

	var state = -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{
		ID:         id,
		Name:       name,
		ModifiedBy: modifiedBy,
		State:      state,
	}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// 删除文章标签
func DeleteTag(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	if err := tagService.Delete(); err != nil {
		appG.Response(http.StatusOK, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)

}
