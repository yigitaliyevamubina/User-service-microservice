package postgres

import (
	"fmt"
	pb "template-service3/genproto/user_service"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

// NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *pb.User) (*pb.User, error) {
	query := `INSERT INTO users(id, first_name, last_name, age, gender, email, password) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, first_name, last_name, age, gender, email, password`

	user.Id = uuid.New().String()

	var respUser pb.User

	rowUser := r.db.QueryRow(query, user.Id,
		user.FirstName,
		user.LastName,
		user.Age,
		user.Gender,
		user.Email,
		user.Password)
	if err := rowUser.Scan(&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Age,
		&respUser.Gender,
		&respUser.Email,
		&respUser.Password); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (r *userRepo) GetUserById(userId *pb.GetUserId) (*pb.UserWithPostsAndComments, error) {
	query := `SELECT id, first_name, last_name, age, gender, email, password FROM users WHERE id = $1 AND deleted_at IS NULL`

	var respUser pb.UserWithPostsAndComments

	rowUser := r.db.QueryRow(query, userId.UserId)
	if err := rowUser.Scan(&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Age,
		&respUser.Gender,
		&respUser.Email,
		&respUser.Password); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (r *userRepo) DeleteUser(userId *pb.GetUserId) (*pb.User, error) {
	query := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP 
             	WHERE id = $1 
             	RETURNING id, first_name, last_name, age, gender, email, password`
	rowUser := r.db.QueryRow(query, userId.UserId)

	var deletedUser pb.User
	if err := rowUser.Scan(&deletedUser.Id,
		&deletedUser.FirstName,
		&deletedUser.LastName,
		&deletedUser.Age,
		&deletedUser.Gender,
		&deletedUser.Email,
		&deletedUser.Password); err != nil {
		return nil, err
	}

	return &deletedUser, nil
}

func (r *userRepo) UpdateUser(user *pb.User) (*pb.User, error) {
	query := `UPDATE users SET first_name = $1, last_name = $2, age = $3, email = $4, password = $5 WHERE id = $6 AND deleted_at IS NULL RETURNING id, first_name, last_name, age, gender, email, password`

	rowUser := r.db.QueryRow(query, user.FirstName, user.LastName, user.Age, user.Id)
	if err := rowUser.Scan(&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.Gender,
		&user.Email,
		&user.Password); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) GetAllUsers(req *pb.GetAllUsersRequest) (*pb.AllUsers, error) {
	query := `SELECT id, first_name, last_name, age, gender, email, password FROM users WHERE deleted_at IS NULL`

	var users pb.AllUsers

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user pb.UserWithPostsAndComments
		if err := rows.Scan(&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Age,
			&user.Gender,
			&user.Email,
			&user.Password); err != nil {
			return nil, err
		}

		users.Users = append(users.Users, &user)
	}

	return &users, nil
}

// CheckField
func (r *userRepo) CheckField(req *pb.Request) (*pb.Response, error) {
	query := fmt.Sprintf(`SELECT count(1) FROM users WHERE %s = $1 AND deleted_at IS NULL`, req.Field)

	var isExists int

	row := r.db.QueryRow(query, req.Data)

	if err := row.Scan(&isExists); err != nil {
		return nil, err
	}

	if isExists == 1 {
		return &pb.Response{Resp: true}, nil
	}

	return &pb.Response{Resp: false}, nil
}
