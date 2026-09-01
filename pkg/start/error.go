package start

import "errors"

var (
	ErrErrInvalidApiKey     = errors.New("please enter a valid KEY")
	ErrErrInvalidApiSecret  = errors.New("please enter a valid SECRET")
	ErrErrInvalidLogin      = errors.New("invalid KEY or SECRET")
	ErrErrInvalidNoLogin    = errors.New("please sign in to your account, not logged in")
	ErrArgumentEmail        = errors.New("please enter a email, for example 'name@company.com'")
	ErrArgumentsDomain      = errors.New("please enter a domain name, for example 'tomba.io'")
	ErrArgumentsDomainLimit = errors.New("the limit value is not a valid enum value of (10,20,50)")
	ErrArgumentsURL         = errors.New("please enter a valid article URL")
	ErrArgumentsLinkedinURL = errors.New("please enter a valid linkedin URL")
	ErrArgumentsFinder      = errors.New("please enter the full name of the person you'd like to find the email address")
)
