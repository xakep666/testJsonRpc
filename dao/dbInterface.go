package dao

import "testJsonRpc/model"

var dbImpl UsersDb

type UsersDb interface {
	// Register attempts to register user with given login.
	// Returns non-nil error value if registration failed, nil otherwise
	Register(login string) error
	// GetByLogin searches user with given login.
	// Returns user info and nil error if search was successful
	GetByLogin(login string) (model.User, error)
	// Edit edits user data with given login.
	// Returns non-nil value if edition was failed, nil otherwise.
	Edit(login string, newData model.User) error
	// Close closes connection to database
	Close()
}

func GetDb() UsersDb {
	return dbImpl
}