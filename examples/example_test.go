package examples

/*
Base test suite configuration for testing the examples.
*/

import (
	"mongoid"
	"mongoid/log"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	// "fmt"
)

var gLiveTestDatabase bool // when true, a live test database is configured and attached (allowing more tests to occurr)

func TestCriteria(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongoid Examples")
}

var _ = BeforeSuite(func() {
	// LookupEnv and GetEnv are missing a supplied default option...
	fetchEnv := func(envName, defaultVal string) string {
		if val, found := os.LookupEnv(envName); found == true {
			return val
		}
		return defaultVal
	}

	testDbHost := fetchEnv("MONGOID_TEST_DBHOST", "localhost:27017")           // assume localhost for tests
	testDbName := fetchEnv("MONGOID_TEST_DBNAME", "_mongoid_example_test_db_") // made up database name

	gTestMongoidConfig := mongoid.Config{
		Clients: []mongoid.Client{
			{
				Name:     "default",
				Hosts:    []string{testDbHost},
				Database: testDbName,
			},
		},
	}
	mongoid.Configure(&gTestMongoidConfig)

	if !mongoid.Configured() {
		log.Fatal("Tried to configure Mongoid but failed... :(")
	}
	// ConnectionTest() error
	if mongoid.Configuration().DefaultClient.ConnectionTest() != nil {
		log.Fatal("MongoDB failed connection test. Cannot continue.")
	}

	gLiveTestDatabase = true
})

var _ = BeforeEach(func() {
	// wipe the collection (ie, database-cleaner-style)
})
