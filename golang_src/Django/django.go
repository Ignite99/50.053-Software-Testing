package Django

import (
	// "bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	logger "github.com/50.053-Software-Testing/HTML_Logger"
	fuzzer "github.com/50.053-Software-Testing/fuzzer/json_mutator"
	interesting "github.com/50.053-Software-Testing/IsInteresting"
)

var loggerInstance *logger.HTMLLogger

// var error_queue []json_seed

var inputQ []interesting.Json_seed

func responseFileInit(path string) (*os.File, error) {
	outputFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening output file: %v", err)
	}

	return outputFile, nil
}


func checkResponse(httpCode int, requestType string, body string, file *os.File, resp *http.Response) {
	var row []string

	now := time.Now()
	dt := now.Format(time.RFC3339)
	row = []string{dt, requestType, body, fmt.Sprint(httpCode)}

	// Log responses in HTML logger after printing out status message
	switch httpCode {
	case 200:
		fmt.Printf("%s request succeeded! HTTP Status: %d\n", requestType, httpCode)
		fmt.Printf("Row: %s\n", row)
		loggerInstance.AddRowWithStyle("background-color:palegreen", row)
		return
	case 201:
		fmt.Printf("%s request succeeded! HTTP Status: %d\n", requestType, httpCode)
		fmt.Printf("Row: %s\n", row)
		loggerInstance.AddRowWithStyle("background-color:palegreen", row)
		return
	case 202:
		fmt.Printf("%s request succeeded! HTTP Status: %d\n", requestType, httpCode)
		fmt.Printf("Row: %s\n", row)
		loggerInstance.AddRowWithStyle("background-color:palegreen", row)
		return
	default:
		fmt.Printf("++ failed %s %d\n", requestType, httpCode)
		fmt.Printf("%s request FAILED! HTTP Status: %d\n", requestType, httpCode)
		fmt.Printf("Row: %s\n", row)
		loggerInstance.AddRowWithStyle("background-color:tomato", row)

		// TODO - isInteresting - if an error, append to the error_queue
		// error_queue = append(error_queue, cur_seed)

		// Write the response body to the file for fucked up responses
		_, _ = file.WriteString("\n")
		_, err := io.Copy(file, resp.Body)
		if err != nil {
			fmt.Printf("error writing to output file: %v\n", err)
			return
		}

		return
	}
}

func requestSender(outputFile *os.File, requestType string, body string, url string) (int, error) {
	var httpCode int
	var req *http.Request
	var err error

	client := &http.Client{}

	if requestType == "GET" {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	} else if requestType == "POST" {
		fmt.Printf("++ in POST")

		req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		return httpCode, fmt.Errorf("invalid request type: %s", requestType)
	}

	if err != nil {
		return httpCode, fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Do da good request shit heheheheheh
	resp, err := client.Do(req)
	if err != nil {
		return httpCode, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Get the lastmost (current) seed from queue, parse responses, and add it to the current seed's output criteria
	curSeed := inputQ[len(inputQ)-1]
	curSeed.OC.ContentType, curSeed.OC.StatusCode, curSeed.OC.ResponseBody = interesting.ResponseParser(*resp)
	inputQ[len(inputQ)-1] = curSeed
	
	// Get the http request shit
	httpCode = resp.StatusCode
	checkResponse(httpCode, requestType, body, outputFile, resp)

	return httpCode, nil
}

func htmlFileInit() {
	var column_names []string

	outputFilePath := "./fuzzing_responses/"
	outputFileName := "logs.html"
	projectType := "DJANGO"

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
	column_names = []string{"Time", "Request type", "Sent Contents", "HTTP Code"}
	loggerInstance.CreateTableHeadings("background-color:lightgrey", column_names)

	fmt.Println("HTML logger created and used successfully.")
}

func Django_Test_Driver(energy int, url string, request_type string, input_file_path string, output_file_path string) {
	var accumulated_iterations int
	var testing_incomplete bool
	var data map[string]interface{}
	// var errorQ []json_seed
	// var filename string
	var responseFile *os.File

	// Create html logger method
	footerFilePath := "./HTML_Logger/formats/footer.html"
	htmlFileInit()

	// // Set default filename if outputFilePath is empty
	// if output_file_path == "" || output_file_path == "./" {
	// 	filename = "./fuzzing_responses/response.txt"
	// } else {
	// 	filename = output_file_path
	// }

	responseFile, err := responseFileInit(output_file_path)
	if err != nil {
		fmt.Println("die from no response file", err)
		return
	}

	testing_incomplete = true
	accumulated_iterations = 0

	jsonFile, err := os.Open(input_file_path)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()

	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	var contentTypeTemp []string

	for key, _ := range data {
		seed := interesting.Json_seed{
			Data:          data,
			Key_to_mutate: key,
			Energy:        3,
		}
		inputQ = append(inputQ, seed)
		contentTypeTemp = append(contentTypeTemp, key)
	}

	for testing_incomplete {
		if len(inputQ) == 0 {
			break
		}

		curSeed := inputQ[0]
		inputQ = inputQ[1:]

		for i := 0; i < curSeed.Energy; i++ {
			curSeed.Data = fuzzer.MutateRequests(request_type, curSeed.Data)

			inputQ = append(inputQ, curSeed)

			jsonData, err := json.Marshal(curSeed.Data)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			jsonString := string(jsonData)

			httpCode, err := requestSender(responseFile, request_type, jsonString, url)
			if err != nil {
				fmt.Println("FUCK IT WE BALLING IN REQUEST SENDER AND DIE", err)
				break
			}
			fmt.Printf("++ print httpcode for the lols %s", httpCode)
			
			curSeed.IC = interesting.RequestParser(url, request_type, contentTypeTemp)
			
			
			if i != 0 { // TODO,  wrong implementation since this i refers to energy.
				prevSeed := inputQ[i-1]
				var isInteresting = interesting.CheckIsInteresting(curSeed, prevSeed)

				if isInteresting == false {
					// Not interesting, so remove new mutated input.
					inputQ = inputQ[:len(inputQ)-1]
				}
			}
		}

		accumulated_iterations++

		if accumulated_iterations > 10 {
			testing_incomplete = false
			break
		}
	}

	if err := loggerInstance.CloseFile(footerFilePath); err != nil {
		log.Fatalf("failed to close output file: %v", err)
	}

	responseFile.Close()

}
