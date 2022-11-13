package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ahmedalkabir/kolibrio-infrastrucutre/collector"

	"github.com/ahmedalkabir/kolibrio-infrastrucutre/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type subject interface {
	register(Observer observer)
	deregister(Observer observer)
	notifyAll()
}

type observer interface {
	update(interface{})
	getTopic() string
	getID() string
	close()
}

// an experiement way
type influxStreamer struct {
	observerList map[string]map[string]observer
	name         string
}

func newInfluxStreamer(name string) *influxStreamer {
	return &influxStreamer{
		observerList: make(map[string]map[string]observer),
		name:         name,
	}
}

// to register topic to mqttStreamer
func (m *influxStreamer) register(o observer) {
	id, ok := m.observerList[o.getTopic()]
	if !ok {
		id = make(map[string]observer)
		m.observerList[o.getTopic()] = id
	}
	id[o.getID()] = o
}

func (m *influxStreamer) deregister(o observer) {
	// delete id
	delete(m.observerList[o.getTopic()], o.getID())
	// delete topic
	// place that cause abnormal event
	// yep, it's really stupid idea to
	// remove topic because it will cause confliction
	// in synchronization of connection

	// TODO: figure out the way of deleting
	// topics in efficient way
	// delete(m.observerList, o.getTopic())
}

func (m *influxStreamer) notifyObserver(topic string, msg interface{}) {
	for _, observer := range m.observerList[topic] {
		observer.update(msg)
	}
}

type mqttStreamer struct {
	observerList map[string]map[string]observer
	name         string
}

func newStreamer(name string) *mqttStreamer {
	return &mqttStreamer{
		observerList: make(map[string]map[string]observer),
		name:         name,
	}
}

// to register topic to mqttStreamer
func (m *mqttStreamer) register(o observer) {
	// subscribe to the topic provided with observer
	webServerToMqtt <- MqttChannel{o.getTopic(), []byte(""), Sub}
	id, ok := m.observerList[o.getTopic()]
	if !ok {
		id = make(map[string]observer)
		m.observerList[o.getTopic()] = id
	}
	id[o.getID()] = o
}

func (m *mqttStreamer) deregister(o observer) {
	// delete id
	delete(m.observerList[o.getTopic()], o.getID())
	// delete topic
	// place that cause abnormal event
	// yep, it's really stupid idea to
	// remove topic because it will cause confliction
	// in synchronization of connection

	// TODO: figure out the way of deleting
	// topics in efficient way
	// delete(m.observerList, o.getTopic())
}

func (m *mqttStreamer) notifyObserver(topic string, msg interface{}) {
	for _, observer := range m.observerList[topic] {
		observer.update(msg)
	}
}

func (m *mqttStreamer) closeObserverGoRoutine(topic, id string) {
	observer, ok := m.observerList[topic][id]
	if ok {
		observer.close()
	}
}

//influxObserver
type influxObserver struct {
	Email string
	Topic string
}

// write the point to the influx db
func (influx *influxObserver) update(value interface{}) {
	// do some work here
	// we'll received value in form of json
	// unmarshal the received value
	type _value struct {
		Value float32 `json:"value"`
	}
	concerteValue := _value{}
	json.Unmarshal(value.([]byte), &concerteValue)
	// write it to the influxDB
	collector.NewPoint(influx.getID(), map[string]string{"topic": influx.getTopic()},
		map[string]interface{}{"value": concerteValue.Value}).Write()
}

func (influx *influxObserver) getTopic() string {
	return influx.Topic
}

func (influx *influxObserver) getID() string {
	return influx.Email
}

func (influx *influxObserver) close() {

}

//

type topicObserver struct {
	ID        string
	topic     string
	msgChan   chan MqttChannel
	closeChan chan bool
}

func (t *topicObserver) update(msg interface{}) {
	// log.Println("OBSERVER -- ", t.getID(), " - ", msg)
	msgToSend := msg.(MqttChannel)
	t.msgChan <- msgToSend
}

func (t *topicObserver) close() {
	t.closeChan <- true
}

func (t *topicObserver) getTopic() string {
	return t.topic
}

func (t *topicObserver) getID() string {
	return t.ID
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func startMqttService(listener *mqttStreamer, influxStream *influxStreamer, inChannel chan MqttChannel, outChannel chan MqttChannel) {
	broker := fmt.Sprintf("tcp://%s:%s", config.MqttConf.Address, config.MqttConf.Port)
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(config.MqttConf.ID)

	// opts.SetKeepAlive(2 * time.Second)
	// opts.SetDefaultPublishHandler(f)
	// opts.SetPingTimeout(1 * time.Second)

	clinet := mqtt.NewClient(opts)
	if token := clinet.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		select {
		case msg1 := <-inChannel:
			switch msg1._Type {
			case Pub:
				token := clinet.Publish(msg1.topic, 0, false, string(msg1.payload))
				token.Wait()
			case Sub:
				log.Println("SUB TO --- ", msg1.topic)
				if token := clinet.Subscribe(msg1.topic, 0, func(client mqtt.Client, msg mqtt.Message) {
					// log.Println("REC ", msg.Topic(), "--", msg.Payload())
					listener.notifyObserver(msg.Topic(), MqttChannel{msg.Topic(), msg.Payload(), Recv})
					influxStream.notifyObserver(msg.Topic(), msg.Payload())
					// pass it to influxdb asynchonosly

				}); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
					os.Exit(1)
				}
			case UnSub:
				log.Println("UNSUB FROM --- ", msg1.topic)
				clinet.Unsubscribe(msg1.topic)
			}
		}
	}
}

func listenChannel() chan MqttChannel {
	out := make(chan MqttChannel)
	go func() {
		for {
			select {
			case msg := <-middleawareToWebSocket:
				out <- msg
			}
		}
	}()
	return out
}
