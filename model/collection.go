package model

import (
	"context"
	"encoding/hex"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Type int

// const (
// 	//Temperature ...
// 	Temperature Type = iota
// 	//Humidity ...
// 	Humidity
// 	//SoilMositure ...
// 	SoilMositure
// )

//Collection ...
type Collection struct {
	ID          primitive.ObjectID `bson:"_id"`
	IDCollector string             `bson:"id_collector"`
	Name        string             `bson:"name"`
	Email       string             `bson:"email"`
	Sources     []string           `bson:"sources"`
}

//CreateCollection ...
func CreateCollection(Email, Name string, Sources []string) *Collection {
	src := []byte(Email)

	hexedEmail := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(hexedEmail, src)

	id := strings.Join([]string{strings.ReplaceAll(Name, " ", "-"),
		string(hexedEmail[:10])}, "-")

	return &Collection{primitive.NewObjectID(), id, Name, Email, Sources}
}

func (coll Collection) insert() error {
	_, err := collCollection.InsertOne(context.TODO(), coll)
	return err
}

func (coll Collection) update() error {
	return nil
}

func (coll Collection) delete() error {
	res := collCollection.FindOneAndDelete(context.TODO(), bson.D{{"_id", coll.ID}}, nil)

	return res.Err()
}

func (coll Collection) findOne(filter map[string]interface{}) (interface{}, error) {
	var receivedColl Collection

	err := collCollection.FindOne(context.TODO(), filter).Decode(&receivedColl)

	if err != nil {
		return nil, err
	}
	return receivedColl, nil
}

func (coll Collection) findAll(filter map[string]interface{}) []interface{} {
	var colls []interface{}
	var cur_coll Collection
	cursor, err := collCollection.Find(context.TODO(), filter)

	if err != nil {
		log.Println(err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&cur_coll)

		if err != nil {
			log.Println(err)
		}

		colls = append(colls, cur_coll)
	}

	return colls
}

//
func ToCollection(data interface{}) Collection {
	return data.(Collection)
}

//
func AllToCollection(data []interface{}) []Collection {
	var colls []Collection
	for _, v := range data {
		colls = append(colls, ToCollection(v))
	}
	return colls
}
