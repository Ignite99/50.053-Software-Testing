package coap

import (
	"context"
	"log"
	"time"

	"github.com/plgd-dev/go-coap/v3/udp"
)

const URL = "127.0.0.1:5683"

func ClientConnect(url string, endpoint string) {
	co, err := udp.Dial(url)

	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.Get(ctx, endpoint)
	if err != nil {
		log.Fatalf("Cannot get response: %v", err)
		return
	}
	log.Printf("Response: %+v", resp)

}

func CoAPTestDriver() {
	ClientConnect(URL, "/basic/")
}
