package dialogservice

import (
	"context"
	"github.com/GalahadKingsman/messenger_dialog/internal/repositories/dialog_repo"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"github.com/redis/go-redis/v9"
	"log"
)

type Service struct {
	pb.UnimplementedDialogServiceServer
	rdb        *redis.Client
	dialogRepo *dialog_repo.Repo
}

func New(dialogRepo *dialog_repo.Repo) *Service {

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", // можно взять из env
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}

	return &Service{
		dialogRepo: dialogRepo,
		rdb:        rdb,
	}
}
