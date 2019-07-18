package forest

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

func GenerateSerialNo() string {

	now := time.Now()

	format := now.Format("20060101150405")

	suffer := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	return fmt.Sprintf("%s%s", format, suffer)

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
