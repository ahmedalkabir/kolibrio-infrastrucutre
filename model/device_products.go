package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DeviceProducts
type DeviceProducts struct {
	ID        primitive.ObjectID `bson:"_id"`
	CodeName  string             `bson:"code_name"`
	Name      string             `bson:"name"`
	Sensors   []string           `bson:"sensors"`
	Actuators []string           `bson:"actuators"`
}

//CreateDeviceProduct
func CreateDeviceProduct(codeName, name string, sensors, actuators []string) *DeviceProducts {
	return &DeviceProducts{primitive.NewObjectID(), codeName, name, sensors, actuators}
}

//update
func (div DeviceProducts) update() error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", div.ID}}
	update := bson.D{{"$set", div}}

	res := deviceProductsCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	return res.Err()
}

//delete
func (div DeviceProducts) delete() error {
	res := deviceProductsCollection.FindOneAndDelete(context.TODO(), bson.D{{"code_name", div.CodeName}}, nil)

	return res.Err()
}

func (div DeviceProducts) insert() error {
	_, err := deviceProductsCollection.InsertOne(context.TODO(), div)
	return err
}

func (div DeviceProducts) findOne(filter map[string]interface{}) (interface{}, error) {
	var receivedProduct DeviceProducts

	err := deviceProductsCollection.FindOne(context.TODO(), filter).Decode(&receivedProduct)

	if err != nil {
		return nil, err
	}

	return receivedProduct, nil
}

//
func (div DeviceProducts) findAll(filter map[string]interface{}) []interface{} {
	var products []interface{}
	var cur_product DeviceProducts
	cursor, err := deviceProductsCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&cur_product)

		if err != nil {
			log.Println(err)
		}

		products = append(products, cur_product)
	}

	return products
}

//
func ToDeviceProduct(data interface{}) DeviceProducts {
	return data.(DeviceProducts)
}

//
func AllToDeviceProduct(data []interface{}) []DeviceProducts {
	var products []DeviceProducts
	for _, v := range data {
		products = append(products, ToDeviceProduct(v))
	}

	return products
}
