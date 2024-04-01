package Django

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	fuzzer "github.com/50.053-Software-Testing/fuzzer/json_mutator"
)

type json_seed struct {
	data          map[string]interface{}
	key_to_mutate string
	energy        int
}

func requestSender(outputFilePath string, requestType string, body string, url string) (int, error) {
	var httpCode int

	// Create a new HTTP client
	client := &http.Client{}

	// Create a new HTTP request
	var req *http.Request
	var err error

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

	// Perform the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		return httpCode, fmt.Errorf("error performing HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Get the HTTP status code
	httpCode = resp.StatusCode

	// Print the response body
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
		curSeed := inputQ[0]
		inputQ = inputQ[1:]

		for i := 0; i < curSeed.energy; i++ {
			curSeed.data = fuzzer.MutateRequests("", curSeed.data)
			inputQ = append(inputQ, curSeed)
			jsonData, err := json.Marshal(curSeed.data)
			if err != nil {
				fmt.Println("Error marshalling JSON:", err)
				return
			}
			jsonString := string(jsonData)

			httpCode, err := requestSender(output_file_path, request_type, jsonString, url)
			if err != nil {
				fmt.Println("FUCK IT WE BALLING AND DIE", err)
				return
			}
			fmt.Println(httpCode)

			// resString := getLastLine("")
			// if !isInteresting(resString, httpCode) {
			// 	// Not interesting, so remove the new mutated input
			// 	inputQ = inputQ[:len(inputQ)-1]
			// }
		}
	}

	fmt.Println(accumulated_iterations)
}
