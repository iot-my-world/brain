package main

import (
	"flag"
	"github.com/iot-my-world/brain/internal/log"
	"github.com/iot-my-world/brain/internal/migration"
	"github.com/iot-my-world/brain/internal/migration/v1_v2/zx303"
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

	// run migrations
	if err := zx303.Migrate(brainDb); err != nil {
		log.Fatal(err.Error())
	}
}
