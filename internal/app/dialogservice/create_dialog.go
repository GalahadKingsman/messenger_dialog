package dialogservice

import (
	"context"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) CreateDialog(ctx context.Context, req *pb.CreateDialogRequest) (*pb.CreateDialogResponse, error) {
	// Валидация
	if req.UserId == req.PeerId {
		return nil, status.Error(codes.InvalidArgument, "нельзя создать диалог с самим собой")
	}

	// Вызов слоя БД
	dialogID, err := s.userRepo.CreateDialog(ctx, req.UserId, req.PeerId)
	if err != nil {
		return nil, status.Error(codes.Internal, "ошибка создания диалога")
	}

	return &pb.CreateDialogResponse{
		DialogId: dialogID,
		Success:  true,
	}, nil
}
