package forest

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	"math/rand"
	"net"
	"time"
)

func GenerateSerialNo() string {

	now := time.Now()

	format := now.Format("20060101150405")

	suffer := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	return fmt.Sprintf("%s%s", format, suffer)

}

func ToDateString(date time.Time) string {

	return date.Format("2006-01-01 15:04:05")
}
func ParkJobConf(jobConf *JobConf) (value []byte, err error) {

	value, err = json.Marshal(jobConf)
	return
}
func UParkJobConf(value []byte) (jobConf *JobConf, err error) {

	jobConf = new(JobConf)
	err = json.Unmarshal(value, jobConf)
	return
}

func ParkGroupConf(groupConf *GroupConf) (value []byte, err error) {

	value, err = json.Marshal(groupConf)
	return
}

func UParkGroupConf(value []byte) (groupConf *GroupConf, err error) {

	groupConf = new(GroupConf)
	err = json.Unmarshal(value, groupConf)
	return
}

func ParkJobSnapshot(snapshot *JobSnapshot) (value []byte, err error) {

	value, err = json.Marshal(snapshot)
	return
}

func UParkJobSnapshot(value []byte) (snapshot *JobSnapshot, err error) {

	snapshot = new(JobSnapshot)
	err = json.Unmarshal(value, snapshot)
	return
}

func GetLocalIpAddress() (ip string) {

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Warnf("err:%#v", err)
		return
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}

	ip = "127.0.0.1"
	return
}
