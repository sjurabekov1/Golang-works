package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/sjurabekov1/homeworks/template-service/genproto"
	l "github.com/sjurabekov1/homeworks/template-service/pkg/logger"
	"github.com/sjurabekov1/homeworks/template-service/storage"
)

//UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
}

//NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	return nil, nil
}
