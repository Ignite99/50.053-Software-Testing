package coap

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	logger "github.com/50.053-Software-Testing/HTML_Logger"
	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/udp"
)

var fuzzingLogger *logger.HTMLLogger
var uniqueLogger *logger.HTMLLogger

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

var uniqueQ []Seed
var inputQ []Seed

type CoAPFuzzer struct {
	target_ip               string
	target_port             int
	target_paths            []string
	total_test_cases        int
	total_bugs_found        int
	interesting_cases_found int
}

// initialise html file for fuzzingLogger
func fuzzingLoggerInit() {
	var columnNames []string

	outputFilePath := "./fuzzing_responses/"
	outputFileName := "logs.html"
	projectType := "COAP"

	outputFile, err := os.Create(filepath.Join(outputFilePath, outputFileName))
	if err != nil {
		// log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Call constructor
	fuzzingLogger = logger.NewHTMLLogger(outputFilePath, outputFileName, projectType, outputFile)

	fuzzingLogger.CreateFile()

	fuzzingLogger.AddText("text-align:center; font-size:26px;", "CoAP Fuzzer Output")

	// Initialise headings
	columnNames = []string{"Time", "Path", "Method", "Request Payload", "Response Body", "Response Payload", "Message Type", "CoAP Code"}
	fuzzingLogger.CreateTableHeadings("background-color:lightgrey", columnNames)

	fmt.Println("Fuzzing HTML logger created and used successfully.")
}

// initialise html file for uniqueLogger
func uniqueLoggerInit() {
	var columnNames []string

	outputFilePath := "./fuzzing_responses/"
	outputFileName := "unique_logs.html"
	projectType := "COAP"

	outputFile, err := os.Create(filepath.Join(outputFilePath, outputFileName))
	if err != nil {
		// log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// Call constructor
	uniqueLogger = logger.NewHTMLLogger(outputFilePath, outputFileName, projectType, outputFile)

	uniqueLogger.CreateFile()

	uniqueLogger.AddText("text-align:center; font-size:26px", "CoAP Unique Responses List")

	// Initialise headings
	columnNames = []string{"Time", "Path", "Method", "Request Payload", "Response Body", "Response Payload", "Message Type", "CoAP Code"}
	uniqueLogger.CreateTableHeadings("background-color:lightgrey", columnNames)

	fmt.Println("Error HTML logger created and used successfully.")
}

func (fuzzer *CoAPFuzzer) IsInteresting(currSeed Seed) {
	// if interesting, add to inputQ
	if CheckIsInteresting(currSeed, inputQ) {
		fuzzer.interesting_cases_found++
		inputQ = append(inputQ, currSeed)
	}
}

func (fuzzer *CoAPFuzzer) get_paths() {
	// create a coap request and send to .well-known/core
	// get the response and parse the response to get the paths
	// return the paths
	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		// log.Fatalf("Error dialing: %v", err)
	}

	path := ".well-known/core"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		// log.Fatalf("Error sending request: %v", err)
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

func (fuzzer *CoAPFuzzer) send_get_request(currSeed Seed) {
	fuzzer.total_test_cases++
	path := currSeed.IC.Path

	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		// log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		// log.Fatalf("Error sending request: %v", err)
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
	row := []string{currTime, path, "GET", "", resp.String(), responseString, typeResponse, codeResponse}
	if CheckCodeSuccess(codeResponse) {
		fuzzingLogger.AddRowWithStyle("background-color:honeydew", row)
	} else {
		fuzzingLogger.AddRowWithStyle("background-color:lightgreen", row)
	}

	// modify currSeed's output and input criteria
	oc := OutputCriteria{"text", codeResponse, typeResponse}
	ic := InputCriteria{path, "POST", "text"}
	currSeed.IC = ic
	currSeed.OC = oc

	// check if anything is interesting and should be put inside uniqueQ and uniqueLogger
	if CheckIsInteresting(currSeed, uniqueQ) {
		uniqueQ = append(uniqueQ, currSeed)
		uniqueLogger.AddRowWithStyle("background-color:honeydew", row)
	}

	co.Close()
}

func (fuzzer *CoAPFuzzer) send_post_request(currSeed Seed) {
	fuzzer.total_test_cases++
	path := currSeed.IC.Path
	payload := currSeed.Data

	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		// log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	content_format := message.TextPlain
	msg_payload := bytes.NewReader([]byte(payload))

	resp, err := co.Post(ctx, path, content_format, msg_payload)
	if err != nil {
		// log.Fatalf("Error sending request: %v", err)
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

	// append to HTML Logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "POST", payload, resp.String(), responseString, typeResponse, codeResponse}
	if CheckCodeSuccess(codeResponse) {
		fuzzingLogger.AddRowWithStyle("background-color:lightCyan", row)
	} else {
		fuzzingLogger.AddRowWithStyle("background-color:skyBlue", row)
	}

	// make output and input criteria
	oc := OutputCriteria{"text", codeResponse, typeResponse}
	ic := InputCriteria{path, "POST", "text"} // TODO - get type of request payload
	currSeed.IC = ic
	currSeed.OC = oc

	// check if anything is interesting and should be put inside uniqueQ and uniqueLogger
	if CheckIsInteresting(currSeed, uniqueQ) {
		uniqueQ = append(uniqueQ, currSeed)
		uniqueLogger.AddRowWithStyle("background-color:lightCyan", row)
	}

	// check if currSeed isInteresting, if yes, put in inputQ
	fuzzer.IsInteresting(currSeed)

	co.Close()
}

func (fuzzer *CoAPFuzzer) send_put_request(currSeed Seed) {
	fuzzer.total_test_cases++
	path := currSeed.IC.Path
	payload := currSeed.Data

	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		// log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	content_format := message.TextPlain
	msg_payload := bytes.NewReader([]byte(payload))

	resp, err := co.Put(ctx, path, content_format, msg_payload)
	if err != nil {
		// log.Fatalf("Error sending request: %v", err)
		fuzzer.total_bugs_found++
	}
	body, err := resp.ReadBody()
	responseString := string(body)

	codeResponse := resp.Code().String()
	typeResponse := resp.Type().String()

	log.Printf("PUT Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	// append to HTML Logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "PUT", payload, resp.String(), responseString, typeResponse, codeResponse}
	if CheckCodeSuccess(codeResponse) {
		fuzzingLogger.AddRowWithStyle("background-color:cornsilk", row)
	} else {
		fuzzingLogger.AddRowWithStyle("background-color:lightSalmon", row)
	}

	// make output and input criteria
	oc := OutputCriteria{"text", codeResponse, typeResponse}
	ic := InputCriteria{path, "PUT", "text"} // TODO - get type of request payload
	currSeed.IC = ic
	currSeed.OC = oc

	// check if anything is interesting and should be put inside uniqueQ and uniqueLogger
	if CheckIsInteresting(currSeed, uniqueQ) {
		uniqueQ = append(uniqueQ, currSeed)
		uniqueLogger.AddRowWithStyle("background-color:cornsilk", row)
	}

	// check if currSeed isInteresting, if yes, put in inputQ
	fuzzer.IsInteresting(currSeed)

	co.Close()
}

func (fuzzer *CoAPFuzzer) send_delete_request(currSeed Seed) {
	fuzzer.total_test_cases++
	path := currSeed.IC.Path

	co, err := udp.Dial(fmt.Sprintf("%s:%d", fuzzer.target_ip, fuzzer.target_port))
	if err != nil {
		// log.Fatalf("Error dialing: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	resp, err := co.Delete(ctx, path)
	if err != nil {
		// log.Fatalf("Error sending request: %v", err)
		fuzzer.total_bugs_found++
	}
	body, err := resp.ReadBody()
	responseString := string(body)

	codeResponse := resp.Code().String()
	typeResponse := resp.Type().String()

	log.Printf("DELETE Request to %s", path)
	log.Printf("Response: %v", resp.String())
	log.Printf("Response body: %v", responseString)

	// append to fuzzing logger
	currTime := time.Now().Format(time.RFC3339)
	row := []string{currTime, path, "DELETE", "", resp.String(), responseString, typeResponse, codeResponse}
	if CheckCodeSuccess(codeResponse) {
		fuzzingLogger.AddRowWithStyle("background-color:lavenderBlush", row)
	} else {
		fuzzingLogger.AddRowWithStyle("background-color:lightPink", row)
	}

	// make output and input criteria
	oc := OutputCriteria{"text", codeResponse, typeResponse}
	ic := InputCriteria{path, "DELETE", "text"} // TODO - get type of request payload
	currSeed.IC = ic
	currSeed.OC = oc

	// check if anything is interesting and should be put inside uniqueQ and uniqueLogger
	if CheckIsInteresting(currSeed, uniqueQ) {
		uniqueQ = append(uniqueQ, currSeed)
		uniqueLogger.AddRowWithStyle("background-color:lavenderBlush", row)
	}

	co.Close()
}

func (fuzzer *CoAPFuzzer) run_fuzzer(path string) {
	// check if inputQ is empty. If empty, exit program.
	if len(inputQ) == 0 {
		log.Printf("Input Queue is empty! Exiting the program...")
		footerFilePath := "./HTML_Logger/formats/footer.html"
		if err := fuzzingLogger.CloseFile(footerFilePath); err != nil {
			// log.Fatalf("failed to close output file: %v", err)
		}
		if err := uniqueLogger.CloseFile(footerFilePath); err != nil {
			// log.Fatalf("failed to close output file: %v", err)
		}
		os.Exit(0)
	}

	// take the first seed in the queue
	currSeed := inputQ[0]
	inputQ = inputQ[1:]

	currSeed.IC.Path = path

	// send a GET request
	fuzzer.send_get_request(currSeed)

	// send a DELETE request
	// fuzzer.send_delete_request(currSeed)

	// send a POST request
	fuzzer.send_post_request(currSeed)

	// send a PUT request
	fuzzer.send_put_request(currSeed)

	for i := 0; i < currSeed.Energy; i++ {
		mutated_payload := mutate_add_byte(currSeed.Data)
		currSeed.Data = mutated_payload

		// send a GET request
		fuzzer.send_get_request(currSeed)

		// send a DELETE request
		// fuzzer.send_delete_request(currSeed)

		// send a POST request
		fuzzer.send_post_request(currSeed)

		// send a PUT request
		fuzzer.send_put_request(currSeed)
	}
}

func CoAPTestDriver(ip_addr string, port int, input_file_path string, energy int) {
	fuzzer := CoAPFuzzer{target_ip: ip_addr, target_port: port, total_test_cases: 0}
	fuzzer.get_paths()

	payloadFile, err := os.Open(input_file_path)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		return
	}
	defer payloadFile.Close()

	payload, err := io.ReadAll(payloadFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	log.Printf("Energy Assigned: %d", energy)

	// append the first payload (seed input) to the inputQ, giving the specified energy
	inputQ = append(inputQ, Seed{string(payload), energy, OutputCriteria{"text", "", ""}, InputCriteria{"", "", "text"}})

	// create html instance
	fuzzingLoggerInit()
	uniqueLoggerInit()

	// fuzz the target
	for _, path := range fuzzer.target_paths {
		fuzzer.run_fuzzer(path)
	}

	// fuzz the inputQ until it is empty or until terminated by user
	for len(inputQ) > 0 {
		log.Println("Number of test cases: ", fuzzer.total_test_cases)
		log.Println("Number of bugs found: ", fuzzer.total_bugs_found)
		log.Println("Number of interesting cases found: ", fuzzer.interesting_cases_found)

		currSeed := inputQ[0]
		inputQ = inputQ[1:]

		// send a GET request
		fuzzer.send_get_request(currSeed)

		// uncomment for DELETE request
		// fuzzer.send_delete_request(currSeed)

		// send a POST request
		fuzzer.send_post_request(currSeed)

		// send a PUT request
		fuzzer.send_put_request(currSeed)

		for i := 0; i < currSeed.Energy; i++ {
			mutated_payload := mutate_random(currSeed.Data, i)
			currSeed.Data = mutated_payload

			// send a GET request
			fuzzer.send_get_request(currSeed)

			// uncomment for DELETE request
			// fuzzer.send_delete_request(currSeed)

			// send a POST request
			fuzzer.send_post_request(currSeed)

			// send a PUT request
			fuzzer.send_put_request(currSeed)
		}
	}

	log.Printf("Total test cases: %d", fuzzer.total_test_cases)
	log.Printf("Total bugs found: %d", fuzzer.total_bugs_found)

	fuzzingLogger.AddText("text-align:center;", fmt.Sprintf("Total test cases found: %d", fuzzer.total_test_cases))
	fuzzingLogger.AddText("text-align:center;", fmt.Sprintf("Total bugs found: %d", fuzzer.total_bugs_found))

	// close html instances
	footerFilePath := "./HTML_Logger/formats/footer.html"
	if err := fuzzingLogger.CloseFile(footerFilePath); err != nil {
		// log.Fatalf("failed to close output file: %v", err)
	}
	if err := uniqueLogger.CloseFile(footerFilePath); err != nil {
		// log.Fatalf("failed to close output file: %v", err)
	}
}
