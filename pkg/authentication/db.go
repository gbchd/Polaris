package authentication

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database
var ctx = context.TODO()

var userCollection *mongo.Collection

var clientCollection *mongo.Collection
var scopeCollection *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := os.Getenv("DB_ADDR")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")

	url := "mongodb://" + username + ":" + password + "@" + addr + ":" + port

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database(database)

	userCollection = db.Collection("user")

	clientCollection = db.Collection("client")
	scopeCollection = db.Collection("scope")

}

/*
	USER
*/
func addUser(user User) (User, error) {
	res, err := userCollection.InsertOne(ctx, user)

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.Id = oid
		return user, err
	} else {
		return User{}, errors.New("Did not find id of inserted user")
	}
}

func GetUser(mail string) (User, error) {
	var user User
	err := userCollection.FindOne(ctx, bson.M{"email": mail}).Decode(&user)
	return user, err
}

func isUserInDatabase(mail string) bool {
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": mail})
	if err != nil {
		panic(err)
	}
	return count >= 1
}

func SetUserAdmin(mail string) error {
	update := bson.M{"$set": bson.M{"role": "admin"}}
	res := userCollection.FindOneAndUpdate(ctx, bson.M{"email": mail}, update)
	return res.Err()
}

func deleteUser(mail string) error {
	res, err := userCollection.DeleteMany(ctx, bson.M{"email": mail})
	if res.DeletedCount < 1 {
		return errors.New("No user with mail : " + mail)
	}
	return err
}

/*
	SCOPES
*/

func AddScope(scope Scope) (Scope, error) {
	res, err := scopeCollection.InsertOne(ctx, scope)

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		scope.Id = oid
		return scope, err
	} else {
		return Scope{}, errors.New("Did not find id of inserted scope")
	}
}

func GetScope(name string) (Scope, error) {
	var scope Scope
	err := scopeCollection.FindOne(ctx, bson.M{"name": name}).Decode(&scope)
	return scope, err
}

func IsScopeInDatabase(name string) bool {
	count, err := scopeCollection.CountDocuments(ctx, bson.M{"name": name})
	if err != nil {
		panic(err)
	}
	return count >= 1
}

func DeleteScope(name string) error {
	res, err := scopeCollection.DeleteMany(ctx, bson.M{"name": name})
	if res.DeletedCount < 1 {
		return errors.New("No scope with name : " + name)
	}
	return err
}

/*
	Client
*/

func addClient(client Client) (Client, error) {
	res, err := clientCollection.InsertOne(ctx, client)

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		client.Id = oid
		return client, err
	} else {
		return Client{}, errors.New("Did not find id of inserted client")
	}
}

func GetClient(client_id string) (Client, error) {
	var client Client
	err := clientCollection.FindOne(ctx, bson.M{"client_id": client_id}).Decode(&client)
	return client, err
}

func IsClientInDatabase(client_id string) bool {
	count, err := clientCollection.CountDocuments(ctx, bson.M{"client_id": client_id})
	if err != nil {
		panic(err)
	}
	return count >= 1
}

func DeleteClient(client_id string) error {
	res, err := scopeCollection.DeleteMany(ctx, bson.M{"client_id": client_id})
	if res.DeletedCount < 1 {
		return errors.New("No client with client_id : " + client_id)
	}
	return err
}
