package redisrepo

import pb "template-service3/genproto/user_service"

// UserStorage I ..
type UserStorageI interface {
	Create(*pb.User) (*pb.User, error)
	GetUserById(*pb.GetUserId) (*pb.User, error)
	UpdateUser(*pb.User) (*pb.User, error)
	DeleteUser(*pb.GetUserId) (*pb.User, error)
	GetAllUsers(*pb.GetAllUsersRequest) (*pb.AllUsersResp, error)
}
