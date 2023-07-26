package domain

import (
	"strconv"
	"strings"

	_tld "github.com/tomba-io/tomba/pkg/validation/tld"
)

// isValidDomain Check if domain looks valid
func IsValidDomain(domain string) bool {
	split := strings.Split(domain, ".")
	if len(split) < 2 {
		return false
	}

	ending := split[len(split)-1]
	if len(ending) < 2 {
		return false
	}
	// check if TLD name actually exists and is not some image ending
	if !_tld.IsValidTLD(ending) {
		return false
	}

	if _, err := strconv.Atoi(ending); err == nil {
		return false
	}

	return true
}
