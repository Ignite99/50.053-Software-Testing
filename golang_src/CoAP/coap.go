package coap

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
)

type CoAPFuzzer struct {
	target_ip        string
	target_port      int
	target_paths     []string
	total_test_cases int
	total_bugs_found int
}

func (fuzzer *CoAPFuzzer) get_paths() {
	// create a coap request and send to .well-known/core
	// get the response and parse the response to get the paths
	// return the paths
	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	path := ".well-known/core"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	body, err := resp.ReadBody()
	responseString := string(body)
	log.Printf("Response: %v", responseString)

	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindAllStringSubmatch(responseString, -1)

	for _, match := range matches {
		fuzzer.target_paths = append(fuzzer.target_paths, match[1])
	}
	log.Printf("Paths: %v", fuzzer.target_paths)
}

func (fuzzer *CoAPFuzzer) send_get_request(path string) {
	fuzzer.total_test_cases++
	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		fuzzer.total_bugs_found++
	}
	body, err := resp.ReadBody()
	responseString := string(body)

	log.Printf("GET Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	co.Close()
}

func (fuzzer *CoAPFuzzer) send_post_request(path string, payload string) {
	fuzzer.total_test_cases++
	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	content_format := message.TextPlain
	msg_payload := bytes.NewReader([]byte(payload))

	resp, err := co.Post(ctx, path, content_format, msg_payload)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		fuzzer.total_bugs_found++
	}
	body, err := resp.ReadBody()
	responseString := string(body)

	log.Printf("Post Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	co.Close()
}

func (fuzzer *CoAPFuzzer) send_put_request(path string, payload string) {
	fuzzer.total_test_cases++
	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	content_format := message.TextPlain
	msg_payload := bytes.NewReader([]byte(payload))

	resp, err := co.Put(ctx, path, content_format, msg_payload)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		fuzzer.total_bugs_found++
	}
	body, err := resp.ReadBody()
	responseString := string(body)

	log.Printf("PUT Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	co.Close()
}

func (fuzzer *CoAPFuzzer) send_delete_request(path string) {
	fuzzer.total_test_cases++
	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	resp, err := co.Delete(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		fuzzer.total_bugs_found++
	}
	body, err := resp.ReadBody()
	responseString := string(body)

	log.Printf("DELETE Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	co.Close()
}

func (fuzzer *CoAPFuzzer) run_fuzzer(path string, payload string) {
	// send a GET request
	fuzzer.send_get_request(path)

	// send a DELETE request
	fuzzer.send_delete_request(path)

	// send a POST request
	fuzzer.send_post_request(path, payload)

	// send a PUT request
	fuzzer.send_put_request(path, payload)

	for i := 0; i < 1; i++ {
		mutated_payload := mutate_add_byte(payload)
		// send a GET request
		fuzzer.send_get_request(path)

		// send a DELETE request
		fuzzer.send_delete_request(path)

		// send a POST request
		fuzzer.send_post_request(path, mutated_payload)

		// send a PUT request
		fuzzer.send_put_request(path, mutated_payload)
	}
}

func mutate_payload(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// perform byte manipulation on the payloadBytes
	// e.g., reverse the bytes
	for i, j := 0, len(payloadBytes)-1; i < j; i, j = i+1, j-1 {
		payloadBytes[i], payloadBytes[j] = payloadBytes[j], payloadBytes[i]
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_add_byte(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// add a random byte at a random position in the payloadBytes
	rand.Seed(time.Now().UnixNano())
	randomByte := byte(rand.Intn(256))
	randomPosition := rand.Intn(len(payloadBytes) + 1)
	payloadBytes = append(payloadBytes[:randomPosition], append([]byte{randomByte}, payloadBytes[randomPosition:]...)...)

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func CoAPTestDriver(ip_addr string, port int) {
	fuzzer := CoAPFuzzer{target_ip: ip_addr, target_port: port, total_test_cases: 0}
	fuzzer.get_paths()
	payload := "Hello World"

	// fuzz the target
	for _, path := range fuzzer.target_paths {
		fuzzer.run_fuzzer(path, payload)
	}

	log.Printf("Total test cases: %d", fuzzer.total_test_cases)
	log.Printf("Total bugs found: %d", fuzzer.total_bugs_found)
}
