package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

//global variables to share between different services
var middlewareToMqtt chan MqttChannel
var mqttToMiddleWare chan MqttChannel
var middleawareToWebSocket chan MqttChannel
var webServerToMqtt chan MqttChannel

var listener *mqttStreamer
var influxStream *influxStreamer

func main() {

	listener = newStreamer("mqtt-1")
	influxStream = newInfluxStreamer("influx")

	measure := &influxObserver{"admin@lambda.ly", "/meshly_1/temperature"}
	influxStream.register(measure)

	middlewareToMqtt = make(chan MqttChannel)
	mqttToMiddleWare = make(chan MqttChannel)
	middleawareToWebSocket = make(chan MqttChannel)
	webServerToMqtt = make(chan MqttChannel)

	go startMqttService(listener, influxStream, middlewareToMqtt, mqttToMiddleWare)

	go startMiddleWareService(mqttToMiddleWare, middlewareToMqtt,
		middleawareToWebSocket, webServerToMqtt)

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	// create a new router to handle our requests
	router := httprouter.New()

	// user := model.CreateUser("admin", "test2020", "admin@lambda.ly", "Ahmed Alkabir")
	// model.InsertToDB(user)

	// coll := model.CreateCollection("admin@lambda.ly", "Meshly Collector", []string{"/meshly_1/temperature"})
	// model.InsertToDB(coll)

	// model.CreateDashboard("My Home 2", "To Monitor and Control My Home", "korabekaali@outlook.com").InsertDashToDB()
	// fmt.Println()
	// model.InsertToDB(
	// 	model.CreateDashboard("My Home 4", "To Monitor and Control My Home", "test@test.com"))
	// widget := model.CreateWidget("My-Home-2-323032", "Temperature", "meshly_324", "temperature", 1)
	// model.InsertToDB(widget)

	// serve static files
	// TODO: you should take path of static files from config file
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	// list of router
	router.GET("/", auth(index))

	router.GET("/login", loginGET)
	router.POST("/login", loginPOST)
	router.GET("/logout", auth(logoutGET))

	router.GET("/dashboards", auth(dashboards))
	router.GET("/dashboards/:dash", auth(getDashboard))
	router.POST("/dashboards/", auth(addDashboard))

	router.GET("/devices", auth(devices))
	router.POST("/devices/", auth(addDevice))
	router.DELETE("/devices/:device", auth(deleteDevice))

	router.POST("/dashboards/widgets/", auth(addWidget))
	router.DELETE("/dashboards/widgets/:dash/:device/:name", auth(deleteWidget))

	router.GET("/sources/:device", auth(getSources))
	router.GET("/sources/:device/:type", auth(getSources))

	router.GET("/collection", auth(collection))
	router.POST("/collection/", auth(addCollection))
	router.DELETE("/collection/:collection", auth(deleteCollection))

	router.GET("/collection/:collection", auth(getCollection))

	// websocket

	router.GET("/sensor/:dash", sensor)
	router.GET("/actuator/:dash", actuator)

	// this mandatory part to make sure
	// our app it will work on heroku platfrom
	var port string
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		port = "8080"
	} else {
		port = portEnv
	}

	// starting up the server
	server := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    time.Duration(20 * int64(time.Second)),
		WriteTimeout:   time.Duration(20 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(server.ListenAndServe())
	// <-c

}
