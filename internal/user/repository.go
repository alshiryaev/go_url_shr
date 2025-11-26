package user

import "go_purple/pkg/db"

type UserRepository struct {
	Database *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Database: db,
	}
}

func (r *UserRepository) Create(user *User) (*User, error) {
	res := r.Database.Create(user)
	if res.Error != nil {
		return nil, res.Error
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var user User
 	res := r.Database.First(&user, "email = ?", email)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}
