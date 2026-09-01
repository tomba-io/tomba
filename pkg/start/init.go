package start

import (
	"fmt"
	"os"

	"github.com/tomba-io/go/tomba"

	"github.com/tomba-io/tomba/pkg/auth"
	"github.com/tomba-io/tomba/pkg/config"
	"github.com/tomba-io/tomba/pkg/util"
	_key "github.com/tomba-io/tomba/pkg/validation/key"
)

// Conn
type Conn struct {
	Parameters
	*tomba.Tomba
}

// Parameters configuration for the cli
type Parameters struct {
	Key     string
	Secret  string
	Target  string // Can pass email, Domain, URL, Linkedin URL.
	Output  string
	Port    int
	JSON    bool
	YAML    bool
	CSV     bool
	Color   bool
	Pretty  bool
	Verbose bool
	Use     string
	Search
	Finder
}

type Search struct {
	Page       int
	Limit      int
	Department string
}

type Finder struct {
	FullName     string
	FirstName    string
	LastName     string
	EnrichMobile bool
}

// New parameters
// Initiate cli parameters
func New(conn Conn) *Conn {
	// Read the config file
	conf, err := config.ReadConfigFile()
	if conn.YAML {
		conn.JSON = false
	}
	if err != nil {
		fmt.Println("Error reading config file:", err)
	}

	// Try OAuth authentication first
	if conf != nil && conf.AuthMethod == "oauth" && conf.AccessToken != "" {
		accessToken, err := auth.EnsureValidToken(conf)
		if err != nil {
			if conn.Use != "login" {
				fmt.Println(util.WarningIcon(), util.Yellow("OAuth token expired. Please run 'tomba login' to re-authenticate."))
				os.Exit(0)
			}
		} else {
			// Fetch API key+secret via OAuth token, then use the SDK normally.
			// The /me endpoint returns the user's API key and secret when
			// authenticated, so we exchange the OAuth token for credentials
			// that the Go SDK can use.
			apiKey, apiSecret, fetchErr := auth.FetchAPICredentials(accessToken)
			if fetchErr != nil {
				fmt.Println(util.WarningIcon(), util.Yellow("Could not retrieve API credentials: "+fetchErr.Error()))
				os.Exit(0)
			}
			t := tomba.New(apiKey, apiSecret)
			conn.Key = apiKey
			conn.Secret = apiSecret
			return &Conn{
				Parameters: conn.Parameters,
				Tomba:      t,
			}
		}
	}

	// Fall back to API key authentication
	if conf != nil && (conf.Key != "" || conf.Secret != "") {
		tomba := tomba.New(conf.Key, conf.Secret)
		conn.Key = conf.Key
		conn.Secret = conf.Secret
		return &Conn{
			Parameters: conn.Parameters,
			Tomba:      tomba,
		}
	}
	if conn.Use != "login" {
		if conn.Key == "" || conn.Secret == "" {
			fmt.Println(util.WarningIcon(), util.Yellow(ErrErrInvalidNoLogin.Error()))
			os.Exit(0)
		}
		if !_key.IsValidAPI(conn.Key) && !_key.IsValidAPI(conn.Secret) {
			fmt.Println(util.WarningIcon(), util.Yellow(ErrErrInvalidLogin.Error()))
			os.Exit(0)
		}
	}
	tomba := tomba.New(conn.Key, conn.Secret)
	return &Conn{
		Parameters: conn.Parameters,
		Tomba:      tomba,
	}
}
