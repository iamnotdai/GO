package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"team-service/dtos"
	"team-service/services"
)

type teamHandler struct {
	service services.ITeamService
}

func TeamHandler(service services.ITeamService) *teamHandler {
	return &teamHandler{service}
}

// --- Routes binding ---

func (h *teamHandler) RegisterRoutes(r *gin.RouterGroup) {
	teams := r.Group("/teams")
	{
		teams.POST("", h.CreateTeam)                                  // POST /teams
		teams.POST("/:teamId/members", h.AddMember)                   // POST /teams/:teamId/members
		teams.DELETE("/:teamId/members/:memberId", h.RemoveMember)    // DELETE /teams/:teamId/members/:memberId
		teams.POST("/:teamId/managers", h.AddManager)                 // POST /teams/:teamId/managers
		teams.DELETE("/:teamId/managers/:managerId", h.RemoveManager) // DELETE /teams/:teamId/managers/:managerId
	}
}

// CreateTeam godoc
// @Summary Tạo team mới
// @Tags Teams
// @Accept json
// @Produce json
// @Param data body dtos.CreateTeamInput true "Team data"
// @Success 201 {object} map[string]interface{}
// @Failure 400,500 {object} map[string]string
// @Router /teams [post]
func (h *teamHandler) CreateTeam(c *gin.Context) {
	var input dtos.CreateTeamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	team, err := h.service.CreateTeam(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":   team.ID,
		"name": team.Name,
	})
}

// AddManager godoc
// @Summary Thêm manager vào team
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param data body dtos.AddUserToTeamInput true "Manager data"
// @Success 201
// @Failure 400,500 {object} map[string]string
// @Router /teams/{teamId}/managers [post]
func (h *teamHandler) AddManager(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var input struct {
		UserID uuid.UUID `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddManager(teamID, input.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add manager"})
		return
	}
	c.Status(http.StatusCreated)
}

// RemoveManager godoc
// @Summary Xoá manager khỏi team
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param managerId path string true "Manager ID"
// @Success 204
// @Failure 400,500 {object} map[string]string
// @Router /teams/{teamId}/managers/{managerId} [delete]
func (h *teamHandler) RemoveManager(c *gin.Context) {
	teamID, err1 := uuid.Parse(c.Param("teamId"))
	userID, err2 := uuid.Parse(c.Param("managerId"))
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teamId or managerId"})
		return
	}

	if err := h.service.RemoveManager(teamID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove manager"})
		return
	}
	c.Status(http.StatusOK)
}

// AddMember godoc
// @Summary Thêm thành viên vào team
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param data body dtos.AddUserToTeamInput true "Member data"
// @Success 201
// @Failure 400,500 {object} map[string]string
// @Router /teams/{teamId}/members [post]
func (h *teamHandler) AddMember(c *gin.Context) {
	teamID, err := uuid.Parse(c.Param("teamId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid team ID"})
		return
	}

	var input struct {
		UserID uuid.UUID `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddMember(teamID, input.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add member"})
		return
	}
	c.Status(http.StatusCreated)
}

// RemoveMember godoc
// @Summary Xoá thành viên khỏi team
// @Tags Teams
// @Accept json
// @Produce json
// @Param teamId path string true "Team ID"
// @Param memberId path string true "Member ID"
// @Success 204
// @Failure 400,500 {object} map[string]string
// @Router /teams/{teamId}/members/{memberId} [delete]
func (h *teamHandler) RemoveMember(c *gin.Context) {
	teamID, err1 := uuid.Parse(c.Param("teamId"))
	userID, err2 := uuid.Parse(c.Param("memberId")) // memberId theo RESTful path
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teamId or memberId"})
		return
	}

	if err := h.service.RemoveMember(teamID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove member"})
		return
	}
	c.Status(http.StatusOK)
}
