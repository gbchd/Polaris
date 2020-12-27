package authentication

import (
	"errors"

	"github.com/guillaumebchd/polaris/pkg/crypto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Creator      string             `bson:"creator" json:"creator"`
	ClientId     string             `bson:"client_id" json:"client_id"`
	ClientSecret string             `bson:"client_secret" json:"client_secret"`
	Name         string             `bson:"name" json:"name"`
	RedirectUri  string             `bson:"redirect_uri" json:"redirect_uri"`
	Scopes       []string           `bson:"scopes" json:"scopes"`
}

type Scope struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
}

func GenerateClient(creator User, name string, redirect_uri string, scopes []string) (Client, error) {

	client_id, err := generateUniqueClientID()
	if err != nil {
		return Client{}, err
	}

	client_secret, err := crypto.RandomHex(32)
	if err != nil {
		return Client{}, err
	}

	for _, s := range scopes {
		if !IsScopeInDatabase(s) {
			return Client{}, errors.New("Invalid scopes")
		}
	}

	client := Client{
		Creator:      creator.Email,
		ClientId:     client_id,
		ClientSecret: client_secret,
		Name:         name,
		RedirectUri:  redirect_uri,
		Scopes:       scopes,
	}

	// add client to db
	return addClient(client)
}

func generateUniqueClientID() (string, error) {
	client_id, err := crypto.RandomHex(16)
	if err != nil {
		return "", err
	}

	for IsClientInDatabase(client_id) {
		client_id, err = crypto.RandomHex(16)
		if err != nil {
			return "", err
		}
	}

	return client_id, nil
}
