package repositories

import (
	"team-service/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ITeamRepository interface {
	// Team
	GetByID(id uuid.UUID) (*models.Team, error)
	Create(team *models.Team) error

	// Managers
	AddManager(teamID, userID uuid.UUID) error
	RemoveManager(teamID, userID uuid.UUID) error
	IsManager(teamID, userID uuid.UUID) (bool, error)

	// Members
	AddMember(teamID, userID uuid.UUID) error
	RemoveMember(teamID, userID uuid.UUID) error
}

type teamRepository struct {
	db *gorm.DB
}

func TeamRepository(db *gorm.DB) ITeamRepository {
	return &teamRepository{db}
}

// --- Team ---

func (r *teamRepository) GetByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	if err := r.db.Preload("Members").Preload("Managers").
		First(&team, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &team, nil
}

func (r *teamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

// --- Managers ---

func (r *teamRepository) AddManager(teamID, userID uuid.UUID) error {
	return r.db.Create(&models.TeamManager{
		TeamID: teamID,
		UserID: userID,
	}).Error
}

func (r *teamRepository) RemoveManager(teamID, userID uuid.UUID) error {
	return r.db.Delete(&models.TeamManager{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}

func (r *teamRepository) IsManager(teamID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&models.TeamManager{}).
		Where("team_id = ? AND user_id = ?", teamID, userID).
		Count(&count).Error
	return count > 0, err
}

// --- Members ---

func (r *teamRepository) AddMember(teamID, userID uuid.UUID) error {
	return r.db.Create(&models.TeamMember{
		TeamID: teamID,
		UserID: userID,
	}).Error
}

func (r *teamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	return r.db.Delete(&models.TeamMember{}, "team_id = ? AND user_id = ?", teamID, userID).Error
}
