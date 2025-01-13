package dto

import "mime/multipart"

type FileUploadRequestPayload struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

type FileUploadRespondPayload struct {
	Uri string `json:"uri"`
}
