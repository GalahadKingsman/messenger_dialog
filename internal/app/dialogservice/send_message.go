package dialogservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/GalahadKingsman/messenger_dialog/internal/models"
	pb "github.com/GalahadKingsman/messenger_dialog/pkg/messenger_dialog_api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
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
	messageID, createTime, err := s.dialogRepo.SendMessage(req.DialogId, req.UserId, req.Text)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
		return nil, status.Error(codes.Internal, "failed to send message")
	}
	peerID, err := s.dialogRepo.GetPeerID(req.DialogId, req.UserId)
	if err != nil {
		log.Printf("[Dialog] GetPeerID failed: %v", err)
	} else {

		notif := models.Notification{
			From:     fmt.Sprintf("%d", req.UserId),
			Message:  req.Text,
			DialogID: req.DialogId,
		}
		payload, _ := json.Marshal(notif)
		channel := fmt.Sprintf("notifications:%d", peerID)

		go func(ch string, data []byte) {
			ctxPub, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := s.rdb.Publish(ctxPub, ch, data).Err(); err != nil {
				log.Printf("[Dialog] Redis publish error: %v", err)
			} else {
				log.Printf("[Dialog] Published to %s payload=%s", ch, string(data))
			}
		}(channel, payload)
	}

	// Формируем ответ с временем из БД
	return &pb.SendMessageResponse{
		MessageId: messageID,
		Timestamp: timestamppb.New(createTime), // Используем время из БД
	}, nil
}
