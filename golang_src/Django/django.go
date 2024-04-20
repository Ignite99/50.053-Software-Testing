package Django

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	logger "github.com/50.053-Software-Testing/HTML_Logger"
	fuzzer "github.com/50.053-Software-Testing/fuzzer/json_mutator"
)

// alternative x
type FakeOpenApi struct {
	Path *Paths
}

type Paths map[string]*PathItem

type PathItem struct {
	ExtensionProps
	Ref         string
	Summary     string
	Description string
	Connect     *Operation
	Delete      *Operation
	Get         *Operation
	Head        *Operation
	Options     *Operation
	Patch       *Operation
	Post        *Operation
	Put         *Operation
	Trace       *Operation
	Servers     Servers
	Parameters  Parameters
}

// Coverage records the test coverage level of each path.
type Coverage struct {
	Levels []int
}

// Returns the concatenated string representing the coverage levels
func (c *Coverage) String() string {
	r := ""
	for _, level := range c.Levels {
		r += strconv.Itoa(level) + " "
	}
	return r
}

type ResponseInfo struct {
	// request *base.Node
	Request string
	Code    int
	Type    string
	Body    string
}

var loggerInstance *logger.HTMLLogger

type json_seed struct {
	data          map[string]interface{}
	key_to_mutate string
	energy        int
}

func responseFileInit(path string) (*os.File, error) {
	outputFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening output file: %v", err)
	}

	return outputFile, nil
}

func isInteresting(line string, httpCode int) bool {
	if strings.Contains(line, "\"success\": false") || httpCode != 200 {
		return true
	}

	return false
}

func getLastLine(filename string) (string, error) {
	var lastLine string

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
		fmt.Printf("%s request succeeded! HTTP Status: %d\n", requestType, httpCode)
		fmt.Printf("Row: %s\n", row)
		loggerInstance.AddRowWithStyle("background-color:tomato", row)

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

func requestSender(outputFile *os.File, requestType string, body string, url string) (int, error, *ResponseInfo) {
	var httpCode int
	var req *http.Request
	var err error

	info := &ResponseInfo{
		Request: "req string",
		Code:    httpCode,
		Type:    "Header.Get(content-type)",
		Body:    "resBody",
		// request: node,
		// Code:    res.StatusCode,
		// Type:    res.Header.Get("Content-Type"),
		// Body:    resBody,
	}

	client := &http.Client{}

	if requestType == "GET" {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	} else if requestType == "POST" {
		req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		// return httpCode, fmt.Errorf("invalid request type: %s", requestType)
		return httpCode, fmt.Errorf("invalid request type: %s", requestType), info
	}

	if err != nil {
		return httpCode, fmt.Errorf("error creating HTTP request: %v", err), info
	}

	// Do da good request shit heheheheheh
	resp, err := client.Do(req)
	if err != nil {
		return httpCode, fmt.Errorf("error performing HTTP request: %v", err), info
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return httpCode, err, info
	}

	// Get the http request shit
	httpCode = resp.StatusCode
	checkResponse(httpCode, requestType, body, outputFile, resp)

	// Populate the ResponseInfo struct
	info.Code = resp.StatusCode
	info.Type = resp.Header.Get("Content-Type")
	info.Body = string(respBody)

	return httpCode, nil, info
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
	var inputQ []json_seed
	var filename string
	var responseFile *os.File

	// Create html logger method
	footerFilePath := "./HTML_Logger/formats/footer.html"
	htmlFileInit()

	// Set default filename if outputFilePath is empty
	if output_file_path == "" || output_file_path == "./" {
		filename = "./fuzzing_responses/response.txt"
	} else {
		filename = output_file_path
	}

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

			httpCode, err, info := requestSender(responseFile, request_type, jsonString, url)
			if err != nil {
				fmt.Println("FUCK IT WE BALLING IN REQUEST SENDER AND DIE", err)
				break
			}

			resString, err := getLastLine(filename)
			if err != nil {
				fmt.Println("FUCK IT WE BALLING IN LAST LINE AND DIE", err)
				break
			}

			// mapCodes := map[int]int{}
			// We need to store the response information here:
			mapInfos := map[int][]*ResponseInfo{}
			mapInfos[node.Group] = append(mapInfos[node.Group], info)

			// Get test coverage levels
			cov := getCoverageLevels(mapInfos)
			endCov := Coverage{}
			strictMode := false

			// Compare with each test coverage level
			isIncrease, newCov := isIndividualIncrease(cov.Levels, endCov.Levels, strictMode)

			if isIncrease {

				// Update coverage levels
				endCov.Levels = newCov.Levels

				// if guided {
				// 	// Sava as new corpus
				// 	b, err := proto.Marshal(x.grammar)
				// 	if err != nil {
				// 		panic(err)
				// 	}
				// 	x.corpus.Add(gofuzz.Artifact{Data: b})
				// }
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

	if err := loggerInstance.CloseFile(footerFilePath); err != nil {
		log.Fatalf("failed to close output file: %v", err)
	}

	responseFile.Close()

}
