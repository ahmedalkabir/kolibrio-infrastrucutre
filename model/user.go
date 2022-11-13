package model

import (
	"context"
	"fmt"
	"log"

	"github.com/ahmedalkabir/kolibrio-infrastrucutre/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//User model
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserName string             `bson:"username"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
	Name     string             `bson:"name"`
}

//CreateUser
func CreateUser(userName, Password, Email, Name string) *User {
	hashedPassword := utils.EncryptPassword("test!test", Password)
	return &User{primitive.NewObjectID(), userName, string(hashedPassword), Email, Name}
}

func (user User) update() error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", user.ID}}
	update := bson.D{{"$set", user}}

	res := userCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	return res.Err()
}

func (user User) delete() error {
	opts := options.Delete()
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{"email", user.Email}}, opts)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("deleted %v documents\n", res.DeletedCount)

	return err
}

func (user User) insert() error {
	_, err := userCollection.InsertOne(context.TODO(), user)
	return err
}

func (user User) findOne(filter map[string]interface{}) (interface{}, error) {
	var receivedUser User

	err := userCollection.FindOne(context.TODO(), filter).Decode(&receivedUser)

	if err != nil {
		return nil, err
	}

	return receivedUser, nil
}

func (user User) findAll(filter map[string]interface{}) []interface{} {
	return nil
}

//ToUser
func ToUser(data interface{}) User {
	return data.(User)
}

//PUser ToGet Around return address
func PUser(v User) *User {
	return &v
}
