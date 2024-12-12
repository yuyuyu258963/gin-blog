package e

var MsgFlags = map[int]string{
	SUCCESS:              "ok",
	ERROR:                "fail",
	INVALID_PARAMS:       "请求参数错误",
	ERROR_NOTFOUND_TOKEN: "未找到Token",

	ERROR_EXIST_TAG:        "已存在该名称的标签",
	ERROR_NOT_EXITS_TAG:    "该标签不存在",
	ERROR_CHECK_TAG_EXISTS: "检查标签是否存在失败",
	ERROR_EXPORT_TAG_FAIL:  "导出标签失败",
	ERROR_IMPORT_TAG_FAIL:  "导入标签失败",

	ERROR_NOT_EXIST_ARTICLE:    "该文章不存在",
	ERROR_GET_ARTICLE_FAIL:     "获取文章失败",
	ERROR_CHECK_ARTICLE_EXISTS: "检查文件是否存在失败",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",

	ERROR_UPLOAD_SAVE_IMAGE_FAIL:  "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL: "检查图片失败",
	ERROR_UPLOAD_IMAGE_FORMAT:     "校验图片错误，图片格式或大小有问题",
}

// 对外暴露错误编码到错误信息的映射
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
