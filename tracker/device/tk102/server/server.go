package server

import (
	"bufio"
	"fmt"
	"github.com/go-errors/errors"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/search/criterion"
	"gitlab.com/iotTracker/brain/search/criterion/text"
	"gitlab.com/iotTracker/brain/search/identifier/device/tk102"
	"gitlab.com/iotTracker/brain/search/identifier/id"
	"gitlab.com/iotTracker/brain/search/query"
	"gitlab.com/iotTracker/brain/security/claims/login"
	"gitlab.com/iotTracker/brain/tracker/device"
	tk1022 "gitlab.com/iotTracker/brain/tracker/device/tk102"
	tk102RecordHandler "gitlab.com/iotTracker/brain/tracker/device/tk102/recordHandler"
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	"net"
	"strconv"
	"strings"
	"time"
)

// TK102Server is a TK102 Tracking Device tcp/id unix socket server
type TK102Server struct {
	readingRecordHandler readingRecordHandler.RecordHandler
	tk102RecordHandler   tk102RecordHandler.RecordHandler
	systemClaims         *login.Login
	ip                   string
	port                 string
}

// New TK102 Server
func New(
	readingRecordHandler readingRecordHandler.RecordHandler,
	systemClaims *login.Login,
	tk102RecordHandler tk102RecordHandler.RecordHandler,
	ip string,
	port string,
) *TK102Server {

	return &TK102Server{
		tk102RecordHandler:   tk102RecordHandler,
		readingRecordHandler: readingRecordHandler,
		ip:                   ip,
		port:                 port,
		systemClaims:         systemClaims,
	}
}

func convertData(raw string) (*reading.Reading, *tk102.Identifier, error) {
	// e.g. raw 027045053055BR00190223A2609.0026S02808.1030E000.21226590.000000000000L00000000
	// will be converted to: -26.150043 28.135050
	newReading := reading.Reading{}
	tk102Identifier := tk102.Identifier{}

	// minimum length of 70?
	if len(raw) < 70 {
		return nil, nil, errors.New("raw data too short")
	}
	newReading.Raw = raw

	newReading.TimeStamp = time.Now().UTC().Unix()

	// check for BR which is used to split out Manufacturer ID
	if strings.Contains(raw, "BR") {
		tk102Identifier.ManufacturerId = raw[2:strings.Index(raw, "BR")]
	} else {
		return nil, nil, errors.New("could not find BR separator in raw data")
	}

	// check for A which separates date and co-ordinates
	if strings.Contains(raw, "A") {
		// confirm N only appears once
		if strings.Count(raw, "A") > 1 {
			return nil, nil, errors.New("more than 1 A in raw data")
		}
	}

	// check for N or S
	nOrS := ""
	if strings.Contains(raw, "N") {
		// confirm N only appears once
		if strings.Count(raw, "N") > 1 {
			return nil, nil, errors.New("more than 1 N in raw data")
		}

		// perform processing for N
		nOrS = "N"
		north := raw[strings.Index(raw, "A")+1 : strings.Index(raw, "N")]
		minutes := north[strings.Index(north, ".")-2:]
		degrees := north[:strings.Index(north, ".")-2]
		floatDegrees, err := strconv.ParseFloat(degrees, 32)
		if err != nil {
			return nil, nil, errors.New("error converting degrees string to float: " + err.Error())
		}
		floatMinutes, err := strconv.ParseFloat(minutes, 32)
		if err != nil {
			return nil, nil, errors.New("error converting string minutes to float: " + err.Error())
		}
		newReading.Latitude = float32(floatDegrees + (floatMinutes / 60))
	} else if strings.Contains(raw, "S") {
		// confirm S only appears once
		if strings.Count(raw, "S") > 1 {
			return nil, nil, errors.New("more than 1 S in raw data")
		}
		nOrS = "S"
		// perform processing for S
		south := raw[strings.Index(raw, "A")+1 : strings.Index(raw, "S")]
		minutes := south[strings.Index(south, ".")-2:]
		degrees := south[:strings.Index(south, ".")-2]
		floatDegrees, err := strconv.ParseFloat(degrees, 32)
		if err != nil {
			return nil, nil, errors.New("error converting degrees string to float: " + err.Error())
		}
		floatMinutes, err := strconv.ParseFloat(minutes, 32)
		if err != nil {
			return nil, nil, errors.New("error converting string minutes to float: " + err.Error())
		}
		newReading.Latitude = -1 * float32(floatDegrees+(floatMinutes/60))
	} else {
		return nil, nil, errors.New("could not find N or S in raw data data")
	}

	// check for E or W
	if strings.Contains(raw, "E") {
		// confirm N only appears once
		if strings.Count(raw, "E") > 1 {
			return nil, nil, errors.New("more than 1 E in raw data")
		}
		// process E
		east := raw[strings.Index(raw, nOrS)+1 : strings.Index(raw, "E")]
		minutes := east[strings.Index(east, ".")-2:]
		degrees := east[:strings.Index(east, ".")-2]
		floatDegrees, err := strconv.ParseFloat(degrees, 32)
		if err != nil {
			return nil, nil, errors.New("error converting degrees string to float: " + err.Error())
		}
		floatMinutes, err := strconv.ParseFloat(minutes, 32)
		if err != nil {
			return nil, nil, errors.New("error converting string minutes to float: " + err.Error())
		}
		newReading.Longitude = float32(floatDegrees + (floatMinutes / 60))
	} else if strings.Contains(raw, "W") {
		// confirm S only appears once
		if strings.Count(raw, "W") > 1 {
			return nil, nil, errors.New("more than 1 W in raw data")
		}
		// process W
		west := raw[strings.Index(raw, nOrS)+1 : strings.Index(raw, "W")]
		minutes := west[strings.Index(west, ".")-2:]
		degrees := west[:strings.Index(west, ".")-2]
		floatDegrees, err := strconv.ParseFloat(degrees, 32)
		if err != nil {
			return nil, nil, errors.New("error converting degrees string to float: " + err.Error())
		}
		floatMinutes, err := strconv.ParseFloat(minutes, 32)
		if err != nil {
			return nil, nil, errors.New("error converting string minutes to float: " + err.Error())
		}
		newReading.Latitude = float32(floatDegrees + (floatMinutes / 60))
	} else {
		return nil, nil, errors.New("could not find W or E in raw data data")
	}

	return &newReading, &tk102Identifier, nil
}

func (ts *TK102Server) handleConnection(c net.Conn) {
	log.Info(fmt.Sprintf("TK102 server serving %s", c.RemoteAddr().String()))
	// initialise session variables
	lastReading := reading.Reading{}
	tk102Device := tk1022.TK102{}
	reader := bufio.NewReader(c)
	invalidDataCount := 0
	for {
		data, err := reader.ReadString(')')
		if err != nil {
			fmt.Println("error", err.Error())
			break
		} else {
			newReading, tk102Identifier, err := convertData(data)
			if err != nil {
				invalidDataCount++
				// only allow 3 instances of invalid data
				if invalidDataCount < 3 {
					continue
				} else {
					log.Warn("too many instances of invalid data, terminating connection")
					break
				}
			}

			// check if the device associated with this reading has been retrieved yet
			if tk102Device.Id == "" {
				// if not, retrieve the associated device
				tk102RetrieveResponse := tk102RecordHandler.RetrieveResponse{}
				if err := ts.tk102RecordHandler.Retrieve(&tk102RecordHandler.RetrieveRequest{
					Claims:     *ts.systemClaims,
					Identifier: *tk102Identifier,
				},
					&tk102RetrieveResponse); err != nil {
					log.Warn("cannot find device for reading: ", err.Error())
					break
				}
				tk102Device = tk102RetrieveResponse.TK102
			}

			// if the last reading was not retrieved yet, retrieve it now
			if lastReading.Id == "" {
				collectQuery := query.Query{
					Limit:  1,
					Offset: 0,
					Order:  []query.SortOrder{query.SortOrderDescending},
					SortBy: []string{"timeStamp"},
				}
				collectCriterion := text.Criterion{
					Field: "deviceId.id",
					Text:  tk102Device.Id,
				}
				readingCollectResponse := readingRecordHandler.CollectResponse{}
				if err := ts.readingRecordHandler.Collect(&readingRecordHandler.CollectRequest{
					Claims:   ts.systemClaims,
					Query:    collectQuery,
					Criteria: []criterion.Criterion{collectCriterion},
				},
					&readingCollectResponse); err != nil {
					log.Warn("unable to perform collect for last reading: ", err.Error())
				}
				if len(readingCollectResponse.Records) > 0 {
					lastReading = readingCollectResponse.Records[0]
				}
			}

			diff := reading.DifferenceBetween(newReading, &lastReading)
			if diff < 30 {
				log.Info(fmt.Sprintf("TK102 %s with no. %s new reading not different enough, not recording", tk102Device.ManufacturerId, tk102Device.SimNumber))
				continue
			}

			// use the device to complete the readings associations
			newReading.DeviceId = id.Identifier{Id: tk102Device.Id}
			newReading.DeviceType = device.TK102
			newReading.OwnerPartyType = tk102Device.OwnerPartyType
			newReading.OwnerId = tk102Device.OwnerId
			newReading.AssignedPartyType = tk102Device.AssignedPartyType
			newReading.AssignedId = tk102Device.AssignedId

			// create the reading
			createReadingResponse := readingRecordHandler.CreateResponse{}
			if err := ts.readingRecordHandler.Create(&readingRecordHandler.CreateRequest{
				Reading: *newReading,
			},
				&createReadingResponse); err != nil {
				fmt.Println("error creating new reading: ", err.Error())
				continue
			}

			// set the last reading equal to this one
			lastReading = createReadingResponse.Reading
			log.Info(fmt.Sprintf("%s: %f %f",
				time.Unix(createReadingResponse.Reading.TimeStamp, 0).Format("2006-01-02 3:04PM"),
				createReadingResponse.Reading.Latitude,
				createReadingResponse.Reading.Longitude,
			))
		}
	}

	fmt.Printf("%s disconnected\n", c.RemoteAddr().String())
}

// Start the TK102 Device Server
func (ts *TK102Server) Start() error {
	listener, err := net.Listen("tcp4", fmt.Sprintf("%s:%s", ts.ip, ts.port))
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		c, err := listener.Accept()
		if err != nil {
			return err
		}
		go ts.handleConnection(c)
	}
}
