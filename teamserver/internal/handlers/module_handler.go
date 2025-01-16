package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ModuleController struct {
}

func NewModuleController() *ModuleController {
	return &ModuleController{}
}

func (mc *ModuleController) UploadModule(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	defer file.Close()

	filename := header.Filename
	moduleName := strings.TrimSuffix(filename, filepath.Ext(filename))
	modulePath := filepath.Join("./teamserver/modules", moduleName)

	outPath := filepath.Join("./teamserver/modules", filename)
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

	ctx.String(http.StatusOK, "upload successful")
}

func (mc *ModuleController) DeleteModule(ctx *gin.Context) {

}
