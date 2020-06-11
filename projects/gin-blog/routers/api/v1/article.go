package v1

import (
	"gin-blog/pkg/app"
	"gin-blog/pkg/codec"
	"gin-blog/pkg/qrcode"
	"gin-blog/pkg/setting"
	"gin-blog/pkg/util"
	"gin-blog/service/article_service"
	"gin-blog/service/tag_service"
	"github.com/astaxie/beego/validation"
	"github.com/boombuler/barcode/qr"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
	"net/http"
)
// @Summary Get a single article 获取单个文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response 
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors(){
		app.MarkErrors(valid.Errors)
		app.ResponseFunc(c,http.StatusOK, codec.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_GET_ARTICLE_FAIL, article)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, article)
 
}


// @Summary Get multiple articles 获取多个文章
// @Produce  json
// @Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response 
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	valid := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id")
	}

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.ResponseFunc(c,http.StatusBadRequest, codec.INVALID_PARAMS, nil)
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
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, data)
}




type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add article 新增文章
// @Produce  json
// @Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response 
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	var (
		form AddArticleForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != codec.SUCCESS {
		app.ResponseFunc(c,httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		CreatedBy:     form.CreatedBy,
	}
	if err := articleService.Add(); err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}



type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Update article 修改文章
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} app.Response 
// @Failure 500 {object} app.Response 
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	var (

		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != codec.SUCCESS {
		app.ResponseFunc(c,httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		ModifiedBy:    form.ModifiedBy,
		State:         form.State,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}

// @Summary Delete article 删除文章
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response 
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		app.ResponseFunc(c,http.StatusOK, codec.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		app.ResponseFunc(c,http.StatusOK, codec.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, nil)
}

const (
	QRCODE_URL = "https://github.com/gyb333/practicego"
)

func GenerateArticlePoster(c *gin.Context) {

	article := &article_service.Article{}
	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := article_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	articlePoster := article_service.NewArticlePoster(posterName, article, qr)
	articlePosterBgService := article_service.NewArticlePosterBg(
		"bg.jpg",
		articlePoster,
		&article_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&article_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := articlePosterBgService.Generate()
	if err != nil {
		app.ResponseFunc(c,http.StatusInternalServerError, codec.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	app.ResponseFunc(c,http.StatusOK, codec.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}