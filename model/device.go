package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//
type Device struct {
	ID         primitive.ObjectID `bson:"_id"`
	Email      string             `bson:"email"`
	Name       string             `bson:"name"`
	IDDEVICE   string             `bson:"id_device"`
	DeviceType DeviceProducts     `bson:"device_type"`
}

// CreateDevice
func CreateDevice(email, name, id string, typeDevice DeviceProducts) *Device {
	return &Device{primitive.NewObjectID(), email, name, id, typeDevice}
}

func (div Device) update() error {
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", div.ID}}
	update := bson.D{{"$set", div}}

	res := deviceCollection.FindOneAndUpdate(context.TODO(), filter, update, opts)

	return res.Err()
}

// Delete
func (div Device) delete() error {
	res := deviceCollection.FindOneAndDelete(context.TODO(), bson.D{{"_id", div.ID}}, nil)

	return res.Err()
}

func (div Device) insert() error {
	_, err := deviceCollection.InsertOne(context.TODO(), div)
	return err
}

func (div Device) findOne(filter map[string]interface{}) (interface{}, error) {
	var receivedDevice Device

	err := deviceCollection.FindOne(context.TODO(), filter).Decode(&receivedDevice)

	if err != nil {
		return nil, err
	}
	return receivedDevice, nil
}

//
func (div Device) findAll(filter map[string]interface{}) []interface{} {
	var devices []interface{}
	var cur_device Device
	cursor, err := deviceCollection.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&cur_device)

		if err != nil {
			log.Println(err)
		}

		devices = append(devices, cur_device)
	}

	return devices
}

//
func ToDevice(data interface{}) Device {
	return data.(Device)
}

//
func AllToDevice(data []interface{}) []Device {
	var devices []Device
	for _, v := range data {
		devices = append(devices, ToDevice(v))
	}

	return devices
}
