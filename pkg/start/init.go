package start

import (
	"fmt"

	"github.com/tomba-io/tomba/pkg/config"
	"github.com/tomba-io/go/tomba"
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
	Target string // Can pass email, Domain, URL, Linkedin URL or TXT file for bulk.
	Output string
	JSON   bool
	YAML   bool
	Color  bool
	Pretty bool
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

	tomba := tomba.New(conn.Key, conn.Secret)
	return &Conn{
		Parameters: conn.Parameters,
		Tomba:      tomba,
	}

}
