package main

import (
	"fmt"
	"os"
	"strings"

	coap "github.com/50.053-Software-Testing/CoAP"
	"github.com/50.053-Software-Testing/Django"
)

var (
	validProjectTypes = map[string]bool{
		"COAP":   true,
		"DJANGO": true,
		"BLE":    true,
	}
	validRequestTypes = map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"HEAD":   true,
		"DELETE": true,
	}
)

func main() {
	var projectType string
	var inputFilePath string
	var outputFilePath string
	var url string
	var requestType string

	if len(os.Args) != 6 && len(os.Args) != 1 {
		fmt.Println("Error: Invalid arguments provided.")
		fmt.Println("The expected format is:")
		fmt.Println(os.Args[0], "<project type> <url> <request_type> <input_file_path> <output_file_path>")
		return
	}

	if len(os.Args) == 6 {
		projectType = os.Args[1]
		url = os.Args[2]
		requestType = strings.ToUpper(os.Args[3])
		inputFilePath = os.Args[4]
		outputFilePath = os.Args[5]

		if !validProjectTypes[projectType] {
			fmt.Println("Invalid project type. Please enter COAP/DJANGO/BLE.")
			return
		}

		if !validRequestTypes[requestType] {
			fmt.Println("Invalid request type. Please enter GET/POST/PUT/HEAD/DELETE.")
			return
		}

		inputFile, err := os.Open(os.Args[4])
		if err != nil {
			fmt.Println("Error: Could not open file", os.Args[4], err)
			return
		}
		defer inputFile.Close()

		outputFile, err := os.Open(os.Args[5])
		if err != nil {
			fmt.Println("Error: Could not open file", os.Args[5], err)
			return
		}
		defer outputFile.Close()
	} else {
		fmt.Println("What fuzzing target are you testing? [Options: COAP, DJANGO, BLE]")
		var err error
		for {
			_, err = fmt.Scanln(&projectType)
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			projectType = strings.ToUpper(projectType)
			if validProjectTypes[projectType] {
				break
			}
			fmt.Println("Invalid project type. Please enter COAP/DJANGO/BLE.")
		}

		fmt.Println("\nWhat URL are you testing on?")
		_, err = fmt.Scanln(&url)
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		fmt.Println("\nWhat is your request type? [Options: GET, POST, PUT, HEAD, DELETE]")
		for {
			_, err = fmt.Scanln(&requestType)
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			requestType = strings.ToUpper(requestType)
			if validRequestTypes[requestType] {
				break
			}
			fmt.Println("Invalid request type. Please enter GET/POST/PUT/HEAD/DELETE.")
		}

		for {
			fmt.Println("\nWhat is your seed input's file path?")
			_, err = fmt.Scanln(&inputFilePath)
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			inputFile, err := os.Open(inputFilePath)
			if err == nil {
				defer inputFile.Close()
				break
			}
		}

		for {
			fmt.Println("\nWhat is your output file path?")
			_, err = fmt.Scanln(&outputFilePath)
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			outputFile, err := os.Open(outputFilePath)
			if err == nil {
				defer outputFile.Close()
				break
			}
		}
	}

	if projectType == "COAP" {
		fmt.Println("[COAP] Fuzzer has initiated call to COAP!")
		coap.CoAPTestDriver()
		// CoAP_Handler()
	} else if projectType == "DJANGO" {
		fmt.Println("[DJANGO] Fuzzer has initiated call to DJANGO Web Application!")
		Django.Django_Test_Driver(1, url, requestType, inputFilePath, outputFilePath)

	} else if projectType == "BLE" {
		fmt.Println("[BLE] Fuzzer has initiated call to BLE_Zephyr!")
		// BLE_Zephyr_Handler()

	} else {
		fmt.Println("Project type mutated. Project type: " + projectType + ". Check code now!")
		return
	}

}
