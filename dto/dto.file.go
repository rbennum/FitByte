package dto

import "mime/multipart"

type FileUploadRequestPayload struct {
	File     *multipart.FileHeader `form:"file" binding:"required"`
	FileType string                `form:"filetype binding:"required"`
}

type FileUploadRespondPayload struct {
	Uri string `form:"uri"`
}
