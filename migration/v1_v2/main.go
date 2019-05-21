package v1_v2

import (
	"flag"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/migration"
)

func main() {
	mongoNodes := flag.String("mongoNodes", "localhost:27017", "the nodes in the db cluster")
	mongoUser := flag.String("mongoUser", "", "brains mongo db user")
	mongoPassword := flag.String("mongoPassword", "", "passwords for brains mongo db")

	databaseName := "brain"

	// connect to database
	databaseSession, err := migration.GetDatabaseSession(
		*mongoNodes,
		*mongoUser,
		*mongoPassword,
		databaseName,
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer databaseSession.Close()
	brainDb := databaseSession.DB(databaseName)

}
