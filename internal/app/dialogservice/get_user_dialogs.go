package dialogservice

import (
	"context"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Service) GetUserDialogs(ctx context.Context, req *pb.GetUserDialogsRequest) (*pb.GetUserDialogsResponse, error) {
	// Установка значений по умолчанию
	limit := 20
	if req.Limit != nil {
		limit = int(*req.Limit)
	}

	offset := 0
	if req.Offset != nil {
		offset = int(*req.Offset)
	}

	// Вызов слоя БД
	dialogs, err := s.userRepo.GetUserDialogs(ctx, req.UserId, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "ошибка получения диалогов")
	}

	// Конвертация в protobuf-формат
	pbDialogs := make([]*pb.DialogInfo, 0, len(dialogs))
	for _, d := range dialogs {
		pbDialogs = append(pbDialogs, &pb.DialogInfo{
			DialogId:     d.ID,
			PeerId:       d.PeerID,
			PeerLogin:    d.PeerLogin,
			LastMessage:  d.LastMessage,
			LastActivity: d.LastActivity.Unix(),
		})
	}

	return &pb.GetUserDialogsResponse{Dialogs: pbDialogs}, nil
}
