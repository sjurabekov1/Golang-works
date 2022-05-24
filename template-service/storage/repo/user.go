package repo

import (
	pb "github.com/sjurabekov1/homeworks/template-service/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	Create(*pb.User) (*pb.User, error)
}
