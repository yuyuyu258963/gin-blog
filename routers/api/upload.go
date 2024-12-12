package api

import (
	"gin_example/pkg/e"
	"gin_example/pkg/logging"
	"gin_example/pkg/upload"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
)

// @Summary UploadFile
// @Description 上传图片
// @Tags file
// @Accept json
// @Produce json
// @Success 200  {string} string "ok"
// @Router /api/upload [get]
func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]interface{})

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		goto repl
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := path.Join(fullPath, imageName)
		// 检查文件的格式和文件大小
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				// 尝试保存到本地文件失败
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = path.Join(savePath, imageName)
			}
		}
	}

repl:
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
