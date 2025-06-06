package dialogservice

import (
	"github.com/GalahadKingsman/messenger_dialog/internal/repositories/dialog_repo"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
)

type Service struct {
	pb.UnimplementedDialogServiceServer

	userRepo *dialog_repo.Repo
}

func New(userRepo *dialog_repo.Repo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
