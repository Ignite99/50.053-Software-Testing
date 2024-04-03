package fuzzer

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())

	randomString := make([]byte, length)
	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}
	return string(randomString)
}

func randomizeValue(data map[string]interface{}, key string) {
	rand.Seed(time.Now().UnixNano())

	value := data[key]

	switch v := value.(type) {
	case int:
		data[key] = rand.Int()
	case float64:
		data[key] = rand.Float64()
	case string:
		data[key] = generateRandomString(len(v))
	case bool:
		data[key] = rand.Intn(2) == 1
	case nil:
		data[key] = rand.Int()
	case []interface{}:
		// TODO: Find ways to fuzz arrays
		data[key] = rand.Int()
	default:
		// Unsupported type
		fmt.Println("value type is not covered")
	}
}

func changeValueType(data map[string]interface{}, key string) {
	if value, ok := data[key]; ok {
		switch v := value.(type) {
		case uint, int:
			rand := rand.Intn(2)
			if rand == 0 {
				data[key] = float64(v.(int))
			} else {
				data[key] = strconv.Itoa(v.(int))
			}
		case float64:
			data[key] = int(v)
		case string:
			if intVal, err := strconv.Atoi(v); err == nil {
				data[key] = intVal
			} else if v == "true" || v == "false" {
				data[key] = v == "true"
			} else {
				data[key] = generateRandomString(len(v))
			}
		case bool:
			if v {
				data[key] = "true"
			} else {
				data[key] = "false"
			}
		default:
			fmt.Println("No valid type found for this value")
		}
	}
}

func MutateRequests(requestType string, data map[string]interface{}) map[string]interface{} {
	rand.Seed(time.Now().UnixNano())

	if requestType == "POST" {
		randomMutateType := rand.Intn(4) + 1

		switch randomMutateType {
		// 1. Add new JSON field
		case 1:
			newKey := generateRandomString(5)
			randomJsonField := rand.Intn(4) + 1 // Random number between 1 and 4

			var newVal interface{}
			switch randomJsonField {
			case 1:
				newVal = generateRandomString(5)
			case 2:
				newVal = rand.Intn(100)
			case 3:
				newVal = rand.Float64() * 100
			case 4:
				newVal = rand.Intn(2) == 1
			}

			data[newKey] = newVal

		// 2. Remove existing JSON field
		case 2:
			keys := make([]string, 0, len(data))
			for key := range data {
				keys = append(keys, key)
			}

			if len(keys) > 0 {
				index := rand.Intn(len(keys))
				delete(data, keys[index])
			}

		// 3. Keep all fields as it is but change value type
		case 3:
			keys := make([]string, 0, len(data))
			for key := range data {
				keys = append(keys, key)
			}

			if len(keys) > 0 {
				index := rand.Intn(len(keys))
				changeValueType(data, keys[index])
			}

		// 4. Randomize current JSON value
		case 4:
			keys := make([]string, 0, len(data))
			for key := range data {
				keys = append(keys, key)
			}

			if len(keys) > 0 {
				index := rand.Intn(len(keys))
				randomizeValue(data, keys[index])
			}
		}
	} else if requestType == "GET" {
		// TODO: handle GET request JSON mutations
		fmt.Println("HRMMMMMMMMM GET REQUETS DED HRMMM")
	}

	return data
}
