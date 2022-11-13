package model

import (
	"context"
	"fmt"
	"log"

	"github.com/ahmedalkabir/kolibrio-infrastrucutre/config"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// a global varibles to hold reference of mongodb connection
// for all models
var clientOptions *options.ClientOptions
var client *mongo.Client
var database *mongo.Database
var userCollection *mongo.Collection
var dashboardCollection *mongo.Collection
var deviceCollection *mongo.Collection
var deviceProductsCollection *mongo.Collection
var widgetCollection *mongo.Collection

// Collection Collection it's all about
// creating document of collection that will be
// to store list of tables
// and tables are all about getting data from influxdb
var collCollection *mongo.Collection

var err error

// prepare our database before start our
// server
func init() {
	// prepare a connection to mongodb
	var uri string
	// let's check if we have user and password
	// or not
	if config.MongoDBConf.User == "" {
		uri = fmt.Sprintf("mongodb://%s:%s", config.MongoDBConf.Address, config.MongoDBConf.Port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", config.MongoDBConf.User, config.MongoDBConf.Password,
			config.MongoDBConf.Address, config.MongoDBConf.Port)
	}

	clientOptions = options.Client().ApplyURI(uri)
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// let's link our link to database
	database = client.Database("kolibrio-aplha")

	// now setup our collection
	userCollection = database.Collection("User")
	dashboardCollection = database.Collection("Dashboard")
	deviceCollection = database.Collection("Device")
	deviceProductsCollection = database.Collection("Device_Products")
	widgetCollection = database.Collection("Widget")
	collCollection = database.Collection("Collection")

	// to avoid of repeating some fields
	// TODO: look for better way of creating index
	// user
	_, err := userCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"email", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Println(err.Error())
	}

	//dashboard
	_, err = dashboardCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"id_dash", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Println(err.Error())
	}

	//device
	// this is really stupid idea, look for
	// better one bitch
	// _, err = deviceCollection.Indexes().CreateOne(
	// 	context.TODO(),
	// 	mongo.IndexModel{
	// 		Keys:    bsonx.Doc{{"id_device", bsonx.Int32(1)}},
	// 		Options: options.Index().SetUnique(true),
	// 	},
	// )

	// if err != nil {
	// 	log.Println(err.Error())
	// }

	//device_products
	_, err = deviceProductsCollection.Indexes().CreateOne(
		context.TODO(),
		mongo.IndexModel{
			Keys:    bsonx.Doc{{"code_name", bsonx.Int32(1)}},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		log.Println(err.Error())
	}
}

//MongoOperation
type MongoOperation interface {
	insert() error
	update() error
	delete() error
	findAll(filter map[string]interface{}) []interface{}
	findOne(filter map[string]interface{}) (interface{}, error)
}

//GetFromDB
// func GetFromDB(value string, opr MongoOperation) (interface{}, error) {
// 	return opr.get(value)
// }

//InsertToDB
func InsertToDB(opr MongoOperation) error {
	return opr.insert()
}

//DeleteFromDB
func DeleteFromDB(opr MongoOperation) error {
	return opr.delete()
}

//UpdateValues
func UpdateValues(opr MongoOperation) error {
	return opr.update()
}

//GetDocument
func GetDocument(filter map[string]interface{}, opr MongoOperation) (interface{}, error) {
	return opr.findOne(filter)
}

//GetAllDocuments
func GetAllDocuments(filter map[string]interface{}, opr MongoOperation) []interface{} {
	return opr.findAll(filter)
}

//Filter well it's just a helper function to make things simple
// I used it with GetDocument and GetAllDocuments as both require filter
func Filter(key, value string) map[string]interface{} {
	return map[string]interface{}{key: value}
}
