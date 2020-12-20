package mongoid

import (
	mongoidError "mongoid/errors"
	"mongoid/log"
)

// Configure sets up the given Config to be used by future requests to the module.
// Calling Configure more than once is specifically prohibited to call attention to multiple-initialization scenarios,
// but if that is your intent, you can safely call ReConfigure() for any configureation call (whether first or subsequent)
func Configure(config *Config) {
	log.Debugf("Configure(%+v)\n", config)
	if config == nil {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Configure",
			Reason:     "Given config was nil",
		})
	}
	mongoidConfigMutex.Lock()
	defer mongoidConfigMutex.Unlock()
	if mongoidConfig != nil { // fail badly if there is already an existing Config
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Configure",
			Reason:     "Already configured",
		})
	}
	if err := configure(config); err != nil {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Configure",
			Reason:     err.Error(),
		})
	}
}

// ReConfigure clears the existing Config and reperforms the Configure process with the newly provided Config.
// This is safe to perform as the first-time configure action, if that is desired.
func ReConfigure(config *Config) {
	if config == nil {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Configure",
			Reason:     "Given config was nil",
		})
	}

	mongoidConfigMutex.Lock()
	defer mongoidConfigMutex.Unlock()
	mongoidConfig = nil // blow away any existing Config

	if err := configure(config); err != nil {
		log.Panic(&mongoidError.InvalidOperation{
			MethodName: "Configure",
			Reason:     err.Error(),
		})
	}
}

// The actual configure function -- assumes that it is performing a clean Config build, meaning that all parts of the config are the new expected value.
// The immediate in-module callers should most likely bubble up all returned error conditions when presented.
func configure(config *Config) error {

	// since this is a (forcefully) dereferenced assignment, the given ptr ref obj can change without affecting our stored value
	newConfigObj := *config
	mongoidConfig = &newConfigObj

	// add a default client if none are specified
	if len(mongoidConfig.Clients) == 0 {
		newClient := new(Client)
		newClient.Hosts = []string{"localhost:27017"}
		newClient.Name = "default"
		newClient.Database = "mongoid"
		mongoidConfig.Clients = append(mongoidConfig.Clients, *newClient)
		log.Error("Added a default Mongoid client configuration -- you should specify this information")
	}

	// if needed, determine the 'default' client, either by name or by order
	if mongoidConfig.DefaultClient == nil {
		var firstClient *Client
		var defaultClient *Client
		// any of them have a name like 'default' (first wins) if so that would be it, else just pick the first overall
		for i, client := range mongoidConfig.Clients {
			if firstClient == nil {
				firstClient = &mongoidConfig.Clients[i]
			}
			if defaultClient == nil {
				if client.Name == "default" {
					defaultClient = &mongoidConfig.Clients[i]
				}
			}
		}
		if defaultClient == nil { // seems no 'default' exists, so just elect the first entry
			mongoidConfig.DefaultClient = firstClient
			if firstClient.Name == "" {
				firstClient.Name = "default"
			}
		} else { // a 'default' was found, so use that
			mongoidConfig.DefaultClient = defaultClient
		}
	}

	for i := range mongoidConfig.Clients {
		clientPrt := &mongoidConfig.Clients[i]
		err := clientPrt.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

// InitializeFromYaml sets up a new Config (or group of configs) that may be used by future requests, using a YAML file as the source
// func InitializeFromYaml(yamlFile string) {}
