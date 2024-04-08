package coap

// check if current seed contains the same attributes as those in inputQ
func CheckIsInteresting(currSeed Seed, inputQ []Seed) bool {
	if len(inputQ) == 0 {
		return true
	}
	for _, seed := range inputQ {
		if currSeed.IC.Path == seed.IC.Path &&
			currSeed.IC.Method == seed.IC.Method &&
			currSeed.IC.ContentType == seed.IC.ContentType &&
			currSeed.OC.ContentType == seed.OC.ContentType &&
			currSeed.OC.StatusCode == seed.OC.StatusCode &&
			currSeed.OC.MessageType == seed.OC.MessageType {
			return false
		}
	}
	return true
}
