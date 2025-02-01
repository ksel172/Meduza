package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

type TeamController struct {
	dal dal.ITeamDAL
}

func NewTeamController(dal dal.ITeamDAL) *TeamController {
	return &TeamController{dal: dal}
}

func (tc *TeamController) CreateTeam(ctx *gin.Context) {
	var team models.Team
	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.dal.CreateTeam(ctx.Request.Context(), &team); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, team)
}

func (tc *TeamController) UpdateTeam(ctx *gin.Context) {
	var team models.Team
	if err := ctx.ShouldBindJSON(&team); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tc.dal.UpdateTeam(ctx.Request.Context(), &team); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, team)
}

func (tc *TeamController) DeleteTeam(ctx *gin.Context) {
	teamID := ctx.Param("id")
	if err := tc.dal.DeleteTeam(ctx.Request.Context(), teamID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "team deleted"})
}

func (tc *TeamController) GetTeams(ctx *gin.Context) {
	teams, err := tc.dal.GetTeams(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, teams)
}

func (tc *TeamController) AddTeamMember(ctx *gin.Context) {
	var member models.TeamMember
	if err := ctx.ShouldBindJSON(&member); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member.AddedAt = time.Now()

	if err := tc.dal.AddTeamMember(ctx.Request.Context(), &member); err != nil {
		if err.Error() == fmt.Sprintf("user '%s' is already a member of team '%s'", member.UserID, member.TeamID) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusCreated, member)
}

func (tc *TeamController) RemoveTeamMember(ctx *gin.Context) {
	memberID := ctx.Param("id")
	if err := tc.dal.RemoveTeamMember(ctx.Request.Context(), memberID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "team member removed"})
}

func (tc *TeamController) GetTeamMembers(ctx *gin.Context) {
	teamID := ctx.Param("id")
	members, err := tc.dal.GetTeamMembers(ctx.Request.Context(), teamID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, members)
}
