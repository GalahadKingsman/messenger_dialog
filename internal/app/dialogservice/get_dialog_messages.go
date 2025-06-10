package dialogservice

import (
	"context"
	"database/sql"
	"errors"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) GetDialogMessages(ctx context.Context, req *pb.GetDialogMessagesRequest) (*pb.GetDialogMessagesResponse, error) {
	// Валидация запроса
	if req.DialogId == 0 {
		return nil, status.Error(codes.InvalidArgument, "dialog_id is required")
	}
	if req.Login == "" {
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	
	// Получаем сообщения из БД
	messages, err := s.dialogRepo.GetDialogMessages(ctx, req.DialogId, req.Login, req.Limit, req.Offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &pb.GetDialogMessagesResponse{Messages: []*pb.Message{}}, nil
		}
		return nil, status.Error(codes.Internal, "failed to get messages: "+err.Error())
	}

	// Конвертируем наши сообщения в protobuf сообщения
	pbMessages := make([]*pb.Message, 0, len(messages))
	for _, msg := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			Id:        msg.ID,
			UserId:    msg.UserID,
			Text:      msg.Text,
			Timestamp: timestamppb.New(msg.CreateDate),
		})
	}

	return &pb.GetDialogMessagesResponse{Messages: pbMessages}, nil
}
