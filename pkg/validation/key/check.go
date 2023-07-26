package key

// IsValidAPI Check if Tomba.io kEY or SECRET is valid
func IsValidAPI(credential string) bool {
	return len(credential) >= 39
}
