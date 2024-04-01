package Django

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	fuzzer "github.com/50.053-Software-Testing/fuzzer/json_mutator"
)

type json_seed struct {
	data          map[string]interface{}
	key_to_mutate string
	energy        int
}

func isInteresting(line string, httpCode int) bool {
	if strings.Contains(line, "\"success\": false") || httpCode != 200 {
		return true
	}

	return false
}

func getLastLine(outputFilePath string) (string, error) {
	var filename string
	var lastLine string

	// Set default filename if outputFilePath is empty
	if outputFilePath == "" || outputFilePath == "./" {
		filename = "./fuzzing_responses/response.txt"
	} else {
		filename = outputFilePath
	}

	file, err := os.Open(filename)
	if err != nil {
		return lastLine, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return lastLine, err
	}

	return lastLine, nil
}

func checkResponse(httpCode int, requestType string, body string) int {
	var row []string

	now := time.Now()
	dt := now.Format(time.RFC3339)
	row = []string{dt, requestType, body, fmt.Sprint(httpCode)}

	// Log responses in HTML logger after printing out status message
	switch httpCode {
	case 200:
		fmt.Printf("%s request succeeded!\n", requestType)
		fmt.Printf("HTTP Status: %d\n", httpCode)
		fmt.Printf("Row: %s\n", row)
		// htmlLogger.addRow("background-color:palegreen", row)
		return 0
	case 201:
		fmt.Printf("%s create request succeeded!\n", requestType)
		fmt.Printf("HTTP Status: %d\n", httpCode)
		fmt.Printf("Row: %s\n", row)
		// htmlLogger.addRow("background-color:palegreen", row)
		return 0
	case 202:
		fmt.Printf("%s accept request succeeded!\n", requestType)
		fmt.Printf("HTTP Status: %d\n", httpCode)
		fmt.Printf("Row: %s\n", row)
		// htmlLogger.addRow("background-color:palegreen", row)
		return 0
	default:
		fmt.Printf("%s request failed!\n", requestType)
		fmt.Printf("HTTP Status code: %d\n", httpCode)
		fmt.Printf("Row: %s\n", row)
		// htmlLogger.addRow("background-color:tomato", row)
		return 0
	}
}

func requestSender(outputFilePath string, requestType string, body string, url string) (int, error) {
	var httpCode int
	var req *http.Request
	var err error

	client := &http.Client{}

	if requestType == "GET" {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	} else if requestType == "POST" {
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

	// Get the http request shit
	httpCode = resp.StatusCode
	checkResponse(httpCode, requestType, body)

	// idk if this is a good idea
	resBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(resBody))

	return httpCode, nil
}

func Django_Test_Driver(energy int, url string, request_type string, input_file_path string, output_file_path string) {
	var accumulated_iterations int
	var testing_incomplete bool
	var data map[string]interface{}
	var inputQ []json_seed

	// Create html logger method

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

	for key, _ := range data {
		seed := json_seed{
			data:          data,
			key_to_mutate: key,
			energy:        3,
		}
		inputQ = append(inputQ, seed)
	}

	for testing_incomplete {
		if len(inputQ) == 0 {
			break
		}

		curSeed := inputQ[0]
		inputQ = inputQ[1:]

		for i := 0; i < curSeed.energy; i++ {
			curSeed.data = fuzzer.MutateRequests(request_type, curSeed.data)
			inputQ = append(inputQ, curSeed)
			jsonData, err := json.Marshal(curSeed.data)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			jsonString := string(jsonData)

			httpCode, err := requestSender(output_file_path, request_type, jsonString, url)
			if err != nil {
				fmt.Println("FUCK IT WE BALLING IN REQUEST SENDER AND DIE", err)
				break
			}

			resString, err := getLastLine(output_file_path)
			if err != nil {
				fmt.Println("FUCK IT WE BALLING IN LAST LINE AND DIE", err)
				break
			}

			if !isInteresting(resString, httpCode) {
				// Not interesting, so remove the new mutated input
				inputQ = inputQ[:len(inputQ)-1]
			}
		}

		accumulated_iterations++

		if accumulated_iterations > 10 {
			testing_incomplete = false
			break
		}
	}

}
