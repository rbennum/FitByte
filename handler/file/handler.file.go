package fileHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	service "github.com/levensspel/go-gin-template/service/user"
)

type FileHandler interface {
	Upload(ctx *gin.Context)
}

type handler struct {
	service service.UserService
	logger  logger.Logger
}

func NewUserHandler(service service.UserService, logger logger.Logger) FileHandler {
	return &handler{service: service, logger: logger}
}

// Upload godoc
// @Tags auth
// @Summary Entry for authentication or create new user
// @Description Either create or login
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param data body dto.FileUploadRequestPayload true "File upload data"
// @Success 200 {object} helper.Response{data=helper.FileUploadRespondPayload} "File uploaded successfully"
// @Success 201 {object} helper.Response{data=helper.FileUploadRespondPayload} "File created successfully"
// @Failure 400 {object} helper.Response{errors=helper.ErrorResponse} "Bad Request"
// @Failure 415 {object} helper.Response{errors=helper.ErrorResponse} "Unsupported Media Type"
// @Failure 413 {object} helper.Response{errors=helper.ErrorResponse} "Payload Too Large"
// @Failure 401 {object} helper.Response{errors=helper.ErrorResponse} "Unauthorized - Missing or invalid token"
// @Router /v1/file [POST]
func (h handler) Upload(ctx *gin.Context) {
	// Verify JWT token (pseudo-code, adjust as per your implementation)
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Token validation logic here (e.g., decode, validate)
	// If invalid, return unauthorized response

	input := new(dto.FileUploadRequestPayload)

	if err := ctx.ShouldBindJSON(&input); err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("FileHandler.Upload"), &input)
		ctx.JSON(http.StatusBadRequest, helper.NewResponse(nil, err))
		return
	}

	// Validasi ekstensi file
	validExtensions := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
	}

	if !validExtensions[input.File.Header.Get("Content-Type")] {
		h.logger.Warn("Invalid file type. Only jpeg, jpg, or png are allowed.", helper.FunctionCaller("FileHandler.Upload"), input)
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 415, "error": "Invalid file type. Only jpeg, jpg, or png are allowed."})
		return
	}

	// Validasi ukuran file (max 100 KiB)
	if input.File.Size > 1024*100 { // 100 KiB
		h.logger.Warn("File size exceeds 100KiB.", helper.FunctionCaller("FileHandler.Upload"), input)
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 413, "error": "File size exceeds 100KiB."})
		return
	}

	// Simpan file atau proses lebih lanjut
	// Misalnya, simpan file ke server atau S3

	// return file URI kembali dari service
	ctx.JSON(http.StatusOK, helper.NewResponse(dto.FileUploadRespondPayload{
		Uri: "cuy",
	}, nil))
}
