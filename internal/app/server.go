package app

import pb "messenger_dialog/pkg/messenger_dialog_api"

type Service struct {
	pb.UnimplementedUserServiceServer

	userRepo *user_repo.Repo
}

func New(userRepo *user_repo.Repo) *Service {
	return &Service{
		userRepo: userRepo,
	}
}
