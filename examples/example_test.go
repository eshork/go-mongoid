package examples

/*
Base test suite configuration for testing the examples.
*/

import (
	"fmt"
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

var MONGOID_TEST_DB bool

var _ = BeforeSuite(func() {
	// set up the test database connection, if we should...

	// LookupEnv and GetEnv are missing a supplied default option...
	fetchEnv := func(envName, defaultVal string) string {
		if val, found := os.LookupEnv(envName); found == true {
			return val
		}
		return defaultVal
	}

	testDb := fetchEnv("MONGOID_TEST_DB", "1")                         // overall live db test toggle, default is "on"
	testDbHost := fetchEnv("MONGOID_TEST_DBHOST", "localhost:27017")   // assume localhost for tests
	testDbName := fetchEnv("MONGOID_TEST_DBNAME", "_mongoid_test_example_db_") // made up database name

	if testDb == "1" {
		fmt.Println("MONGOID TESTS REQUIRING MONGOID_TEST_DB ACTIVATED")
		MONGOID_TEST_DB = true
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
			log.Fatal("MONGOID_TEST_DB=1 but MongoDB failed connection test. Cannot continue.")
		}
	} else {
		fmt.Println(".. skipping tests requiring online database ..")
	}
})

var _ = BeforeEach(func() {
	// wipe the collection (ie, database-cleaner-style)
})

// helper/wrapper function - tests that require a database can enclose their contents with this to prevent execution (and thus failure) when a db is not present
func OnlineDatabaseOnly(f func()) {
	// execute the given block if a database is currently usable, otherwise skip it (so it always passes) and try to announce it was skipped
	if MONGOID_TEST_DB {
		f()
	} else {
		By("skipping test requiring online database (PASS)")
		// fmt.Println("skipping test requiring online database (PASS)") // very loud announcement
	}
}
