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

func CheckCodeSuccess(currCode string) bool {
	success := []string{"Created", "Deleted", "Valid", "Changed", "Content", "Continue"}
	// clientError := []string{"BadRequest", "Unauthorized", "BadOption", "Forbidden", "NotFound", "MethodNotAllowed", "NotAcceptable", "RequestEntityIncomplete", "PreconditionFailed", "RequestEntityTooLarge", "UnsupportedMediaType"}
	// serverError := []string{"TooManyRequests", "InternalServerError", "NotImplemented", "BadGateway", "ServiceUnavailable", "GatewayTimeout", "ProxyingNotSupported"}

	// check if currCode is inside the success codes, return true (no error)
	for _, code := range success {
		if currCode == code {
			return true
		}
	}

	// if not inside success codes, return false (got error)
	return false
}
