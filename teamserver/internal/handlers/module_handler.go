package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ModuleController struct {
}

func NewModuleController() *ModuleController {
	return &ModuleController{}
}

func (mc *ModuleController) UploadModule(ctx *gin.Context) {

	moduleUploadPath := conf.GetModuleUploadPath()

	// Create the module upload path if it doesn't exist
	if _, err := os.Stat(moduleUploadPath); os.IsNotExist(err) {
		err = os.MkdirAll(moduleUploadPath, os.ModePerm)
		if err != nil {
			ctx.String(http.StatusInternalServerError, fmt.Sprintf("create module upload path err: %s", err.Error()))
			return
		}
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	defer file.Close()

	filename := header.Filename
	moduleName := strings.TrimSuffix(filename, filepath.Ext(filename))
	modulePath := filepath.Join(moduleUploadPath, moduleName)

	outPath := filepath.Join(moduleUploadPath, filename)
	outFile, err := os.Create(outPath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("create file err: %s", err.Error()))
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, file)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("write file err: %s", err.Error()))
		return
	}

	err = utils.Unzip(outPath, modulePath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("unzip file err: %s", err.Error()))
		return
	}

	err = os.Remove(outPath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("delete zip file err: %s", err.Error()))
		return
	}

	ctx.String(http.StatusOK, "upload successful")
}

func (mc *ModuleController) DeleteModule(ctx *gin.Context) {

}
