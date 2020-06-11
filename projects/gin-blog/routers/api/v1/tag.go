package v1

import (
	"gin-blog/models"
	"gin-blog/pkg/app"
	"gin-blog/pkg/codec"
	"gin-blog/pkg/export"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)



// @Summary Get multiple article tags 获取多个文章标签
// @Produce  json
// @Param name query string false "Name"
// @Param state query int false "State"
// @Success 200 {object} app.Response 
// @Failure 500 {object} app.Response 
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {

	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}


type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add article tag 新增文章标签
// @Produce  json
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response 
// @Failure 500 {object} app.Response 
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		form AddTagForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != codec.SUCCESS {
		app.ResponseFunc(c,httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}
	exists, err := tagService.ExistByName()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Update article tag 修改文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name body string true "Name"
// @Param state body int false "State"
// @Param modified_by body string true "ModifiedBy"
// @Success 200 {object} app.Response 
// @Failure 500 {object} app.Response 
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != codec.SUCCESS {
		app.ResponseFunc(c,httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exists, err := tagService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}

// @Summary Delete article tag
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response 
// @Failure 500 {object} app.Response 
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.ResponseFunc(c,http.StatusBadRequest, codec.INVALID_PARAMS, nil)
	}

	tagService := tag_service.Tag{ID: id}
	exists, err := tagService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	if err := tagService.Delete(); err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}

// @Summary Export article tag
// @Produce  json
// @Param name body string false "Name"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/export [post]
func ExportTag(c *gin.Context) {
	name := c.PostForm("name")
	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// @Summary Import article tag
// @Produce  json
// @Param file body file true "Excel File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func ImportTag(c *gin.Context) {

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}



// @Summary 修改文章标签
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func EditTagOld(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := codec.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = codec.SUCCESS
		if flag,_:=models.ExistTagByID(id);flag {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			code = codec.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : codec.GetMsg(code),
		"data" : make(map[string]string),
	})
}



