package dtos

type CreateTeamInput struct {
	TeamName string `json:"teamName" binding:"required"`

	Managers []struct {
		ManagerID string `json:"managerId" binding:"required"`
	} `json:"managers"`

	Members []struct {
		MemberID string `json:"memberId" binding:"required"`
	} `json:"members"`
}
