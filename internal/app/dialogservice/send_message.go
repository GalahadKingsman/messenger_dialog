package dialogservice

import (
	"context"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (s *Service) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	// Валидация входных данных
	if req.DialogId == 0 {
		return nil, status.Error(codes.InvalidArgument, "dialog_id is required")
	}
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "sender_id is required")
	}
	if req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "text cannot be empty")
	}

	// Сохраняем сообщение в БД
	messageID, timestamp, err := s.dialogRepo.SendMessage(req.DialogId, req.UserId, req.Text)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return nil, status.Error(codes.Internal, "failed to send message")
	}

	// Формируем ответ
	return &pb.SendMessageResponse{
		MessageId: messageID,
		Timestamp: timestamp.Unix(), // Конвертируем время в Unix timestamp
	}, nil
}
