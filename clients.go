package mongoid

import (
	"context"
	"fmt"
	"mongoid/log"
	"time"

	mongoidError "mongoid/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// Client a client connection performing operations on a given topology
type Client struct {
	fmt.Stringer
	Name     string
	URI      string
	Hosts    []string
	Database string
	AppName  string

	Options            ClientOptions
	mongoClient        *mongo.Client
	mongoClientOptions *options.ClientOptions
}

// ClientOptions define client specific options
type ClientOptions struct {
	//NYI Write              interface{} // write concern // TODO implement this
	//NYI Read               interface{} // read concern  // TODO implement this
	Username           string        // Username for authentication
	Password           string        // Password for authentication
	AuthMech           string        // Change the default authentication mechanism. Valid options are: "scram", "mongodb_cr", "mongodb_x509", and "plain". (default on 3.0 is :scram, default on 2.4 and 2.6 is :plain)
	AuthSource         string        // The database or source to authenticate the user against. (default: admin)
	Connect            string        // Force the driver to connect in a specific way instead of auto-discovering. Can be one of: "direct", "replica_set", "sharded".
	HeartbeatFrequency time.Duration // Change the default time the server monitors refresh their status via ismaster commands. (default: 10 seconds)
	ConnectTimeout     time.Duration // The time to wait to establish a connection before timing out. (default: 10 seconds)
	MaxPoolSize        uint64        // The maximum number of connections in the connection pool. (default: 5)

	//NYI WaitQueueTimeout   int    // The time to wait, in seconds, in the connection pool for a connection to be checked in before timing out. (default: 1)
	SocketTimeout          time.Duration // The timeout to wait to execute operations on a socket before raising an error. (default: 0 / infinite)
	LocalThreshold         time.Duration // The time for selecting servers for a near read preference. (default: 0.015 seconds)
	ServerSelectionTimeout time.Duration // The timeout for selecting a server for an operation. (default: 30 seconds)

	SSL             bool // Whether to connect to the servers via ssl. (default: false)
	SSLVerify       bool // Whether or not to do peer certification validation. (default: false) // DEVIATION was sslinsecure
	SkipConnectTest bool // If true, skips connection tests; otherwise the client will attempt a server Ping check during init, throwing error if the service is unavailable (default: false) // DEVIATION
}

const constDefaultHeartbeatFrequency time.Duration = 10 * time.Second

// GetHeartbeatFrequency gives the effective heartbeat frequency
func (opts *ClientOptions) GetHeartbeatFrequency() time.Duration {
	if opts == nil || opts.HeartbeatFrequency == 0 {
		log.Trace("default GetHeartbeatFrequency")
		return constDefaultHeartbeatFrequency
	}
	return opts.HeartbeatFrequency
}

const constDefaultConnectionTimeout time.Duration = 10 * time.Second

// GetConnectTimeout gives the effective connection timeout
func (opts *ClientOptions) GetConnectTimeout() time.Duration {
	if opts == nil || opts.ConnectTimeout == 0 {
		return constDefaultConnectionTimeout
	}
	return opts.ConnectTimeout
}

const constDefaultSocketTimeout time.Duration = 0 * time.Second

// GetSocketTimeout gives the effective socket timeout
func (opts *ClientOptions) GetSocketTimeout() time.Duration {
	if opts == nil || opts.SocketTimeout == 0 {
		return constDefaultSocketTimeout
	}
	return opts.SocketTimeout
}

const constDefaultLocalThreshold time.Duration = 15 * time.Millisecond

// GetLocalThreshold gives the effective local threshold
func (opts *ClientOptions) GetLocalThreshold() time.Duration {
	if opts == nil || opts.LocalThreshold == 0 {
		return constDefaultLocalThreshold
	}
	return opts.LocalThreshold
}

const constDefaultServerSelectionTimeout time.Duration = 30 * time.Second

// GetServerSelectionTimeout gives the effective sever selection timeout
func (opts *ClientOptions) GetServerSelectionTimeout() time.Duration {
	if opts == nil || opts.ServerSelectionTimeout == 0 {
		return constDefaultServerSelectionTimeout
	}
	return opts.ServerSelectionTimeout
}

const constDefaultMaxPoolSize uint64 = 5

// GetMaxPoolSize gives the effective max pool size
func (opts *ClientOptions) GetMaxPoolSize() uint64 {
	if opts == nil || opts.MaxPoolSize == 0 {
		return constDefaultMaxPoolSize
	}
	return opts.MaxPoolSize
}

// const constDefaultReadConcern string = "majority"
// const constDefaultWriteConcern string = "majority"
// const constDefaultReadPreference string = "secondaryPreferred"
// const constDefaultDatabaseName string = "incidental"

func (c Client) String() string {
	uninitialized := ""
	if c.mongoClient == nil {
		uninitialized = "[uninitialized] "
	}
	return fmt.Sprintf("{%s"+
		"Name: %s,"+
		" URI: %s,"+
		" Hosts: %s"+
		" Database: %s"+
		" AppName: %s"+
		" Options: %+v"+
		"}", uninitialized, c.Name, c.URI, c.Hosts, c.Database, c.AppName, c.Options)
}

// Connect the Client to the server topology
func (c *Client) Connect() error {
	log.Debug("Connect()")
	// ensure client is disconnected (has no affect if was never connected)
	if err := c.Disconnect(); err != nil {
		return err
	}

	c.mongoClientOptions = buildMongoClientOptions(*c)

	client, err := mongo.NewClient(c.mongoClientOptions)
	ctx, cancel := context.WithTimeout(context.Background(), c.Options.GetConnectTimeout())
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return err
	}
	c.mongoClient = client

	if false == c.Options.SkipConnectTest {
		err := c.connectionTest(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect the Client from the server topology
func (c *Client) Disconnect() error {
	if c.mongoClient != nil {
		log.Debug("Disconnect()")
		err := c.mongoClient.Disconnect(nil)
		c.mongoClient = nil
		if err != nil {
			return err
		}
	}
	return nil
}

// ConnectionTest (Client.ConnectionTest) performs a basic connectivity check upon a single Client
// returns nil on Success, a specific error condition on failure
func (c *Client) ConnectionTest() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.Options.GetConnectTimeout())
	defer cancel()
	return c.connectionTest(ctx)
}
func (c *Client) connectionTest(ctx context.Context) error {
	log.Debug("ConnectionTest()")
	err := c.mongoClient.Ping(ctx, nil)
	if err != nil {
		if err == context.DeadlineExceeded {
			log.Error("ConnectionTest() FAILED")
			return &mongoidError.OperationTimedOut{
				MethodName: "connectionTest",
				Reason:     "Deadline Exceeded",
			}
		}
		return err
	}
	return nil
}

// MongoDriverClient returns the underlying Mongo Driver Client (go.mongodb.org/mongo-driver/mongo) for the mongoid.Client it was called upon
func (c *Client) MongoDriverClient() *mongo.Client {
	if c.mongoClient != nil {
		return c.mongoClient
	}
	return nil
}

func buildMongoClientOptions(mongoidClient Client) *options.ClientOptions {
	log.Trace("buildMongoClientOptions() ", mongoidClient)
	clientOpts := options.Client()
	if mongoidClient.URI != "" {
		clientOpts = clientOpts.ApplyURI(mongoidClient.URI)
	}
	clientOpts.Hosts = mongoidClient.Hosts
	// clientOpts.Database = mongoidClient.Database // Database used by mongoid, not by mongo-driver
	if mongoidClient.AppName != "" {
		clientOpts.SetAppName(mongoidClient.AppName)
	}

	if mongoidClient.Options.AuthMech != "" || mongoidClient.Options.Username != "" || mongoidClient.Options.Password != "" {
		log.Trace("building new auth")
		auth := options.Credential{}
		if clientOpts.Auth != nil {
			log.Tracef("has existing auth: %+v\n", clientOpts.Auth)
			auth = *clientOpts.Auth
		}
		auth.AuthSource = mongoidClient.Options.AuthSource
		auth.AuthMechanism = mongoidClient.Options.AuthMech
		auth.Username = mongoidClient.Options.Username
		auth.Password = mongoidClient.Options.Password
		if mongoidClient.Options.Password != "" {
			auth.PasswordSet = true
		}
		log.Tracef("new built auth: %+v\n", auth)
		clientOpts.SetAuth(auth)
	}

	if mongoidClient.Options.Connect == "direct" {
		clientOpts.SetDirect(true)
	}

	clientOpts.SetConnectTimeout(mongoidClient.Options.GetConnectTimeout())
	clientOpts.SetHeartbeatInterval(mongoidClient.Options.GetHeartbeatFrequency())
	clientOpts.SetServerSelectionTimeout(mongoidClient.Options.GetServerSelectionTimeout())
	clientOpts.SetLocalThreshold(mongoidClient.Options.GetLocalThreshold())
	if mongoidClient.Options.GetSocketTimeout() > 0 {
		clientOpts.SetSocketTimeout(mongoidClient.Options.GetSocketTimeout())
	}
	clientOpts.SetMaxPoolSize(mongoidClient.Options.GetMaxPoolSize())

	log.Tracef("%+v\n", clientOpts)

	return clientOpts
}

// ApplyURI parses the provided connection string and sets the values and options accordingly.
func (c *Client) ApplyURI(uri string) (*Client, error) {
	log.Debug("Client.ApplyURI()")

	connString, err := connstring.Parse(uri)
	if err != nil {
		return nil, err
	}
	c.URI = connString.Original
	c.Hosts = connString.Hosts
	c.Database = connString.Database
	c.AppName = connString.AppName

	c.Options.AuthSource = connString.AuthSource
	c.Options.Username = connString.Username
	c.Options.Password = connString.Password

	if len(connString.Options["connect"]) > 0 {
		c.Options.Connect = connString.Options["connect"][0]
	}
	if connString.HeartbeatIntervalSet {
		c.Options.HeartbeatFrequency = connString.HeartbeatInterval
	}
	if connString.ConnectTimeoutSet {
		c.Options.ConnectTimeout = connString.ConnectTimeout
	}
	if connString.MaxPoolSizeSet {
		c.Options.MaxPoolSize = connString.MaxPoolSize
	}
	if connString.SocketTimeoutSet {
		c.Options.SocketTimeout = connString.SocketTimeout
	}
	if connString.LocalThresholdSet {
		c.Options.LocalThreshold = connString.LocalThreshold
	}
	if connString.ServerSelectionTimeoutSet {
		c.Options.ServerSelectionTimeout = connString.ServerSelectionTimeout
	}
	if connString.SSLSet {
		c.Options.SSL = connString.SSL
	}
	if connString.SSLInsecureSet {
		c.Options.SSLVerify = !connString.SSLInsecure
	}

	// log.Printf("%+v\n", connString)
	// log.Printf("%+v\n", c)
	// log.Printf("%+v\n", c.Options)

	// clientOpts := options.Client()
	// clientOpts = clientOpts.ApplyURI(uri)

	return c, nil
}

func (c *Client) getMongoCollectionHandle(databaseName, collectionName string) *mongo.Collection {
	return c.MongoDriverClient().Database(databaseName).Collection(collectionName)
}
