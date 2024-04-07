package Interesting

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
)

type OutputCriteria struct {
	ContentType  string
	StatusCode   int
	ResponseBody string
}

type InputCriteria struct {
	Path        string
	Method      string
	ContentType []string
}

type Json_seed struct {
	Data          map[string]interface{}
	Key_to_mutate string
	Energy        int
	OC            OutputCriteria
	IC            InputCriteria
}

func RequestParser(path string, method string, contentType []string) InputCriteria {
	ic := InputCriteria{Path: path, Method: method, ContentType: contentType}
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

func CheckIsInteresting(currSeed Json_seed, prevSeed Json_seed) bool {
	// If the current seed has the same oc and ic as the prev seed's oc and ic --> not interesting
	
	currIc := currSeed.IC
	prevIc := prevSeed.IC

	currOc := currSeed.OC
	prevOc := prevSeed.OC

	boolIcEqual := currIc.Path == prevIc.Path &&
	currIc.Method == prevIc.Method &&
	isContentTypeSame(currIc.ContentType, prevIc.ContentType)

	boolOcEqual := currOc.ContentType == prevOc.ContentType &&
		currOc.StatusCode == prevOc.StatusCode &&
		currOc.ResponseBody == prevOc.ResponseBody

	return !(boolIcEqual && boolOcEqual)
}

func isContentTypeSame(currIcContentType []string, prevIcContentType []string) bool {
	
	if len(currIcContentType) != len(prevIcContentType) {
		return false
    }

	sort.Strings(currIcContentType)
	sort.Strings(prevIcContentType)

    for i := range currIcContentType {
        if currIcContentType[i] != prevIcContentType[i] {
            return false
        }
    }

    return true
}
