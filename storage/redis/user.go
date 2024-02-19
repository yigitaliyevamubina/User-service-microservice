package redis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	pb "template-service3/genproto/user_service"
)

type redisRepo struct {
	rConn *redis.Pool
}

func NewUserRepo(rd *redis.Pool) *redisRepo {
	return &redisRepo{rConn: rd}
}

func (r *redisRepo) Create(user *pb.User) (*pb.User, error) {
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	conn := r.rConn.Get()
	defer conn.Close()

	_, err = conn.Do("SET", user.Id, bytes)
	if err != nil {
		return nil, err
	}

	resp, err := redis.Bytes(conn.Do("GET", user.Id))
	if err != nil {
		return nil, err
	}

	respUser := pb.User{}
	if err := json.Unmarshal(resp, &respUser); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (r *redisRepo) GetUserById(userId *pb.GetUserId) (*pb.User, error) {
	conn := r.rConn.Get()
	defer conn.Close()

	resp, err := redis.Bytes(conn.Do("GET", userId.UserId))
	if err != nil {
		return nil, err
	}

	respUser := pb.User{}
	if err := json.Unmarshal(resp, &respUser); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (r *redisRepo) DeleteUser(userId *pb.GetUserId) (*pb.User, error) {
	conn := r.rConn.Get()
	defer conn.Close()

	resp, err := redis.Bytes(conn.Do("GET", userId.UserId))
	if err != nil {
		return nil, err
	}

	respUser := pb.User{}
	if err := json.Unmarshal(resp, &respUser); err != nil {
		return nil, err
	}

	_, err = conn.Do("DEL", userId.UserId)
	if err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (r *redisRepo) UpdateUser(user *pb.User) (*pb.User, error) {
	conn := r.rConn.Get()
	defer conn.Close()

	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	_, err = conn.Do("SET", user.Id, bytes)
	if err != nil {
		return nil, err
	}

	resp, err := redis.Bytes(conn.Do("GET", user.Id))
	if err != nil {
		return nil, err
	}

	respUser := pb.User{}
	if err := json.Unmarshal(resp, &respUser); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (r *redisRepo) GetAllUsers(req *pb.GetAllUsersRequest) (*pb.AllUsersResp, error) {
	conn := r.rConn.Get()
	defer conn.Close()
	keys, err := redis.Strings(conn.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}

	users := []*pb.User{}
	for _, key := range keys {
		value, err := redis.Bytes(conn.Do("GET", key))
		if err != nil {
			return nil, err
		}

		user := pb.User{}
		if err := json.Unmarshal(value, &user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	allUsers := pb.AllUsersResp{
		Users: users,
	}

	return &allUsers, nil
}
