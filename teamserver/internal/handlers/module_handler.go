package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ModuleController struct {
	ModuleDAL *dal.ModuleDAL
}

func NewModuleController(moduleDAL *dal.ModuleDAL) *ModuleController {
	return &ModuleController{ModuleDAL: moduleDAL}
}

func (mc *ModuleController) UploadModule(ctx *gin.Context) {
	moduleUploadPath := conf.GetModuleUploadPath()

	// Create the module upload path if it doesn't exist
	if err := os.MkdirAll(moduleUploadPath, os.ModePerm); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("create module upload path err: %s", err.Error()))
		return
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
	if err := saveUploadedFile(file, outPath); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("save uploaded file err: %s", err.Error()))
		return
	}

	if err := utils.Unzip(outPath, modulePath); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("unzip file err: %s", err.Error()))
		return
	}

	if err := os.Remove(outPath); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("delete zip file err: %s", err.Error()))
		return
	}

	moduleConfigPath := filepath.Join(modulePath, moduleName+".json")
	moduleConfig, err := LoadModuleConfig(moduleConfigPath)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("load module config err: %s", err.Error()))
		return
	}

	moduleConfig.Module.Id = uuid.New().String()
	if err := mc.ModuleDAL.CreateModule(ctx.Request.Context(), &moduleConfig.Module); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("save module to database err: %s", err.Error()))
		return
	}

	ctx.String(http.StatusOK, "upload successful")
}

func saveUploadedFile(file io.Reader, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, file); err != nil {
		return err
	}

	return nil
}

func (mc *ModuleController) DeleteModule(ctx *gin.Context) {
	moduleId := ctx.Param("id")
	module, err := mc.ModuleDAL.GetModuleById(ctx.Request.Context(), moduleId)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("get module by id err: %s", err.Error()))
		return
	}
	if module == nil {
		ctx.String(http.StatusNotFound, "module not found")
		return
	}

	moduleUploadPath := conf.GetModuleUploadPath()
	modulePath := filepath.Join(moduleUploadPath, module.Name)

	if err := os.RemoveAll(modulePath); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("delete module files err: %s", err.Error()))
		return
	}

	if err := mc.ModuleDAL.DeleteModule(ctx.Request.Context(), moduleId); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("delete module err: %s", err.Error()))
		return
	}
	ctx.String(http.StatusOK, "delete successful")
}

func (mc *ModuleController) DeleteAllModules(ctx *gin.Context) {
	moduleUploadPath := conf.GetModuleUploadPath()

	if err := mc.ModuleDAL.DeleteAllModules(ctx.Request.Context()); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("delete all modules err: %s", err.Error()))
		return
	}

	if err := os.RemoveAll(moduleUploadPath); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("delete all module files err: %s", err.Error()))
		return
	}

	if err := os.MkdirAll(moduleUploadPath, os.ModePerm); err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("recreate module upload path err: %s", err.Error()))
		return
	}

	ctx.String(http.StatusOK, "delete all successful")
}

func (mc *ModuleController) GetAllModules(ctx *gin.Context) {
	modules, err := mc.ModuleDAL.GetAllModules(ctx.Request.Context())
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("get all modules err: %s", err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, modules)
}

func (mc *ModuleController) GetModuleById(ctx *gin.Context) {
	moduleId := ctx.Param("id")
	module, err := mc.ModuleDAL.GetModuleById(ctx.Request.Context(), moduleId)
	if err != nil {
		ctx.String(http.StatusInternalServerError, fmt.Sprintf("get module by id err: %s", err.Error()))
		return
	}
	if module == nil {
		ctx.String(http.StatusNotFound, "module not found")
		return
	}
	ctx.JSON(http.StatusOK, module)
}

func LoadModuleConfig(filePath string) (*models.ModuleConfig, error) {
	bytes, err := GetModuleBytes(filePath)
	if err != nil {
		return nil, err
	}

	var moduleConfig models.ModuleConfig
	if err := json.Unmarshal(bytes, &moduleConfig); err != nil {
		return nil, err
	}

	return &moduleConfig, nil
}

func GetModuleBytes(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func LoadModule(path string) (*models.ModuleConfig, error) {
	bytes, err := GetModuleBytes(path)
	if err != nil {
		return nil, err
	}

	var moduleConfig models.ModuleConfig
	if err := json.Unmarshal(bytes, &moduleConfig); err != nil {
		return nil, err
	}

	return &moduleConfig, nil
}

func LoadCommands(moduleConfig *models.ModuleConfig) ([]models.Command, error) {
	if moduleConfig == nil {
		return nil, errors.New("moduleConfig is nil")
	}

	return moduleConfig.Module.Commands, nil
}

func LoadAllModules(paths []string) ([]*models.ModuleConfig, error) {
	var modules []*models.ModuleConfig

	for _, path := range paths {
		module, err := LoadModule(path)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}

	return modules, nil
}
