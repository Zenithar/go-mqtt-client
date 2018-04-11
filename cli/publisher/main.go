package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

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
	retain   = flag.Bool("retain", false, "Sets retain attribute to publication")
)

func init() {
	flag.Parse()
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

	// Read message from stdin
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Unable to read data from stdin, %v\n", err)
	}

	// Publish to topic
	if token := cli.Publish(*topic, 0, *retain, data); token.Wait() && token.Error() != nil {
		log.Fatalln(token.Error())
	}

}
