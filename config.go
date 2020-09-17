package mongoid

import (
	// "mongoid/log"
	"fmt"
	"sync"
)

// Config defines all configuration necessary options, including database Client connection details which may span more than a single server/cluster
// Only one (1) Config  may exist in the system at any one time
type Config struct {
	fmt.Stringer
	ConfigOptions ConfigOptions
	Clients       []Client
	DefaultClient *Client
}

// ConfigOptions define Mongoid specific options.
type ConfigOptions struct {
	//TODO UseUTC bool // Ensure all times are UTC in the app side. (default: false)
	//TODO Logger *log.Logger // Specify a custom log.Logger instance
	//TODO LogLevel int // The log level.
	//TODO PreloadModels bool // Preload all models in development, needed when models use inheritance. (default: false) // not really a thing for Go...?
}

// the Config data holders & mutex
var mongoidConfig *Config         // the global config, used for all future unspecified requests
var mongoidConfigMutex sync.Mutex // the global config mutex, used to synchronize access to all Config access

// Stringer implementation -- value as provided by caller
func (c Config) String() string {
	return fmt.Sprintf(""+
		"ConfigOptions: %v,"+
		" DefaultClient: %v,"+
		" Clients: %v,"+
		"", c.ConfigOptions, c.DefaultClient, c.Clients)
}

// Configuration returns a copy of the current running config
func Configuration() Config {
	mongoidConfigMutex.Lock()
	defer mongoidConfigMutex.Unlock()
	// be careful around this mutex
	var retConfig Config
	retConfig = *mongoidConfig
	return retConfig
}

// Configured returns true if mongoid already has a configuration, false otherwise
func Configured() bool {
	mongoidConfigMutex.Lock()
	defer mongoidConfigMutex.Unlock()
	if mongoidConfig != nil {
		return true
	}
	return false
}

// ClientByName returns the handle to the requested named client
func ClientByName(clientName string) *Client {
	mongoidConfigMutex.Lock()
	defer mongoidConfigMutex.Unlock()
	if mongoidConfig == nil {
		return nil
	}
	for i, client := range mongoidConfig.Clients {
		if client.Name == clientName {
			return &mongoidConfig.Clients[i]
		}
	}
	return nil
}

// DefaultClient returns the handle to the default client
func DefaultClient() *Client {
	mongoidConfigMutex.Lock()
	defer mongoidConfigMutex.Unlock()
	if mongoidConfig == nil {
		return nil
	}
	return mongoidConfig.DefaultClient
}
