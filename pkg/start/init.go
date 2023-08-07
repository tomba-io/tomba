package start

import (
	"fmt"
	"os"

	"github.com/tomba-io/go/tomba"
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
	Key    string
	Secret string
	Target string // Can pass email, Domain, URL, Linkedin URL.
	Output string
	Port   int
	JSON   bool
	YAML   bool
	Color  bool
	Pretty bool
	Use    string
	Search
	Finder
}

type Search struct {
	Page       int
	Limit      int
	Department string
}

type Finder struct {
	FullName  string
	FirstName string
	LastName  string
}

// New parameters
// Initiate cli parameters
func New(conn Conn) *Conn {
	// Read the config file
	conf, err := config.ReadConfigFile()
	if conn.YAML {
		conn.Parameters.JSON = false
	}
	if err != nil {
		fmt.Println("Error reading config file:", err)
	}
	if conf.Key != "" || conf.Secret != "" {
		tomba := tomba.New(conf.Key, conf.Secret)
		conn.Parameters.Key = conf.Key
		conn.Parameters.Secret = conf.Secret
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
