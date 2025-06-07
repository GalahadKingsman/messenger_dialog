package dialogservice

import (
	"github.com/GalahadKingsman/messenger_dialog/internal/repositories/dialog_repo"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
)

type Service struct {
	pb.UnimplementedDialogServiceServer

	dialogRepo *dialog_repo.Repo
}

func New(dialogRepo *dialog_repo.Repo) *Service {
	return &Service{
		dialogRepo: dialogRepo,
	}
}
