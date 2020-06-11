package api

import (
	"gin-blog/pkg/codec"
	"gin-blog/pkg/logging"
	"gin-blog/pkg/upload"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadImage(c *gin.Context) {
	code := codec.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = codec.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  codec.GetMsg(code),
			"data": data,
		})
	}

	if image == nil {
		code = codec.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if ! upload.CheckImageExt(imageName) || ! upload.CheckImageSize(file) {
			code = codec.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = codec.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = codec.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  codec.GetMsg(code),
		"data": data,
	})
}
