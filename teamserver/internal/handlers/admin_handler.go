package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var (
	curAdminCount = 0        // Tracks the number of admins created
	maxAdmin      = 4        // Maximum number of allowed admins
	mu            sync.Mutex // Ensures thread safety for token counter
)

type AdminController struct {
	dal *dal.AdminDAL
}

func NewAdminController(dal *dal.AdminDAL) *AdminController {
	return &AdminController{
		dal: dal,
	}
}

func (aC *AdminController) CreateAdmin(ctx *gin.Context) {

	var adminReq models.ResAdmin

	// it checks if route is restricted
	if isRouteRestricted() {
		ctx.JSONP(http.StatusForbidden, gin.H{
			"Message": "Route Permanently restricted: Admin creation is not allowed",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	authToken := ctx.GetHeader("Authorization")
	if authToken == "" {
		ctx.JSONP(http.StatusForbidden, gin.H{
			"Message": "Authorization token is missing",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	// It's Validate Token
	if err := validateToken(authToken); err != nil {
		remaining := maxAdmin - curAdminCount
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"Error":            err.Error(),
			"Message":          "Token Validation Error",
			"Status":           "Failed",
			"admins_left":      remaining,
			"max_admin_tokens": maxAdmin,
			"token":            authToken,
		})
		ctx.Abort()
		return
	}

	// parsing request body
	if err := ctx.ShouldBindJSON(&adminReq); err != nil {
		ctx.JSONP(http.StatusUnprocessableEntity, gin.H{
			"Error":   err.Error(),
			"Message": "Unable to parse request body",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	hashPassword, err := utils.HashPassword(adminReq.PasswordHash)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Invalid Request Error",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}
	adminReq.PasswordHash = hashPassword

	validate := utils.NewValidatorService()
	if err := validate.ValidateStruct(adminReq); err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"Message": "Validation failed",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	err = aC.dal.CreateDefaultAdmins(ctx.Request.Context(), &adminReq)
	if err != nil {
		ctx.JSONP(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"Message": "Error Adding Users",
			"Status":  "Failed",
		})
		ctx.Abort()
		return
	}

	adminCount()

	remaining := maxAdmin - curAdminCount

	ctx.JSONP(http.StatusCreated, gin.H{
		"Message":          "User Created Successfully",
		"Status":           "Success",
		"admins_left":      remaining,
		"max_admin_tokens": maxAdmin,
	})
}

func isRouteRestricted() bool {
	return curAdminCount >= maxAdmin
}

func adminCount() {
	mu.Lock()
	curAdminCount++
	mu.Unlock()
}

func validateToken(reqToken string) error {
	envToken := conf.GetMeduzaAdminSecret()
	if envToken != "" {
		return fmt.Errorf("error loading admin secret, empty")
	}

	if reqToken != envToken {
		return fmt.Errorf("invalid or expired token")
	}

	return nil
}
