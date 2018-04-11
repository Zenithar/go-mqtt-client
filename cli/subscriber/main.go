package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/davecgh/go-spew/spew"
	MQTT "github.com/eclipse/paho.mqtt.golang"

	"go.zenithar.org/mqtt/pkg/tlsconfig"
)

var (
	broker   = flag.String("broker", "127.0.0.1:8883", "Broker address")
	clientID = flag.String("client-id", "go-subscriber", "Client identifier")
	certPath = flag.String("cert", "cert.pem", "Client certificate to use")
	keyPath  = flag.String("key", "key.pem", "Client private key ")
	caPath   = flag.String("ca", "ca.pem", "Root CA")
	topic    = flag.String("topic", "livingroom/#/temperature", "Topic to subscribe")
)

func init() {
	flag.Parse()
}

func topicListener(client MQTT.Client, msg MQTT.Message) {
	spew.Dump(msg.Payload())
}

func main() {

	// Prepare client options
	opts := MQTT.NewClientOptions().AddBroker(*broker)
	opts.SetAutoReconnect(true)
	opts.SetClientID(*clientID)
	opts.SetTLSConfig(tlsconfig.ClientDefault(*certPath, *keyPath, *caPath))

	// Build MQTT Client
	cli := MQTT.NewClient(opts)

	// Connect to MQTT Server
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	// Subscribe to topic
	if token := cli.Subscribe(*topic, 0, topicListener); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

	// Wait for ctrl+c
	log.Println("Waiting for messages ... (ctrl+c to interrupt)")
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
