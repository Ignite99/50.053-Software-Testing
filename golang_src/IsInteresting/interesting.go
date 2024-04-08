package Interesting

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"reflect"
)

type OutputCriteria struct {
	ContentType  string
	StatusCode   int
	ResponseBody string
}

type InputCriteria struct {
	Path        string
	Method      string
	ContentType string
	RequestBodyProperties []string // TODO: rename this
}

type Json_seed struct {
	Data          map[string]interface{}
	Key_to_mutate string
	Energy        int
	OC            OutputCriteria
	IC            InputCriteria
}

func RequestParser(path string, method string, contentType string, requestBodyProperties []string) InputCriteria {
	ic := InputCriteria{Path: path, Method: method, ContentType: contentType, RequestBodyProperties: requestBodyProperties}
    return ic
}

func ResponseParser(response http.Response) (string, int, string) {
	// fmt.Println("-------------------------")
	// fmt.Println(response)
	// fmt.Println("-------------------------")

	// get response content type
	contentType := response.Header.Get("Content-Type")
	// fmt.Println("Content Type:", contentType)

	// get response status code
	statusCode := response.StatusCode
	// fmt.Println("Status Code:", statusCode)

	// get response body
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}
	responseBodyStr := string(responseBody)
	// fmt.Println("Response Body:", responseBodyStr)

	return contentType, statusCode, responseBodyStr
}

func CheckIsInteresting(currSeed Json_seed, prevSeed Json_seed, errorQ []Json_seed) bool {
	// Considered not interesting if:
	//		- prev IC same as cur IC &&
	// 		- prev OC same as cur OC && 
	// 		- curSeed exists in errorQ
	
	currIc := currSeed.IC
	prevIc := prevSeed.IC

	currOc := currSeed.OC
	prevOc := prevSeed.OC

	boolIcEqual := currIc.Path == prevIc.Path &&
	currIc.Method == prevIc.Method &&
	currIc.ContentType == prevIc.ContentType &&
	isContentTypeSame(currIc.RequestBodyProperties, prevIc.RequestBodyProperties)

	boolOcEqual := currOc.ContentType == prevOc.ContentType &&
		currOc.StatusCode == prevOc.StatusCode &&
		currOc.ResponseBody == prevOc.ResponseBody

		
	// TODO: Incomplete implementation of checking ErrorQ
	inErrorQ := isExistsInErrorQ(currSeed, errorQ )
	fmt.Printf("inErrorQ %s\n", inErrorQ)

	return !(boolIcEqual && boolOcEqual && inErrorQ)
}

func isContentTypeSame(arr1 []string, arr2 []string) bool {
	
	if len(arr1) != len(arr2) {
		return false
    }

	sort.Strings(arr1)
	sort.Strings(arr2)

    for i := range arr1 {
        if arr1[i] != arr2[i] {
            return false
        }
    }

    return true
}

func isExistsInErrorQ(currSeed Json_seed, errorQ []Json_seed) bool {
	for _, seed := range errorQ {
		if reflect.DeepEqual(seed, currSeed) {
			return true
		}
	}
	return false
}
