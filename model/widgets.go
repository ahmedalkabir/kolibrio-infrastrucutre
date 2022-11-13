package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Widget
// 1 = Gauage
// 2 = ToggleSwitch
type Widget struct {
	ID          primitive.ObjectID `bson:"_id"`
	IDDashboard string             `bson:"id_dash"`
	Name        string             `bson:"name"`
	IDDEVICE    string             `bson:"id_device"`
	Source      string             `bson:"source"`
	WidgetType  int32              `bson:"type"`
}

// CreateWidget
func CreateWidget(idDashboard, name, idDevice, source string, widgetType int32) *Widget {
	return &Widget{primitive.NewObjectID(), idDashboard, name, idDevice, source, widgetType}
}

func (widget Widget) insert() error {
	_, err := widgetCollection.InsertOne(context.TODO(), widget)
	return err
}

func (widget Widget) update() error {
	return nil
}

func (widget Widget) delete() error {
	res := widgetCollection.FindOneAndDelete(context.TODO(), bson.D{{"_id", widget.ID}}, nil)

	return res.Err()
}

func (widget Widget) findOne(filter map[string]interface{}) (interface{}, error) {
	var receivedWidget Widget

	err := widgetCollection.FindOne(context.TODO(), filter).Decode(&receivedWidget)

	if err != nil {
		return nil, err
	}
	return receivedWidget, nil
}

func (widget Widget) findAll(filter map[string]interface{}) []interface{} {
	var widgets []interface{}
	var cur_widget Widget
	cursor, err := widgetCollection.Find(context.TODO(), filter)

	if err != nil {
		log.Println(err)
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&cur_widget)

		if err != nil {
			log.Println(err)
		}

		widgets = append(widgets, cur_widget)
	}

	return widgets
}

//
func ToWidget(data interface{}) Widget {
	return data.(Widget)
}

//
func AllToWidget(data []interface{}) []Widget {
	var widgets []Widget
	for _, v := range data {
		widgets = append(widgets, ToWidget(v))
	}
	return widgets
}
