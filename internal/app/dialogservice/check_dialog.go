package dialogservice

import (
	"context"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"log"
)

func (s *Service) CheckDialog(ctx context.Context, req *pb.CheckDialogRequest) (*pb.CheckDialogResponse, error) {
	// Вызываем существующую функцию поиска диалога
	dialogID, _, err := s.dialogRepo.CheckDialog(req.UserId, req.PeerId)
	if err != nil {
		log.Printf("Failed to check dialog: %v", err)
		return nil, err
	}

	// Формируем ответ
	exists := dialogID > 0
	response := &pb.CheckDialogResponse{
		Exists:   exists,
		DialogId: int32(dialogID),
	}

	return response, nil
}
