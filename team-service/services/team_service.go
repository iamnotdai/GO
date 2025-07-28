package services

import (
	"fmt"
	"team-service/dtos"
	"team-service/models"
	"team-service/repositories"

	"github.com/google/uuid"
)

type ITeamService interface {
	// Team
	CreateTeam(input dtos.CreateTeamInput) (*models.Team, error)

	// Managers
	AddManager(teamID, userID uuid.UUID) error
	RemoveManager(teamID, userID uuid.UUID) error
	IsManager(teamID, userID uuid.UUID) (bool, error)

	// Members
	AddMember(teamID, userID uuid.UUID) error
	RemoveMember(teamID, userID uuid.UUID) error
}

type teamService struct {
	repo repositories.ITeamRepository
}

func TeamService(repo repositories.ITeamRepository) ITeamService {
	return &teamService{repo}
}

// --- Team ---

func (s *teamService) CreateTeam(input dtos.CreateTeamInput) (*models.Team, error) {
	// Validate and parse UUID
	team := &models.Team{
		ID:   uuid.New(),
		Name: input.TeamName,
	}

	err := s.repo.Create(team)
	if err != nil {
		return nil, err
	}

	// Add Managers
	for _, m := range input.Managers {
		managerID, err := uuid.Parse(m.ManagerID)
		if err != nil {
			return nil, fmt.Errorf("invalid manager ID: %v", err)
		}
		err = s.repo.AddManager(team.ID, managerID)
		if err != nil {
			return nil, fmt.Errorf("failed to add manager: %v", err)
		}
	}

	// Add Members
	for _, m := range input.Members {
		memberID, err := uuid.Parse(m.MemberID)
		if err != nil {
			return nil, fmt.Errorf("invalid member ID: %v", err)
		}
		err = s.repo.AddMember(team.ID, memberID)
		if err != nil {
			return nil, fmt.Errorf("failed to add member: %v", err)
		}
	}

	return team, nil
}

func (s *teamService) GetTeamByID(id uuid.UUID) (*models.Team, error) {
	return s.repo.GetByID(id)
}

// --- Managers ---

func (s *teamService) AddManager(teamID, userID uuid.UUID) error {
	return s.repo.AddManager(teamID, userID)
}

func (s *teamService) RemoveManager(teamID, userID uuid.UUID) error {
	return s.repo.RemoveManager(teamID, userID)
}

func (s *teamService) IsManager(teamID, userID uuid.UUID) (bool, error) {
	return s.repo.IsManager(teamID, userID)
}

// --- Members ---

func (s *teamService) AddMember(teamID, userID uuid.UUID) error {
	return s.repo.AddMember(teamID, userID)
}

func (s *teamService) RemoveMember(teamID, userID uuid.UUID) error {
	return s.repo.RemoveMember(teamID, userID)
}
