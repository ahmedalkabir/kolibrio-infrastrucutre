package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmedalkabir/kolibrio-infrastrucutre/model"
	"github.com/ahmedalkabir/kolibrio-infrastrucutre/utils"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

//index   it's going to give back a main page
func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	// get number of dashboards to show in index front page
	dashs := model.AllToDashBoard(model.GetAllDocuments(model.Filter("email", user.Email), model.Dashboard{}))
	numberOfDashes := len(dashs)

	data := struct {
		User    *model.User
		NumDash int
	}{
		user,
		numberOfDashes,
	}

	generatePage(w, "index", data)
}

//auth to check if user is logged in or not
func auth(handler httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if userValue := utils.GetCookie("kolibrio-user", r); userValue != nil {
			handler(w, r, ps)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}

//loginGET
func loginGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	// reason for created custom function of login
	// because of different layout

	// in case we received message from loginPOST
	if msg := w.Header().Get("message"); msg != "" {
		generateLoginPage(w, msg)
	}

	generateLoginPage(w, nil)
}

//loginPOST
func loginPOST(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err.Error())
	}

	// get an email from database
	// user, err := model.GetUserFromDB(r.PostForm.Get("email"))
	getUser, err := model.GetDocument(
		model.Filter("email", r.PostForm.Get("email")),
		model.User{})

	// in case there isn't any user related to that email
	// we should notify it with flash message
	if err != nil {
		log.Println(err.Error())

		w.Header().Add("message", "Make sure you typed your email and password correctly")
		loginGET(w, r, nil)
	} else {
		user := model.ToUser(getUser)
		// now let's compare the hashed password if it is the same
		if user.Password != string(utils.EncryptPassword("test!test", r.PostForm.Get("password"))) {
			w.Header().Add("message", "Make sure you typed your email and password correctly")
			loginGET(w, r, nil)
		}
		// create a new session with updated data
		if cookie := utils.SetCookie("kolibrio-user", "email", user.Email); cookie != nil {
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/", 302)
		}
	}
}

//logoutGET
func logoutGET(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie, err := r.Cookie("kolibrio-user")

	if err != nil {
		log.Println(err)
	}

	cookie.MaxAge = -1

	http.SetCookie(w, cookie)

	// redirect to login page
	http.Redirect(w, r, "/", 302)
}

//Dashboard
func dashboards(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	dash := model.AllToDashBoard(
		model.GetAllDocuments(model.Filter("email", user.Email), model.Dashboard{}))

	data := struct {
		User *model.User
		Dash []model.Dashboard
	}{
		user,
		dash,
	}

	generatePage(w, "list_of_dashboards", data)
}

//
func getDashboard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	// get data from database based on email
	dash, err := model.GetDocument(model.Filter("id_dash", ps.ByName("dash")), model.Dashboard{})
	userDevices := model.AllToDevice(model.GetAllDocuments(model.Filter("email", user.Email), model.Device{}))
	widgets := model.AllToWidget(model.GetAllDocuments(model.Filter("id_dash", ps.ByName("dash")), model.Widget{}))

	if err != nil {
		log.Println(err.Error())
	}
	// TODO: optimize this lines
	data := struct {
		Host    string
		User    *model.User
		Dash    model.Dashboard
		Devices []model.Device
		Widgets []model.Widget
	}{
		r.Host,
		user,
		model.ToDashBoard(dash),
		userDevices,
		widgets,
	}

	generatePage(w, "dashboard", data)
}

//addDashboard
func addDashboard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := getUserDirectly(r)

	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err)
	}

	// parse body
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var jsonDash map[string]json.RawMessage
	err = json.Unmarshal(body, &jsonDash)

	dashboard := model.CreateDashboard(strings.Trim(string(jsonDash["name"]), "\""), strings.Trim(string(jsonDash["desc"]), "\""), user.Email)

	err = model.InsertToDB(dashboard)

	if err != nil {
		log.Println(err)
		sendJSONMessage(w, "error", "make sure that you didn't duplicate the dashboard name")
	} else {
		// return true
		sendJSONMessage(w, "true", "Dashboard has been added successfully")
	}
}

//getSources
// it's used to return sources of data based on device and gauge
// which means of user want device_1 will show up a sources (sensors or actuator) related to that device
// it's used inside js code to update interface based on user selection
func getSources(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := getUserDirectly(r)

	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err)
	}

	idDevice := ps.ByName("device")

	typeWidget := ps.ByName("type")
	// in somecase the type will not be used
	// so we'll set the default value as
	// 1
	if typeWidget == "" {
		typeWidget = "1"
	}

	device, err := model.GetDocument(model.Filter("id_device", idDevice), model.Device{})

	if err != nil {
		log.Println(err)
	}

	if device != nil {
		switch typeWidget {
		case "1":
			generateJSON(w, model.ToDevice(device).DeviceType.Sensors)
		case "2":
			generateJSON(w, model.ToDevice(device).DeviceType.Actuators)
		}
	}

}

//Devices
func devices(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	devices := model.AllToDeviceProduct(model.GetAllDocuments(nil, model.DeviceProducts{}))
	userDevices := model.AllToDevice(model.GetAllDocuments(model.Filter("email", user.Email), model.Device{}))

	data := struct {
		User     *model.User
		Products []model.DeviceProducts
		Devices  []model.Device
	}{
		user,
		devices,
		userDevices,
	}

	generatePage(w, "list_of_devices", data)
}

//addDevice
func addDevice(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	// parse body
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var jsonDevice map[string]json.RawMessage
	err = json.Unmarshal(body, &jsonDevice)

	// i need to trim ""
	name := strings.Trim(string(jsonDevice["name"]), "\"")
	id := strings.Trim(string(jsonDevice["id"]), "\"")
	productName := strings.Trim(string(jsonDevice["product"]), "\"")
	product, err := model.GetDocument(model.Filter("code_name", productName), model.DeviceProducts{})

	device := model.CreateDevice(user.Email, name, id, model.ToDeviceProduct(product))

	err = model.InsertToDB(device)

	if err != nil {
		log.Println(err)
		sendJSONMessage(w, "error", "make sure that you didn't duplicate the device")
	} else {
		sendJSONMessage(w, "true", "The Device Has been added Successfully")
	}

}

//deleteDevice
func deleteDevice(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// get device id from the url and delete it from the db
	deviceID := ps.ByName("device")
	device, err := model.GetDocument(model.Filter("id_device", deviceID), model.Device{})

	if err != nil {
		log.Println(err)
		sendJSONMessage(w, "error", err.Error())
	}

	model.DeleteFromDB(model.ToDevice(device))

	if err != nil {
		log.Println(err)
		sendJSONMessage(w, "error", err.Error())
	} else {
		sendJSONMessage(w, "true", "The Device Has been added Successfully")
	}
}

//add widget
func addWidget(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	// parse body
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)

	var jsonWidget map[string]json.RawMessage
	err = json.Unmarshal(body, &jsonWidget)

	// TODO: optimize this lines
	name := strings.Trim(string(jsonWidget["name"]), "\"")
	typeWidget, err := strconv.Atoi(strings.Trim(string(jsonWidget["type"]), "\""))

	device := strings.Trim(string(jsonWidget["device"]), "\"")
	source := strings.Trim(string(jsonWidget["source"]), "\"")
	dash := strings.Trim(string(jsonWidget["dash"]), "\"")

	widget := model.CreateWidget(dash, name, device, source, int32(typeWidget))

	err = model.InsertToDB(widget)

	// subscribe to specific channel
	webServerToMqtt <- MqttChannel{strings.Join([]string{"", device, source}, "/"), []byte(""), Sub}

	if err != nil {
		log.Println(err)
		sendJSONMessage(w, "error", "make sure that you didn't duplicate the widget")
	} else {
		sendJSONMessage(w, "true", "The Widget Has been added Successfully")
	}
}

func deleteWidget(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	dash := ps.ByName("dash")
	device := ps.ByName("device")
	name := ps.ByName("name")
	// this time we have to provide more than field to filter
	widget, err := model.GetDocument(map[string]interface{}{"id_dash": dash, "id_device": device, "name": name}, model.Widget{})

	if err != nil {
		log.Println(err)
		return
	}

	err = model.DeleteFromDB(model.ToWidget(widget))

	if err != nil {
		log.Println(err)
		sendJSONMessage(w, "error", err.Error())
	} else {
		sendJSONMessage(w, "true", "The Widget Has been added Successfully")
	}
}

func collection(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}
	colls := model.AllToCollection(
		model.GetAllDocuments(model.Filter("email", user.Email), model.Collection{}))
	userDevices := model.AllToDevice(model.GetAllDocuments(model.Filter("email", user.Email), model.Device{}))

	data := struct {
		User    *model.User
		Colls   []model.Collection
		Devices []model.Device
	}{
		user,
		colls,
		userDevices,
	}

	generatePage(w, "list_of_collections", data)
}

func getCollection(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user, err := getUserDirectly(r)
	if user == nil {
		http.Redirect(w, r, "/login", 302)
	}

	if err != nil {
		log.Println(err.Error())
	}

	// get data from database based on email
	coll, err := model.GetDocument(model.Filter("id_collector", ps.ByName("collection")), model.Collection{})

	if err != nil {
		log.Println(err.Error())
	}
	// TODO: optimize this lines
	data := struct {
		Host string
		User *model.User
		Coll model.Collection
	}{
		r.Host,
		user,
		model.ToCollection(coll),
	}

	generatePage(w, "collection", data)
}

func addCollection(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func deleteCollection(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

//internal functions
func getUserDirectly(r *http.Request) (*model.User, error) {
	userFromCookie := utils.GetCookie("kolibrio-user", r)
	if userFromCookie == nil {
		return nil, nil
	}

	user, err := model.GetDocument(
		model.Filter("email", userFromCookie["email"]),
		model.User{})

	if err != nil {
		log.Println(err.Error())
	}

	//TODO: Look for good solution
	return model.PUser(model.ToUser(user)), err
}

func sendJSONMessage(w http.ResponseWriter, status, message string) error {
	return generateJSON(w, map[string]interface{}{"status": status, "message": message})
}

//sensor websocket controller
// Notes: each websocket connection has its own go routine
func sensor(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("WS/SENSOR OPENED")
	IDdash := ps.ByName("dash")

	// we need to get a list of topics to be used in subscribing
	widgets := model.AllToWidget(model.GetAllDocuments(model.Filter("id_dash", IDdash), model.Widget{}))
	topics := make([]string, 0)

	for _, widget := range widgets {
		if widget.WidgetType == 1 {
			topics = append(topics, strings.Join([]string{"", widget.IDDEVICE, widget.Source}, "/"))
		}
	}

	// create ws
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	closeChan := make(chan bool)
	// each websocket connection has own writer
	go writeWS(ws, closeChan, topics)
	go reader(ws, closeChan, topics)
}

//actuator websocket controller
func actuator(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("WS/ACTUATOR OPENED")
	IDdash := ps.ByName("dash")

	// we need to get a list of topics to be used in subscribing
	widgets := model.AllToWidget(model.GetAllDocuments(model.Filter("id_dash", IDdash), model.Widget{}))
	topics := make([]string, 0)

	// subscribe to channel related to the actuators
	for _, widget := range widgets {
		if widget.WidgetType == 2 {
			topics = append(topics, strings.Join([]string{"", widget.IDDEVICE, widget.Source, "st"}, "/"))
		}
	}

	// create ws
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

	closeChan := make(chan bool)
	// each websocket connection has own writer
	// to write to the client
	go writeWS(ws, closeChan, topics)
	// to read from the client
	go reader(ws, closeChan, topics)
}

//WebSocket Writer
func writeWS(ws *websocket.Conn, close chan bool, topics []string) {
	// close websocket when goes out of this
	// routine
	defer func() {
		log.Println("WS/SENSOR CLOSED")
		ws.Close()
	}()

	// get the address memory of underlyingConnection of websocket
	wsAddress := fmt.Sprintf("%p", ws.UnderlyingConn())
	var listTopic []*topicObserver

	// subscribe to all topics
	for i, topic := range topics {
		listTopic = append(listTopic, &topicObserver{ID: wsAddress, topic: topic,
			msgChan: make(chan MqttChannel), closeChan: make(chan bool)})
		listener.register(listTopic[i])
	}

	webSockeMqttChannel := make(chan MqttChannel)

	// create a seprated go routines for each topic
	for i := 0; i < len(listTopic); i++ {
		go func(index int, channelOut chan MqttChannel) {
			defer log.Println(index, " -- CLOSED")
			log.Println("START GO ROUTINE ", index)
			for {
				select {
				// get messages from broker based on topic
				// send it back to outside for preparing and writing
				// to the client
				case msg := <-listTopic[index].msgChan:
					channelOut <- msg

				// if user close the connection
				// we deregister all the topic
				case <-listTopic[index].closeChan:
					listener.deregister(listTopic[index])
					return
				}
			}
		}(i, webSockeMqttChannel)
	}

	for {
		select {
		case msg := <-webSockeMqttChannel:
			msgToSend := struct {
				Topic string `json:"topic"`
				Msg   string `json:"msg"`
			}{
				msg.topic,
				string(msg.payload),
			}

			ws.WriteJSON(msgToSend)
		case <-close:
			return
		}
	}
}

// Reader Socket
func reader(ws *websocket.Conn, close chan bool, topics []string) {
	defer func() {
		log.Println("READER CLOSE")
		ws.Close()
	}()
	// ws.SetReadLimit(512)
	// ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	// ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			// get the address memory of underlyingConnection of websocket
			log.Println(err)

			wsAddress := fmt.Sprintf("%p", ws.UnderlyingConn())
			for _, topic := range topics {
				listener.closeObserverGoRoutine(topic, wsAddress)
			}
			close <- true
			return
		}

		type Message struct {
			Topic   string `json:"topic"`
			Payload string `json:"msg"`
		}

		msg := &Message{}
		err = json.Unmarshal(data, msg)

		if err != nil {
			log.Println(err)
		}

		webServerToMqtt <- MqttChannel{msg.Topic, []byte(msg.Payload), Pub}
		// log.Println(msg)
	}
}
