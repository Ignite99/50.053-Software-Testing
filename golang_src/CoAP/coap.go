package coap

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"time"

	logger "github.com/50.053-Software-Testing/HTML_Logger"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
)

var loggerInstance *logger.HTMLLogger

type OutputCriteria struct {
	ContentType string
	StatusCode  string
	MessageType string
}

type InputCriteria struct {
	Path        string
	Method      string
	ContentType string
}

type Seed struct {
	Data   string
	Energy int
	OC     OutputCriteria
	IC     InputCriteria
}

var errorQ []Seed
var inputQ []Seed

type CoAPFuzzer struct {
	target_ip        string
	target_port      int
	target_paths     []string
	total_test_cases int
	total_bugs_found int
}

// initialise html file
func htmlFileInit() {
	var column_names []string

	outputFilePath := "./fuzzing_responses/"
	outputFileName := "logs.html"
	projectType := "CoAP"

	outputFile, err := os.Create(filepath.Join(outputFilePath, outputFileName))
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Call constructor
	htmlLogger := logger.NewHTMLLogger(outputFilePath, outputFileName, projectType, outputFile)
	loggerInstance = htmlLogger

	loggerInstance.CreateFile()

	// Initialise headings
	column_names = []string{"Time", "Path", "Method", "Request Payload", "Response Body", "Response Payload", "Message Type", "CoAP Code"}
	loggerInstance.CreateTableHeadings("background-color:lightgrey", column_names)

	fmt.Println("HTML logger created and used successfully.")
}

func (fuzzer *CoAPFuzzer) IsInteresting(currSeed Seed) {
	// if interesting, add to inputQ
	if fuzzer.total_test_cases == 0 {
		log.Printf("test case 0 - is interesting")
		inputQ = append(inputQ, currSeed)
	} else if CheckIsInteresting(currSeed, inputQ, errorQ) {
		log.Printf("is interesting")
		inputQ = append(inputQ, currSeed)
	} else {
		log.Printf("not interesting")
	}
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

	codeResponse := resp.Code().String()
	typeResponse := resp.Type().String()

	log.Printf("GET Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	// append to HTML Logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "GET", "N/A", resp.String(), responseString, typeResponse, codeResponse}
	loggerInstance.AddRowWithStyle("background-color:aquamarine", row)

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

	codeResponse := resp.Code().String()
	typeResponse := resp.Type().String()

	log.Printf("Post Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	if err != nil {
		fmt.Printf(err.Error())
	}

	// make output and input criteria
	oc := OutputCriteria{"text", codeResponse, typeResponse}
	ic := InputCriteria{path, "POST", "text"} // TODO - get type of request payload

	// make currSeed
	currSeed := Seed{payload, 3, oc, ic}

	// check if isInteresting, if yes, put in inputQ
	fuzzer.IsInteresting(currSeed)

	// append to HTML Logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "POST", payload, resp.String(), responseString, typeResponse, codeResponse}
	loggerInstance.AddRowWithStyle("background-color:paleTurquoise", row)

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

	codeResponse := resp.Code().String()
	typeResponse := resp.Type().String()

	log.Printf("PUT Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	// make output and input criteria
	oc := OutputCriteria{"text", codeResponse, typeResponse}
	ic := InputCriteria{path, "PUT", "text"} // TODO - get type of request payload

	// make currSeed
	currSeed := Seed{payload, 3, oc, ic}

	// check if isInteresting
	fuzzer.IsInteresting(currSeed)

	// append to HTML Logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "PUT", payload, resp.String(), responseString, typeResponse, codeResponse}
	loggerInstance.AddRowWithStyle("background-color:navajoWhite", row)

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

	codeResponse := resp.Code().String()
	typeResponse := resp.Type().String()

	log.Printf("DELETE Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	// append to HTML Logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "DELETE", "N/A", resp.String(), responseString, typeResponse, codeResponse}
	loggerInstance.AddRowWithStyle("background-color:lightSalmon", row)

	co.Close()
}

func (fuzzer *CoAPFuzzer) run_fuzzer(path string) {
	// check if inputQ is empty
	if len(inputQ) == 0 {
		log.Printf("Input Queue is empty! Exiting the program...")

		// close HTML logger
		footerFilePath := "./HTML_Logger/formats/footer.html"
		if err := loggerInstance.CloseFile(footerFilePath); err != nil {
			log.Fatalf("failed to close output file: %v", err)
		}

		os.Exit(0)
	}

	// take the first seed in the queue
	currSeed := inputQ[0]
	inputQ = inputQ[1:]

	// test inside inputQ
	for _, seed := range inputQ {
		fmt.Println(seed.Data)
	}

	// get payload and energy
	payload := currSeed.Data
	energy := currSeed.Energy

	// send a GET request
	fuzzer.send_get_request(path)

	// send a DELETE request
	fuzzer.send_delete_request(path)

	// send a POST request
	fuzzer.send_post_request(path, payload)

	// send a PUT request
	fuzzer.send_put_request(path, payload)

	for i := 0; i < energy; i++ {
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

	// TODO - seed input should be from user input
	// append the first payload (seed input) to the inputQ, giving an arbitrary energy of 3
	inputQ = append(inputQ, Seed{payload, 3, OutputCriteria{"text", "", ""}, InputCriteria{"", "", "text"}})

	// create html instance
	htmlFileInit()

	// fuzz the target
	for _, path := range fuzzer.target_paths {
		fmt.Println("Path: ", path)
		fuzzer.run_fuzzer(path)
	}

	log.Printf("Total test cases: %d", fuzzer.total_test_cases)
	log.Printf("Total bugs found: %d", fuzzer.total_bugs_found)

	// close html instance
	footerFilePath := "./HTML_Logger/formats/footer.html"
	if err := loggerInstance.CloseFile(footerFilePath); err != nil {
		log.Fatalf("failed to close output file: %v", err)
	}
}
