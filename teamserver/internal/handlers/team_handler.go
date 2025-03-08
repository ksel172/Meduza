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
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	creatorID := ctx.GetString("userID")
	if creatorID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required field", "userID is required")
		return
	}

	if err := tc.dal.CreateTeam(ctx.Request.Context(), &team, creatorID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create team", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Team created successfully", team)
}

func (tc *TeamController) UpdateTeam(ctx *gin.Context) {
	var team models.Team
	if err := ctx.ShouldBindJSON(&team); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := tc.dal.UpdateTeam(ctx.Request.Context(), &team); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to update team", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Team updated successfully", team)
}

func (tc *TeamController) DeleteTeam(ctx *gin.Context) {
	teamID := ctx.Param(models.ParamTeamID)
	if teamID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamTeamID))
		return
	}

	if err := tc.dal.DeleteTeam(ctx.Request.Context(), teamID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete team", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Team deleted successfully", nil)
}

func (tc *TeamController) GetTeams(ctx *gin.Context) {
	teams, err := tc.dal.GetTeams(ctx.Request.Context())
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get teams", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Teams retrieved successfully", teams)
}

func (tc *TeamController) AddTeamMember(ctx *gin.Context) {
	var member models.TeamMember
	if err := ctx.ShouldBindJSON(&member); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	member.AddedAt = time.Now()

	if err := tc.dal.AddTeamMember(ctx.Request.Context(), &member); err != nil {
		if err.Error() == fmt.Sprintf("user '%s' is already a member of team '%s'", member.UserID, member.TeamID) {
			models.ResponseError(ctx, http.StatusConflict, "Member already exists", err.Error())
		} else {
			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to add team member", err.Error())
		}
		return
	}

	models.ResponseSuccess(ctx, http.StatusCreated, "Team member added successfully", member)
}

func (tc *TeamController) RemoveTeamMember(ctx *gin.Context) {
	memberID := ctx.Param(models.ParamMemberID)
	if memberID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamMemberID))
		return
	}

	if err := tc.dal.RemoveTeamMember(ctx.Request.Context(), memberID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to remove team member", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Team member removed successfully", nil)
}

func (tc *TeamController) GetTeamMembers(ctx *gin.Context) {
	teamID := ctx.Param(models.ParamTeamID)
	if teamID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing required parameter", fmt.Sprintf("%s is required", models.ParamTeamID))
		return
	}

	members, err := tc.dal.GetTeamMembers(ctx.Request.Context(), teamID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get team members", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Team members retrieved successfully", members)
}
