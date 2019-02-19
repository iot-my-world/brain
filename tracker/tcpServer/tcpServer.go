package tcpServer

import (
	readingRecordHandler "gitlab.com/iotTracker/brain/tracker/reading/recordHandler"
	"fmt"
	"net"
	"bufio"
	"gitlab.com/iotTracker/brain/tracker/reading"
	"time"
	"strconv"
	"gitlab.com/iotTracker/brain/log"
)

type tcpServer struct {
	readingRecordHandler readingRecordHandler.RecordHandler
	ip                   string
	port                 string
}

func New(
	readingRecordHandler readingRecordHandler.RecordHandler,
	ip string,
	port string,
) *tcpServer {

	return &tcpServer{
		readingRecordHandler: readingRecordHandler,
		ip:                   ip,
		port:                 port,
	}
}

// 027045053055BR00190217A2609.0064S02808.0845E000.30840070.000000000000L00000000
var exampleData = "(027045053055BR03190217A2608.9903S02808.1080E000.70610270.000000000000L00000000)"

/**
2609.0064 ==> 26°09'(60*0.0064)"
 */

 //26°09'0.0313 ==> 26°09'1.878000"
func toDMS(reading string) (string, error) {
	degrees := reading[:2]
	minutes := reading[2:4]
	fractionOfAMinute, err := strconv.ParseFloat(fmt.Sprintf("0.%s", reading[5:]), 32)
	if err != nil {
		return "", err
	}
	seconds := fmt.Sprintf("%f", 60*fractionOfAMinute)

	return fmt.Sprintf("%s°%s'%s\"", degrees, minutes, seconds), nil
}

func (ts *tcpServer) handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	reader := bufio.NewReader(c)
	for {
		data, err := reader.ReadString(')')
		if err != nil {
			fmt.Println("error", err.Error())
			break
		} else {
			if len(data) < len(exampleData) {
				fmt.Printf("Not Recording %s\n", data)
			} else {
				// parse readings
				south, err := toDMS(data[24:33])
				if err != nil {
					log.Error(fmt.Sprintf("unable to parse south reading %s to dms: %s", data[24:33], err.Error()))
					continue
				}
				east, err := toDMS(data[35:44])
				if err != nil {
					log.Error(fmt.Sprintf("unable to parse south reading %s to dms: %s", data[35:44], err.Error()))
					continue
				}

				// create
				createReadingResponse := readingRecordHandler.CreateResponse{}
				if err := ts.readingRecordHandler.Create(&readingRecordHandler.CreateRequest{
					Reading: reading.Reading{
						IMEI:      data[2:13],
						Raw:       data,
						TimeStamp: time.Now().UTC().Unix(),
						SouthCoordinate:south,
						EastCoordinate:east,
					},
				},
					&createReadingResponse);
					err != nil {
					fmt.Println("error creating new reading: ", err.Error())
					continue
				}
				fmt.Printf("%s: %sS %sE\n",
					time.Unix(createReadingResponse.Reading.TimeStamp, 0).Format("2006-01-02 3:04PM"),
					createReadingResponse.Reading.SouthCoordinate,
					createReadingResponse.Reading.EastCoordinate,
				)
			}
		}
	}

	fmt.Printf("%s disconnected\n", c.RemoteAddr().String())
}

func (ts *tcpServer) Start() error {
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
