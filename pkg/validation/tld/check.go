package ltd

// IsValidTLD check if is valid Top-Level Domain
func IsValidTLD(request string) bool {
	for _, validTLD := range availableTLDs() {
		if validTLD == request {
			return true
		}
	}
	return false
}
