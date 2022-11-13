package config

import (
	"log"
	"os"
)

//Server this srtuct it will be use
// to get enviroment variables values to be
// shared over the program
type Server struct {
	Address string
	Port    string
}

//MongoDB
type MongoDB struct {
	Address  string
	Port     string
	User     string
	Password string
}

//MQTT
type MQTT struct {
	Address string
	Port    string
	ID      string
}

//InfluxDB
type InfluxDB struct {
	Address string
	Port    string
}

var (
	ServerConf  Server
	MongoDBConf MongoDB
	MqttConf    MQTT
	InfluxConf  InfluxDB
)

func init() {
	log.Println("CONFIG PACKAGE Called")

	// Read Enviroment Variables of Server
	ServerConf.Address = os.Getenv("ADDR")
	ServerConf.Port = os.Getenv("PORT")

	// Read Enviroment Variables of MongoDB
	MongoDBConf.Address = os.Getenv("MONGO_ADDR")
	MongoDBConf.Port = os.Getenv("MONGO_PORT")
	MongoDBConf.User = os.Getenv("MONGO_USER")
	MongoDBConf.Password = os.Getenv("MONGO_PASS")

	// Read Enviroment Variables of MongoDB
	MqttConf.Address = os.Getenv("MQTT_ADDR")
	MqttConf.Port = os.Getenv("MQTT_PORT")
	MqttConf.ID = os.Getenv("MQTT_ID")

	InfluxConf.Address = os.Getenv("INFLUX_ADDR")
	InfluxConf.Port = os.Getenv("INFLUX_PORT")

	log.Println(MongoDBConf)
}
