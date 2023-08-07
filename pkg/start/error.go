package start

import "errors"

var (
	ErrErrInvalidApiKey     = errors.New("Please enter a valid KEY.")
	ErrErrInvalidApiSecret  = errors.New("Please enter a valid SECRET.")
	ErrErrInvalidLogin      = errors.New("Invalid KEY or SECRET.")
	ErrErrInvalidNoLogin    = errors.New("Please Sign in to your account, not logged in.")
	ErrArgumentEmail        = errors.New("Please enter a email, for example 'name@company.com'.")
	ErrArgumentsDomain      = errors.New("Please enter a domain name, for example 'tomba.io'.")
	ErrArgumentsDomainLimit = errors.New("The limit value is not a valid enum value of (10,20,50):")
	ErrArgumentsURL         = errors.New("Please enter a valid article URL")
	ErrArgumentsLinkedinURL = errors.New("Please enter a valid linkedin URL")
	ErrArgumentsFinder      = errors.New("Please enter the full name of the person you'd like to find the email address.")
)
