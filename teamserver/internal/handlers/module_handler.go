package handlers

import (
	"database/sql"
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
		models.ResponseError(ctx, http.StatusInternalServerError, "Error creating module upload path", err.Error())
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Error getting form file", err.Error())
		return
	}
	defer file.Close()

	filename := header.Filename
	moduleName := strings.TrimSuffix(filename, filepath.Ext(filename))
	modulePath := filepath.Join(moduleUploadPath, moduleName)

	outPath := filepath.Join(moduleUploadPath, filename)
	if err := saveUploadedFile(file, outPath); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Error saving uploaded file", err.Error())
		return
	}

	if err := utils.Unzip(outPath, modulePath); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to unzip file", err.Error())
		return
	}

	if err := os.Remove(outPath); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete zip file", err.Error())
		return
	}

	moduleConfigPath := filepath.Join(modulePath, moduleName+".json")
	moduleConfig, err := LoadModuleConfig(moduleConfigPath)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to load module configuration", err.Error())
		return
	}

	moduleConfig.Module.Id = uuid.New().String()
	if err := mc.ModuleDAL.CreateModule(ctx.Request.Context(), &moduleConfig.Module); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create module", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Module uploaded successfully", nil)
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
	moduleId := ctx.Param(models.ParamModuleID)
	if moduleId == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamModuleID))
		return
	}

	module, err := mc.ModuleDAL.GetModuleById(ctx.Request.Context(), moduleId)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get module", err.Error())
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		models.ResponseError(ctx, http.StatusNotFound, "Module not found", nil)
		return
	}

	moduleUploadPath := conf.GetModuleUploadPath()
	modulePath := filepath.Join(moduleUploadPath, module.Name)

	if err := os.RemoveAll(modulePath); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete module files", err.Error())
		return
	}

	if err := mc.ModuleDAL.DeleteModule(ctx.Request.Context(), moduleId); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete module", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Module deleted successfully", nil)
}

func (mc *ModuleController) DeleteAllModules(ctx *gin.Context) {
	moduleUploadPath := conf.GetModuleUploadPath()

	if err := mc.ModuleDAL.DeleteAllModules(ctx.Request.Context()); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete all modules", err.Error())
		return
	}

	if err := os.RemoveAll(moduleUploadPath); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete all module files", err.Error())
		return
	}

	if err := os.MkdirAll(moduleUploadPath, os.ModePerm); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to recreate module upload path", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "All modules deleted successfully", nil)
}

func (mc *ModuleController) GetAllModules(ctx *gin.Context) {
	modules, err := mc.ModuleDAL.GetAllModules(ctx.Request.Context())
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get all modules", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Modules retrieved successfully", modules)
}

func (mc *ModuleController) GetModuleById(ctx *gin.Context) {
	moduleId := ctx.Param(models.ParamModuleID)
	module, err := mc.ModuleDAL.GetModuleById(ctx.Request.Context(), moduleId)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get module", err.Error())
		return
	}
	if errors.Is(err, sql.ErrNoRows) {
		models.ResponseError(ctx, http.StatusNotFound, "Module not found", nil)
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Module retrieved successfully", module)
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
