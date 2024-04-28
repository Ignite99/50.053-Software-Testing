package Django

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

type OutputCriteria struct {
	ContentType            string
	StatusCode             int
	ResponseBody           string
	ResponseBodyProperties []string
}

type InputCriteria struct {
	Path                  string
	Method                string
	ContentType           string
	RequestBodyProperties []string // TODO: rename this
}

type Json_seed struct {
	Data          map[string]interface{}
	Key_to_mutate string
	Energy        int
	OC            OutputCriteria
	IC            InputCriteria
}

type ErrorQ map[int]map[string][]Json_seed

var histQ []Json_seed
var errorQ ErrorQ

func RequestParser(path string, method string, contentType string, requestBodyProperties []string) InputCriteria {
	ic := InputCriteria{
		Path:                  path,
		Method:                method,
		ContentType:           contentType,
		RequestBodyProperties: requestBodyProperties,
	}
	return ic
}

// func ResponseParser(response http.Response) (string, int, []string, string) {
func ResponseParser(response http.Response) OutputCriteria {
	// fmt.Println("-------------------------")
	// fmt.Println(response)
	// fmt.Println("-------------------------")

	// get response content type
	contentType := response.Header.Get("Content-Type")
	fmt.Println(" ++ Content Type:", contentType)

	// get response status code
	statusCode := response.StatusCode
	fmt.Println("++ Status Code:", statusCode)

	// get response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	responseBodyStr := string(responseBody)
	fmt.Println("Response Body:", responseBodyStr)

	// var respBodyProperties []string

	// if strings.Contains(contentType, "application/json") {
	// 	respBodyProperties = getJSONKeys(responseBody)
	// }
	respBodyProperties := getJSONKeys(responseBody)
	fmt.Println("++ respBodyProperties:", respBodyProperties)

	oc := OutputCriteria{
		ContentType:            contentType,
		StatusCode:             statusCode,
		ResponseBody:           responseBodyStr,
		ResponseBodyProperties: respBodyProperties,
	}

	// return contentType, statusCode, respBodyProperties, responseBodyStr
	return oc
}

func getJSONKeys(respBody []byte) []string {
	// Unmarshal JSON into an interface{}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(respBody, &jsonData); err != nil {
		return nil
	}

	// Extract keys from the JSON data
	var keys []string
	for key := range jsonData {
		keys = append(keys, key)
	}
	return keys

}

// func CheckIsInteresting(currSeed Json_seed, prevSeed Json_seed, errorQ []Json_seed) bool {
func CheckIsInteresting(currSeed Json_seed) bool {

	// 0. Init  isInteresting
	isInteresting := false

	// 1. Check if json seed exists in history Q:
	isSeedExists := seedExistsHistQ(currSeed, histQ)
	if !isSeedExists {
		histQ = append(histQ, currSeed)
		isInteresting = true
	}

	// If error, then check if error exists.
	if ((currSeed.OC.StatusCode/100)%10) == 4 || ((currSeed.OC.StatusCode/100)%10) == 5 {
		isExistingError := seedExistsInErrorQ(currSeed)

		if !isExistingError {
			isInteresting = true
		}
	}
	return isInteresting
}

func seedExistsHistQ(currSeed Json_seed, histQ []Json_seed) bool {

	for _, seed := range histQ {
		if reflect.DeepEqual(seed, currSeed) {
			return true
		}
	}
	return false
}

func seedExistsInErrorQ(currSeed Json_seed) bool {
	httpCode := currSeed.OC.StatusCode
	respBody := currSeed.OC.ResponseBody
	isExistingError := false

	// Initialize ErrorQ if it's nil
	if errorQ == nil {
		errorQ = make(ErrorQ)
	}

	// Initialize ErrorQ[httpCode] if it's nil
	if errorQ[httpCode] == nil {
		errorQ[httpCode] = make(map[string][]Json_seed)
	}

	// Initialize ErrorQ[httpCode][resBody] if it's nil
	if errorQ[httpCode][respBody] == nil {
		errorQ[httpCode][respBody] = make([]Json_seed, 0)
	}

	// Check if currSeed already exists in ErrorQ[httpCode][respBody]
	found := false
	for _, seed := range errorQ[httpCode][respBody] {
		if reflect.DeepEqual(seed, currSeed) {
			found = true
			isExistingError = true
			break
		}
	}
	// If currSeed is not found, append it to ErrorQ[httpCode][respBody]
	if !found {
		errorQ[httpCode][respBody] = append(errorQ[httpCode][respBody], currSeed)
	}
	return isExistingError
}
