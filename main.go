package main

import (
	pubnub "github.com/pubnub/go"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)   // Keep alive.
	subscribed := make(chan bool) // Wait until PubNub is connected.

	config := pubnub.NewConfig()
	config.PublishKey = "pub-c-2af91346-bb0a-4bd6-89b3-d0daf9838bd9"   // YOUR PUBNUB PUBLISH KEY HERE.
	config.SubscribeKey = "sub-c-ca7555cc-d8b1-11e9-aa3a-6edd521294c5" // YOUR PUBNUB SUBSCRIBE KEY HERE.
	channel := "chat"                                                  // Chat channel.
	pn := pubnub.NewPubNub(config)

	box := js.Global().Get("document").Call("getElementById", "box")     // Output.
	input := js.Global().Get("document").Call("getElementById", "input") // Input.
	send := js.Global().Get("document").Call("getElementById", "send")   // Send button.

	listener := pubnub.NewListener() // Create a listener and subscribe.
	go func() {
		for {
			select {
			case status := <-listener.Status:
				switch status.Category {
				case pubnub.PNConnectedCategory:
					// Connected to PubNub.
					subscribed <- true
				}
			case message := <-listener.Message:
				// New message.
				if msg, ok := message.Message.(map[string]interface{}); ok {
					if str, ok := msg["msg"].(string); ok {
						box.Set("innerHTML", box.Get("innerHTML").String()+str+"<br>")
					}
				}
			}
		}
	}()
	pn.AddListener(listener)
	pn.Subscribe().
		Channels([]string{channel}).
		Execute()

	<-subscribed // Ready to receive messages.

	send.Set("onclick", js.FuncOf(func(js.Value, []js.Value) interface{} { // Publish message from input.
		msg := map[string]interface{}{
			"msg": input.Get("value").String(),
		}
		go pn.Publish().Channel(channel).Message(msg).Execute()
		return nil
	}))

	<-c
}
