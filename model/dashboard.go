package model

import (
	"context"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Dashboard model
type Dashboard struct {
	ID primitive.ObjectID `bson:"_id"`
	// IDDashboard it will be generated from CreateDashboard
	IDDashboard string `bson:"id_dash"`
	Name        string `bson:"name"`
	Desc        string `bson:"desc"`
	// Email here it will be used as refernces
	Email       string `bson:"email"`
	CreatedTime string `bson:"created_time"`
}

//
func CreateDashboard(name, desc, email string) *Dashboard {
	dt := time.Now()

	src := []byte(dt.String())
	hexedTime := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(hexedTime, src)

	src = []byte(email)
	hexedEmail := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(hexedEmail, src)

	idDash := strings.Join([]string{strings.ReplaceAll(name, " ", "-"),
		string(hexedTime[:10]), string(hexedEmail[:6])}, "-")

	return &Dashboard{primitive.NewObjectID(), idDash, name, desc, email, dt.String()}
}

func (dash Dashboard) update() error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", dash.ID}}
	update := bson.D{{"$set", dash}}

	res := dashboardCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	return res.Err()
}

// Delete
func (dash Dashboard) delete() error {
	res := dashboardCollection.FindOneAndDelete(context.TODO(), bson.D{{"id_dash", dash.IDDashboard}}, nil)

	return res.Err()
}

func (dash Dashboard) insert() error {
	_, err := dashboardCollection.InsertOne(context.TODO(), dash)
	return err
}

func (dash Dashboard) findOne(filter map[string]interface{}) (interface{}, error) {
	var receivedDash Dashboard

	err := dashboardCollection.FindOne(context.TODO(), filter).Decode(&receivedDash)

	if err != nil {
		return nil, err
	}
	return receivedDash, nil
}

func (dash Dashboard) findAll(filter map[string]interface{}) []interface{} {
	var dashboards []interface{}
	var cur_dashboard Dashboard
	cursor, err := dashboardCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&cur_dashboard)
		if err != nil {
			log.Println(err)
		}

		dashboards = append(dashboards, cur_dashboard)
	}

	if err := cursor.Err(); err != nil {
		return nil
	}

	return dashboards
}

//
func ToDashBoard(data interface{}) Dashboard {
	return data.(Dashboard)
}

//AllToDashBoard
func AllToDashBoard(data []interface{}) []Dashboard {
	var dash []Dashboard
	for _, v := range data {
		dash = append(dash, ToDashBoard(v))
	}
	return dash
}
