package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/utils"
)

var (
	curAdminCount = 0        // Tracks the number of admins created
	maxAdmin      = 4        // Maximum number of allowed admins
	mu            sync.Mutex // Ensures thread safety for token counter
	s             = utils.Status
)

type AdminController struct {
	dal dal.IAdminDal
}

func NewAdminController(dal dal.IAdminDal) *AdminController {
	return &AdminController{
		dal: dal,
	}
}

func (aC *AdminController) CreateAdmin(ctx *gin.Context) {

	var adminReq models.ResAdmin

	// it checks if route is restricted
	if isRouteRestricted() {
		ctx.JSONP(http.StatusForbidden, gin.H{
			"message": "route permanently restricted: admin creation is not allowed",
			"status":  s.ERROR,
		})
		ctx.Abort()
		return
	}

	authToken := ctx.GetHeader("Authorization")
	if authToken == "" {
		ctx.JSONP(http.StatusForbidden, gin.H{
			"message": "authorization token is missing",
			"status":  s.FAILED,
		})
		ctx.Abort()
		return
	}

	// It's Validate Token
	if err := validateToken(authToken); err != nil {
		remaining := maxAdmin - curAdminCount
		logger.Error("Token Validation Error", err)
		ctx.JSONP(http.StatusUnauthorized, gin.H{
			"message":          "token validation error",
			"status":           s.ERROR,
			"admins_left":      remaining,
			"max_admin_tokens": maxAdmin,
			"token":            authToken,
		})
		ctx.Abort()
		return
	}

	// parsing request body
	if err := ctx.ShouldBindJSON(&adminReq); err != nil {
		logger.Error("Unable to parse error of request body", err)
		ctx.JSONP(http.StatusUnprocessableEntity, gin.H{
			"message": "unable to parse request body",
			"status":  s.FAILED,
		})
		ctx.Abort()
		return
	}

	hashPassword, err := utils.HashPassword(adminReq.PasswordHash)
	if err != nil {
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"message": "invalid request error",
			"status":  s.ERROR,
		})
		ctx.Abort()
		return
	}
	adminReq.PasswordHash = hashPassword

	validate := utils.NewValidatorService()
	if err := validate.ValidateStruct(adminReq); err != nil {
		logger.Error("Validation Parsing Failed :", err)
		ctx.JSONP(http.StatusBadRequest, gin.H{
			"message": "validation failed",
			"status":  s.FAILED,
		})
		ctx.Abort()
		return
	}

	err = aC.dal.CreateDefaultAdmins(ctx.Request.Context(), &adminReq)
	if err != nil {
		logger.Error("Error while creating/adding admin: ", err)
		ctx.JSONP(http.StatusInternalServerError, gin.H{
			"message": "unable to add admin",
			"status":  s.FAILED,
		})
		ctx.Abort()
		return
	}

	adminCount()

	remaining := maxAdmin - curAdminCount

	ctx.JSONP(http.StatusCreated, gin.H{
		"message":          "admin created successfully",
		"status":           s.SUCCESS,
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
	reqTokenSplit := strings.Split(reqToken, " ")
	if len(reqTokenSplit) != 2 {
		return fmt.Errorf("invalid token")
	}
	reqToken = reqTokenSplit[1] // Grab only the token, ignore Bearer text

	envToken := conf.GetMeduzaAdminSecret()
	if envToken == "" {
		return fmt.Errorf("error loading admin secret, empty")
	}

	if reqToken != envToken {
		return fmt.Errorf("invalid or expired token")
	}

	return nil
}
