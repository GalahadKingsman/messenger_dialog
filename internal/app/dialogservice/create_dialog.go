package dialogservice

import (
	"context"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
)

func (s *Service) CreateDialog(ctx context.Context, req *pb.CreateDialogRequest) (*pb.CreateDialogResponse, error) {
	// Проверяем, существует ли диалог между пользователями
	dialogID, dialogName, err := s.dialogRepo.CheckDialog(req.UserId, req.PeerId)
	if err != nil {
		return nil, err
	}

	// Если диалог существует - возвращаем его
	if dialogID != 0 {
		return &pb.CreateDialogResponse{
			Success:    true,
			DialogId:   int32(dialogID),
			DialogName: dialogName,
		}, nil
	}

	// Если диалога нет - создаем новый
	NewDialogID, err := s.dialogRepo.CreateDialog(req.UserId, req.PeerId, req.DialogName)
	if err != nil {
		return nil, err
	}

	return &pb.CreateDialogResponse{
		Success:    true,
		DialogId:   int32(NewDialogID),
		DialogName: req.DialogName,
	}, nil
}
