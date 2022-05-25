package services

import (
	"github.com/ko1ke/know-sync-api/datasources/postgres_db"
	"github.com/ko1ke/know-sync-api/domain/users"
)

func CreateUser(user users.User) (*users.User, error) {
	if err := user.Save(postgres_db.Client); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(userID uint) (*users.User, error) {
	u := &users.User{ID: userID}
	if err := u.Get(postgres_db.Client); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*users.User, error) {
	u := &users.User{Email: email}
	if err := u.GetByEmail(postgres_db.Client); err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteUser(userID uint) error {
	user := &users.User{ID: userID}
	return user.Delete(postgres_db.Client)
}
