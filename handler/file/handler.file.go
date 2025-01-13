package fileHandler

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/levensspel/go-gin-template/domain"
	"github.com/levensspel/go-gin-template/dto"
	"github.com/levensspel/go-gin-template/helper"
	"github.com/levensspel/go-gin-template/logger"
	fileService "github.com/levensspel/go-gin-template/service/file"
	"github.com/samber/do/v2"
)

type FileHandler interface {
	Upload(ctx *gin.Context)
}

type handler struct {
	service       fileService.FileService
	logger        logger.Logger
	storageClient domain.StorageClient
}

func NewHandler(service fileService.FileService, logger logger.Logger, storageClient domain.StorageClient) FileHandler {
	return &handler{service: service, logger: logger, storageClient: storageClient}
}

func NewInject(i do.Injector) (FileHandler, error) {
	_service := do.MustInvoke[fileService.FileService](i)
	_logger := do.MustInvoke[logger.LogHandler](i)
	_storageClient := do.MustInvoke[domain.StorageClient](i)
	return NewHandler(_service, &_logger, _storageClient), nil
}

// Upload godoc
// @Tags file
// @Summary Upload an file
// @Description Upload an file
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer JWT token"
// @Param File formData file true "Body with file zip"
// @Success 200 {object} helper.Response{data=dto.FileUploadRespondPayload} "File uploaded successfully"
// @Success 201 {object} helper.Response{data=dto.FileUploadRespondPayload} "File created successfully"
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

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		h.logger.Warn(err.Error(), helper.FunctionCaller("FileHandler.Upload"), header)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle file."})
		return
	}

	// note filename
	// fileExt := filepath.Ext(header.Filename)
	// originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	// now := time.Now()
	// filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
	// note;end

	// Validasi ekstensi file
	validExtensions := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		".jpeg":      true,
		".jpg":       true,
		".png":       true,
	}

	if !validExtensions[header.Header.Get("Content-Type")] {
		h.logger.Warn("Invalid file type. Only jpeg, jpg, or png are allowed. Try another method", helper.FunctionCaller("FileHandler.Upload"), header.Header.Get("Content-Type"))
		if !validExtensions[filepath.Ext(header.Filename)] {
			h.logger.Warn("Invalid file type. Only jpeg, jpg, or png are allowed.", helper.FunctionCaller("FileHandler.Upload"), filepath.Ext(header.Filename))
			ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "Invalid file type. Only jpeg, jpg, or png are allowed."})
			return
		}
	}

	// Validasi ukuran file (max 100 KiB)
	if header.Size > 1024*100 { // 100 KiB
		h.logger.Warn("File size exceeds 100KiB.", helper.FunctionCaller("FileHandler.Upload"), header.Size)
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "error": "File size exceeds 100KiB."})
		return
	}

	// Read file content into a buffer
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Simpan file atau proses lebih lanjut
	// Misalnya, simpan file ke server atau S3
	h.service.Upload(ctx, file, header)

	uploadedURL, err := h.storageClient.PutFile(ctx, header.Filename, header.Header.Get("Content-Type"), buf.Bytes(), true)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// return file URI kembali dari service
	ctx.JSON(
		http.StatusOK,
		dto.FileUploadRespondPayload{
			Uri: uploadedURL},
	)
}
