package server

import (
	"bufio"
	"fmt"
	"gitlab.com/iotTracker/brain/log"
	"gitlab.com/iotTracker/brain/tracker/reading"
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	"net"
	"strconv"
	"time"
	"strings"
	"github.com/go-errors/errors"
)

type tk102Server struct {
	readingRecordHandler readingRecordHandler.RecordHandler
	ip                   string
	port                 string
}

func New(
	readingRecordHandler readingRecordHandler.RecordHandler,
	ip string,
	port string,
) *tk102Server {

	return &tk102Server{
		readingRecordHandler: readingRecordHandler,
		ip:                   ip,
		port:                 port,
	}
}

func convertData(raw string) (*reading.Reading, error) {
	// e.g. raw 027045053055BR00190223A2609.0026S02808.1030E000.21226590.000000000000L00000000
	// will be converted to: -26.150043 28.135050
	newReading := reading.Reading{}

	// minimum length of 70?
	if len(raw) < 70 {
		return nil, errors.New("raw data too short")
	}
	newReading.Raw = raw

	newReading.TimeStamp = time.Now().UTC().Unix()

	// check for BR which is used to split out IMEI
	if strings.Contains(raw, "BR") {
		newReading.IMEI = raw[:strings.Index(raw, "BR")]
	} else {
		return nil, errors.New("could not find BR separator in raw data")
	}

	// check for A which separates date and co-ordinates
	if strings.Contains(raw, "A") {
		// confirm N only appears once
		if strings.Count(raw, "A") > 1 {
			return nil, errors.New("more than 1 A in raw data")
		}
	}

	// check for N or S
	nOrS := ""
	if strings.Contains(raw, "N") {
		// confirm N only appears once
		if strings.Count(raw, "N") > 1 {
			return nil, errors.New("more than 1 N in raw data")
		} else {
			// perform processing for N
			nOrS = "N"
			north := raw[strings.Index(raw, "A")+1 : strings.Index(raw, "N")]
			minutes := north[strings.Index(north, ".")-2:]
			degrees := north[:strings.Index(north, ".")-2]
			floatDegrees, err := strconv.ParseFloat(degrees, 32)
			if err != nil {
				return nil, errors.New("error converting degrees string to float: " + err.Error())
			}
			floatMinutes, err := strconv.ParseFloat(minutes, 32)
			if err != nil {
				return nil, errors.New("error converting string minutes to float: " + err.Error())
			}
			newReading.Latitude = float32(floatDegrees + (floatMinutes / 60))
		}
	} else if strings.Contains(raw, "S") {
		// confirm S only appears once
		if strings.Count(raw, "S") > 1 {
			return nil, errors.New("more than 1 S in raw data")
		} else {
			nOrS = "S"
			// perform processing for S
			south := raw[strings.Index(raw, "A")+1 : strings.Index(raw, "S")]
			minutes := south[strings.Index(south, ".")-2:]
			degrees := south[:strings.Index(south, ".")-2]
			floatDegrees, err := strconv.ParseFloat(degrees, 32)
			if err != nil {
				return nil, errors.New("error converting degrees string to float: " + err.Error())
			}
			floatMinutes, err := strconv.ParseFloat(minutes, 32)
			if err != nil {
				return nil, errors.New("error converting string minutes to float: " + err.Error())
			}
			newReading.Latitude = -1 * float32(floatDegrees+(floatMinutes/60))
		}
	} else {
		return nil, errors.New("could not find N or S in raw data data")
	}

	// check for E or W
	if strings.Contains(raw, "E") {
		// confirm N only appears once
		if strings.Count(raw, "E") > 1 {
			return nil, errors.New("more than 1 E in raw data")
		} else {
			// process E
			east := raw[strings.Index(raw, nOrS)+1 : strings.Index(raw, "E")]
			minutes := east[strings.Index(east, ".")-2:]
			degrees := east[:strings.Index(east, ".")-2]
			floatDegrees, err := strconv.ParseFloat(degrees, 32)
			if err != nil {
				return nil, errors.New("error converting degrees string to float: " + err.Error())
			}
			floatMinutes, err := strconv.ParseFloat(minutes, 32)
			if err != nil {
				return nil, errors.New("error converting string minutes to float: " + err.Error())
			}
			newReading.Longitude = float32(floatDegrees + (floatMinutes / 60))
		}
	} else if strings.Contains(raw, "W") {
		// confirm S only appears once
		if strings.Count(raw, "W") > 1 {
			return nil, errors.New("more than 1 W in raw data")
		} else {
			// process W
			west := raw[strings.Index(raw, nOrS)+1 : strings.Index(raw, "W")]
			minutes := west[strings.Index(west, ".")-2:]
			degrees := west[:strings.Index(west, ".")-2]
			floatDegrees, err := strconv.ParseFloat(degrees, 32)
			if err != nil {
				return nil, errors.New("error converting degrees string to float: " + err.Error())
			}
			floatMinutes, err := strconv.ParseFloat(minutes, 32)
			if err != nil {
				return nil, errors.New("error converting string minutes to float: " + err.Error())
			}
			newReading.Latitude = float32(floatDegrees + (floatMinutes / 60))
		}
	} else {
		return nil, errors.New("could not find W or E in raw data data")
	}

	return &newReading, nil
}

func (ts *tk102Server) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	lastReading := reading.Reading{}
	reader := bufio.NewReader(c)
	for {
		data, err := reader.ReadString(')')
		if err != nil {
			fmt.Println("error", err.Error())
			break
		} else {
			newReading, err := convertData(data)
			if err != nil {
				log.Info("Not Recording. Invalid Data: ", err.Error())
				continue
			}

			// if last reading was saved
			if lastReading.Id != "" {
				diff := reading.DifferenceBetween(newReading, &lastReading)
				if diff < 30 {
					log.Info(fmt.Sprintf("New reading not different enough. Not Recording."))
					continue
				}
			}


			// create
			createReadingResponse := readingRecordHandler.CreateResponse{}
			if err := ts.readingRecordHandler.Create(&readingRecordHandler.CreateRequest{
				Reading: *newReading,
			},
				&createReadingResponse); err != nil {
				fmt.Println("error creating new reading: ", err.Error())
				continue
			}
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

func (ts *tk102Server) Start() error {
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
