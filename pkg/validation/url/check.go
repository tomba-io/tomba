package url

import (
	"regexp"

	"github.com/asaskevich/govalidator"
)

// IsValidURL Check if url is valid
func IsValidURL(url string) bool {
	return govalidator.IsURL(url)
}

// IsValidLinkedInProfile Check if LinkedIn profile is valid
func IsValidLinkedInProfile(url string) bool {
	if govalidator.IsURL(url) {
		pattern := `(?:https?:)?\/\/(?:[\w]+\.)?linkedin\.com\/in\/(?P<permalink>[\w\-\_À-ÿ%]+)\/?`
		regex := regexp.MustCompile(pattern)
		match := regex.FindString(url)
		if match != "" {
			return true
		}
	}

	return false
}
