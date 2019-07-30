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

func ParkJobExecuteSnapshot(snapshot *JobExecuteSnapshot) (value []byte, err error) {
	value, err = json.Marshal(snapshot)
	return

}

func UParkJobExecuteSnapshot(value []byte) (snapshot *JobExecuteSnapshot, err error) {

	snapshot = new(JobExecuteSnapshot)
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

func TimeSubDays(t1, t2 time.Time) int {

	if t1.Location().String() != t2.Location().String() {
		return -1
	}
	hours := t1.Sub(t2).Hours()

	if hours <= 0 {
		return -1
	}
	// sub hours less than 24
	if hours < 24 {
		// may same day
		t1y, t1m, t1d := t1.Date()
		t2y, t2m, t2d := t2.Date()
		isSameDay := (t1y == t2y && t1m == t2m && t1d == t2d)

		if isSameDay {

			return 0
		} else {
			return 1
		}

	} else { // equal or more than 24

		if (hours/24)-float64(int(hours/24)) == 0 { // just 24's times
			return int(hours / 24)
		} else { // more than 24 hours
			return int(hours/24) + 1
		}
	}

}

func ParseInLocation(value string) (dateTime time.Time, err error) {

	dateTime, err = time.Parse("2006-01-01 15:04:05", value)
	return
}
