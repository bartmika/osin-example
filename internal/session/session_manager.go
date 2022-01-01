package session

import (
	// "fmt"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/bartmika/osin-example/internal/models"
)

type SessionManager struct {
	rdb *redis.Client
}

func New() *SessionManager {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //TODO: Add variable config
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	return &SessionManager{
		rdb: rdb,
	}
}

func (sm *SessionManager) SaveUser(ctx context.Context, sessionUUID string, user *models.User, d time.Duration) error {
	userBin, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = sm.rdb.Set(sessionUUID, userBin, d).Err()
	if err != nil {
		return err
	}
	return nil
}

func (sm *SessionManager) GetUser(ctx context.Context, sessionUUID string) (*models.User, error) {
	userString, err := sm.rdb.Get(sessionUUID).Result()
	if err == redis.Nil {
		// fmt.Println("key2 does not exist")
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		userBin := []byte(userString)
		user := &models.User{}
		err = json.Unmarshal(userBin, user)
		return user, err
	}
}
