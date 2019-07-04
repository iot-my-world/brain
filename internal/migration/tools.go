package migration

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/iot-my-world/brain/internal/log"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

func GetDatabaseSession(mongoNodes, mongoUser, mongoPassword, databaseName string) (*mgo.Session, error) {
	mongoCluster := strings.Split(mongoNodes, ",")
	log.Info(fmt.Sprintf("connecting to mongo @ node addresses: [%s]", strings.Join(mongoCluster, ", ")))

	dialInfo := mgo.DialInfo{
		Addrs:     mongoCluster,
		Username:  mongoUser,
		Password:  mongoPassword,
		Mechanism: "SCRAM-SHA-1",
		Timeout:   10 * time.Second,
		Source:    "admin",
		Database:  databaseName,
	}
	mongoSession, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return nil, errors.New("could not connect to mongo cluster: " + err.Error())
	}
	log.Info("Connected to Mongo!")
	return mongoSession, nil
}
