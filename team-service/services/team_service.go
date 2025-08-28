package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"team-service/cache"
	"team-service/config"
	"team-service/dtos"
	"team-service/events"
	"team-service/models"
	"team-service/repositories"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
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
	repo          repositories.ITeamRepository
	kafkaProducer *events.KafkaProducer
	redisClient   *redis.Client
}

func TeamService(repo repositories.ITeamRepository, cfg *config.Config) ITeamService {
	kafkaProducer, err := events.NewKafkaProducer(cfg) // Bây giờ trả về *KafkaProducer
	if err != nil {
		panic(fmt.Sprintf("failed to create kafka producer: %v", err))
	}
	redisClient := cache.NewRedisClient(cfg)
	return &teamService{
		repo:          repo,
		kafkaProducer: kafkaProducer,
		redisClient:   redisClient,
	}
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

	// Gửi sự kiện "team.created" đến Kafka
	teamEvent := map[string]interface{}{
		"team_id":   team.ID.String(),
		"team_name": team.Name,
	}
	teamEventBytes, err := json.Marshal(teamEvent)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal team event: %v", err)
	}
	err = s.kafkaProducer.PublishMessage(events.TopicTeamCreated, team.ID.String(), teamEventBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to publish team created event: %v", err)
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
	ctx := context.Background()
	key := fmt.Sprintf("team:%s", id.String())

	// Thử lấy dữ liệu từ Redis
	data, err := s.redisClient.Get(ctx, key).Result()
	if err == nil {
		// Dữ liệu tồn tại trong Redis, parse JSON
		var teamEvent struct {
			TeamID   string `json:"team_id"`
			TeamName string `json:"team_name"`
		}
		if err := json.Unmarshal([]byte(data), &teamEvent); err != nil {
			log.Printf("failed to unmarshal team data from redis: %v", err)
			// Không trả lỗi, tiếp tục lấy từ DB
		} else {
			// Tạo đối tượng models.Team từ dữ liệu Redis
			team := &models.Team{
				ID:   id,
				Name: teamEvent.TeamName,
			}
			return team, nil
		}
	} else if err != redis.Nil {
		// Lỗi Redis không phải là "key không tồn tại", log lỗi nhưng vẫn tiếp tục
		log.Printf("failed to get team from redis: %v", err)
	}

	// Nếu Redis không có dữ liệu hoặc gặp lỗi, lấy từ DB
	team, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get team from db: %v", err)
	}

	// Lưu dữ liệu vào Redis để sử dụng lần sau
	teamEvent := struct {
		TeamID   string `json:"team_id"`
		TeamName string `json:"team_name"`
	}{
		TeamID:   team.ID.String(),
		TeamName: team.Name,
	}
	jsonBytes, err := json.Marshal(teamEvent)
	if err != nil {
		log.Printf("failed to marshal team data for redis: %v", err)
		// Không trả lỗi, vì dữ liệu đã lấy được từ DB
	} else {
		data = string(jsonBytes)
		err = s.redisClient.Set(ctx, key, data, 0).Err()
		if err != nil {
			log.Printf("failed to save team to redis: %v", err)
		}
	}

	return team, nil
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

// Đóng producer khi service dừng
func (s *teamService) Close() error {
	if err := s.kafkaProducer.Close(); err != nil {
		return err
	}
	return s.redisClient.Close()
}
