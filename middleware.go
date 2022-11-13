package main

import "sync"

type Type int

const (
	None Type = iota
	Pub
	Sub
	UnSub
	Recv
)

// MqttChannel is going to use as message communication between
// middleware and mqtt
type MqttChannel struct {
	topic   string
	payload []byte
	_Type   Type
}

var wsMap sync.Map

func startMiddleWareService(inChannel chan MqttChannel, outChannels ...chan MqttChannel) {

	// receive messages from mqtt
	go func() {
		for {
			select {
			// mqtt to websocket
			case msg1 := <-inChannel:
				// a hany a solution for our issues because
				// that channel cann't share between all
				// goroutines
				// wsMap.Store("value", msg1)
				outChannels[1] <- msg1
			// webserver to mqtt
			case msg2 := <-outChannels[2]:
				outChannels[0] <- msg2
			}
		}
	}()
}
