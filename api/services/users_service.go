package services

import (
	"know-sync-api/domain/users"
)

func CreateUser(user users.User) (*users.User, error) {
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(userID uint) (*users.User, error) {
	u := &users.User{ID: userID}
	if err := u.Get(); err != nil {
		return nil, err
	}
	return u, nil
}

func GetUserByEmail(email string) (*users.User, error) {
	u := &users.User{Email: email}
	if err := u.GetByEmail(); err != nil {
		return nil, err
	}
	return u, nil
}

func DeleteUser(userID uint) error {
	user := &users.User{ID: userID}
	return user.Delete()
}
