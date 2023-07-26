package email

import (
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	_tld "github.com/tomba-io/tomba/pkg/validation/tld"
)

// IsValidEmail Check if email is valid
func IsValidEmail(email string) bool {

	split := strings.Split(email, ".")
	ending := split[len(split)-1]

	if len(ending) < 2 {
		return false
	}

	if len(email) > 254 {
		return false
	}
	// check if TLD name actually exists and is not some image ending
	if !_tld.IsValidTLD(ending) {
		return false
	}
	// check with govalidator
	if !govalidator.IsEmail(email) {
		return false
	}

	if _, err := strconv.Atoi(ending); err == nil {
		return false
	}

	return true
}
