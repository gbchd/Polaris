package authentication

import (
	"errors"

	"github.com/guillaumebchd/polaris/pkg/crypto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Email        string             `bson:"email"`
	HashPassword string             `bson:"hash"`
	Role         string             `bson:"role"`
}

func CreateUser(name string, mail string, password string) (*User, error) {
	// hash password
	if isUserInDatabase(mail) {
		return nil, errors.New("There is already a user with this email in the database")
	}

	hash, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Name:         name,
		Email:        mail,
		HashPassword: hash,
		Role:         "user",
	}

	user, err = addUser(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CheckUser(mail string, password string) bool {
	user, err := GetUser(mail)
	if err != nil {
		return false
	}

	return crypto.CheckPasswordHash(password, user.HashPassword)
}
