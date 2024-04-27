package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
		"GET":  true,
		"POST": true,
	}
)

func main() {
	var projectType string
	var inputFilePath string
	var energy string
	var url string
	var requestType string
	var energyInt int

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
		energy = os.Args[5]

		if !validProjectTypes[projectType] {
			fmt.Println("Invalid project type. Please enter COAP/DJANGO/BLE.")
			return
		}

		if !validRequestTypes[requestType] {
			fmt.Println("Invalid request type. Please enter GET/POST.")
			return
		}

		inputFile, err := os.Open(os.Args[4])
		if err != nil {
			fmt.Println("Error: Could not open file", os.Args[4], err)
			return
		}
		defer inputFile.Close()

		energyInt, err = strconv.Atoi(energy)
		if err != nil {
			fmt.Println("Error: Cannot convert ", os.Args[5], " into integer.")
			return
		}
		if energyInt < 1 {
			fmt.Println("Error: Energy cannot be negative.")
			return
		}

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

		fmt.Println("\nWhat is your request type? [Options: GET, POST]")
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
			fmt.Println("\nWhat is your energy assignment?")
			_, err = fmt.Scanln(&energy)
			energyInt, err = strconv.Atoi(energy)
			if err != nil {
				fmt.Println("Error: Cannot convert into integer.")
			} else if energyInt < 1 {
				fmt.Println("Error: Energy cannot be negative.")
			} else {
				break
			}
		}
	}

	if projectType == "COAP" {
		fmt.Println("[COAP] Fuzzer has initiated call to COAP!")
		ip_addr := strings.Split(url, ":")[0]
		port := strings.Split(url, ":")[1]
		port_num, err := strconv.Atoi(port)
		if err != nil {
			log.Fatalf("Error converting port to integer: %v", err)
		}
		coap.CoAPTestDriver(ip_addr, port_num, inputFilePath, energyInt)
		// CoAP_Handler()
	} else if projectType == "DJANGO" {
		fmt.Println("[DJANGO] Fuzzer has initiated call to DJANGO Web Application!")
		Django.Django_Test_Driver(energyInt, url, requestType, inputFilePath, "./fuzzing_responses/response.txt")

	} else if projectType == "BLE" {
		fmt.Println("[BLE] Fuzzer has initiated call to BLE_Zephyr!")
		// BLE_Zephyr_Handler()

	} else {
		fmt.Println("Project type mutated. Project type: " + projectType + ". Check code now!")
		return
	}

}
