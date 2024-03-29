package zx303

import (
	"errors"
	"github.com/iot-my-world/brain/internal/migration/v1_v2/zx303/v1"
	"github.com/iot-my-world/brain/internal/migration/v1_v2/zx303/v2"
	"gopkg.in/mgo.v2"
)

const zx303DeviceCollection = "zx303Device"

func Migrate(database *mgo.Database) error {

	// get zx303 device collection
	collection := database.C(zx303DeviceCollection)

	// retrieve old devices
	oldZX303s := make([]v1.ZX303, 0)
	if err := collection.Find(nil).All(&oldZX303s); err != nil {
		return errors.New("error finding all zx303s " + err.Error())
	}
	if oldZX303s == nil {
		oldZX303s = make([]v1.ZX303, 0)
	}

	// clear the collection
	if _, err := collection.RemoveAll(nil); err != nil {
		return errors.New("failed to remove all old zx303 devices: " + err.Error())
	}

	// migrate devices
	newZX303s := make([]interface{}, 0)
	for _, oldZX303 := range oldZX303s {
		newZX303s = append(newZX303s,
			v2.ZX303{
				Id:                     oldZX303.Id,
				IMEI:                   oldZX303.IMEI,
				SimCountryCode:         oldZX303.SimCountryCode,
				SimNumber:              oldZX303.SimNumber,
				OwnerPartyType:         oldZX303.OwnerPartyType,
				OwnerId:                oldZX303.OwnerId,
				AssignedPartyType:      oldZX303.AssignedPartyType,
				AssignedId:             oldZX303.AssignedId,
				LoggedIn:               false,
				LogInTimestamp:         0,
				LogOutTimestamp:        0,
				LastHeartbeatTimestamp: 0,
			})
	}
	insertBulkOperation := collection.Bulk()
	insertBulkOperation.Insert(newZX303s...)
	if _, err := insertBulkOperation.Run(); err != nil {
		return errors.New("error inserting bulk new zx303 devices: " + err.Error())
	}

	return nil
}
