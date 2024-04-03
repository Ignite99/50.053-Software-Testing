package Interesting

import (
	"fmt"
	"io/ioutil"
	"net/http"
)


type outputCriteria struct{
	contentType string
	statusCode int
	responseBody string
}

type inputCriteria struct{
	path string
	method string
	contentType string
}

type json_seed struct {
	data          map[string]interface{}
	key_to_mutate string
	energy        int
	oc			  outputCriteria
	ic 			  inputCriteria
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

func CheckIsInteresting(input_queue []json_seed, oc outputCriteria, ic inputCriteria) bool{
	// TODO - check if the seeds in the queue have the same oc and ic as the current json_seed's oc and ic
	return false;
}