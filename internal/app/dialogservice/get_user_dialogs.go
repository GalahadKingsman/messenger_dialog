package dialogservice

import (
	"context"
	"github.com/GalahadKingsman/messenger_dialog/internal/models"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *Service) GetUserDialogs(ctx context.Context, req *pb.GetUserDialogsRequest) (*pb.GetUserDialogsResponse, error) {
	// Устанавливаем значения по умолчанию
	limit := int32(50) // значение по умолчанию
	if req.Limit != nil && *req.Limit > 0 {
		limit = *req.Limit
	}

	offset := int32(0) // значение по умолчанию
	if req.Offset != nil {
		offset = *req.Offset
	}

	// Получаем данные из репозитория
	dialogs, err := s.dialogRepo.GetUserDialogs(req.UserId, limit, offset)
	if err != nil {
		log.Printf("Failed to get user dialogs: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get dialogs: %v", err)
	}

	return &pb.GetUserDialogsResponse{Dialogs: convertToPbDialogs(dialogs)}, nil
}

func convertToPbDialogs(modelsDialogs []*models.DialogInfo) []*pb.DialogInfo {
	pbDialogs := make([]*pb.DialogInfo, len(modelsDialogs))
	for i, d := range modelsDialogs {
		pbDialogs[i] = &pb.DialogInfo{
			DialogId:    d.ID,
			PeerId:      d.PeerID,
			PeerLogin:   d.PeerLogin,
			LastMessage: d.LastMessage,
		}
	}
	return pbDialogs
}
