package mongoid_test

import (
	"fmt"
	"log"
	"mongoid"
	"os"

	. "github.com/onsi/ginkgo"
)

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
	testDbName := fetchEnv("MONGOID_TEST_DBNAME", "_mongoid_test_db_") // made up database name

	if testDb == "1" {
		fmt.Println("MONGOID TESTS REQUIRING MONGOID_TEST_DB ACTIVATED")
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
	}
})

var _ = BeforeEach(func() {
	// wipe the collection (ie, database-cleaner)
})
