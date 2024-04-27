package coap

import (
	"math/rand"
	"time"
)

const symbols = "!@#$%^&*()_+-=[]{}|;:',.<>?~"

func mutate_random(payload string) string {
	// call a random mutation function
	rand.Seed(time.Now().UnixNano())

	mutationFunctions := []func(string) string{
		mutate_add_byte,
		mutate_delete_byte,
		mutate_add_bytes,
		mutate_delete_bytes,
		mutate_replace_bytes,
		mutate_flip_bytes,
		mutate_reverse_bytes,
		mutate_add_symbols,
		mutate_empty_string,
	}

	randomIndex := rand.Intn(len(mutationFunctions))
	mutatedPayload := mutationFunctions[randomIndex](payload)

	return mutatedPayload
}

func mutate_add_byte(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// add a random byte at a random position in the payloadBytes
	rand.Seed(time.Now().UnixNano())
	randomByte := byte(rand.Intn(256))
	randomPosition := rand.Intn(len(payloadBytes) + 1)
	payloadBytes = append(payloadBytes[:randomPosition], append([]byte{randomByte}, payloadBytes[randomPosition:]...)...)

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_add_bytes(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	rand.Seed(time.Now().UnixNano())

	randomNum := rand.Intn(100)

	// add random number of random bytes
	for i := 0; i < randomNum; i++ {
		randomByte := byte(rand.Intn(256))
		randomPosition := rand.Intn(len(payloadBytes) + 1)
		payloadBytes = append(payloadBytes[:randomPosition], append([]byte{randomByte}, payloadBytes[randomPosition:]...)...)
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_delete_byte(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// ensure there's at least one byte in the payload
	if len(payloadBytes) == 0 {
		return payload
	}

	// select a random position to delete
	rand.Seed(time.Now().UnixNano())
	randomPosition := rand.Intn(len(payloadBytes))

	// delete the byte at the random position
	payloadBytes = append(payloadBytes[:randomPosition], payloadBytes[randomPosition+1:]...)

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_delete_bytes(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// ensure there's at least one byte in the payload
	if len(payloadBytes) == 0 {
		return payload
	}

	rand.Seed(time.Now().UnixNano())

	// upper bound is length of payloadBytes
	randomNum := rand.Intn(len(payloadBytes))

	// delete random number of random bytes
	for i := 0; i < randomNum; i++ {
		randomPosition := rand.Intn(len(payloadBytes))
		payloadBytes = append(payloadBytes[:randomPosition], payloadBytes[randomPosition+1:]...)
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_replace_bytes(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// ensure there's at least one byte in the payload
	if len(payloadBytes) == 0 {
		return payload
	}

	rand.Seed(time.Now().UnixNano())
	numToReplace := rand.Intn(len(payloadBytes))

	// replace the specified number of random bytes
	for i := 0; i < numToReplace; i++ {
		randomPosition := rand.Intn(len(payloadBytes))
		randomByte := byte(rand.Intn(256))
		payloadBytes[randomPosition] = randomByte
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_flip_bytes(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// ensure there's at least one byte in the payload
	if len(payloadBytes) == 0 {
		return payload
	}

	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(len(payloadBytes))

	// flip the specified number of random bytes
	for i := 0; i < randomNum; i++ {
		randomPosition := rand.Intn(len(payloadBytes))
		payloadBytes[randomPosition] = ^payloadBytes[randomPosition]
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_reverse_bytes(payload string) string {
	// convert the payload to a byte slice
	payloadBytes := []byte(payload)

	// ensure there's at least one byte in the payload
	if len(payloadBytes) == 0 {
		return payload
	}

	// reverse the order of bytes in the payload
	for i, j := 0, len(payloadBytes)-1; i < j; i, j = i+1, j-1 {
		payloadBytes[i], payloadBytes[j] = payloadBytes[j], payloadBytes[i]
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_add_symbols(payload string) string {
	payloadBytes := []byte(payload)

	rand.Seed(time.Now().UnixNano())
	numToAdd := rand.Intn(len(payloadBytes))

	// add symbols to the payload
	for i := 0; i < numToAdd; i++ {
		randomSymbol := symbols[rand.Intn(len(symbols))]
		payloadBytes = append(payloadBytes, randomSymbol)
	}

	// convert the byte slice back to a string
	mutatedPayload := string(payloadBytes)

	return mutatedPayload
}

func mutate_empty_string(payload string) string {
	rand.Seed(time.Now().UnixNano())

	emptyList := []string{
		"",
		"{}",
		" ",
		"null",
	}

	randomIndex := rand.Intn(len(emptyList))
	return emptyList[randomIndex]
}
