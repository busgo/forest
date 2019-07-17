package forest

import (
	"github.com/labstack/gommon/log"
	"github.com/robfig/cron"
	"sync"
	"time"
)

// job scheduler
type JobScheduler struct {
	node          *JobNode
	eventChan     chan *JobChangeEvent
	schedulePlans map[string]*SchedulePlan
	lk            *sync.RWMutex
}

func NewJobScheduler(node *JobNode) (sch *JobScheduler) {

	sch = &JobScheduler{
		node:          node,
		eventChan:     make(chan *JobChangeEvent, 250),
		schedulePlans: make(map[string]*SchedulePlan),
		lk:            &sync.RWMutex{},
	}

	go sch.lookup()

	return
}

// look up the job change event chan
func (sch *JobScheduler) lookup() {

	log.Printf("the job scheduler init success!")
	for event := range sch.eventChan {

		sch.handleJobChangeEvent(event)
	}
}

// handle the job change event
func (sch *JobScheduler) handleJobChangeEvent(event *JobChangeEvent) {

	switch event.Type {
	case JobCreateChangeEvent:
		sch.handleJobCreateEvent(event)
	case JobUpdateChangeEvent:
		sch.handleJobUpdateEvent(event)
	case JobDeleteChangeEvent:
		sch.handleJobDeleteEvent(event)
	}
}

// handle the job create event
func (sch *JobScheduler) handleJobCreateEvent(event *JobChangeEvent) {

	var (
		err      error
		schedule cron.Schedule
	)
	sch.lk.Lock()
	defer sch.lk.Unlock()

	jobConf := event.Conf

	if _, ok := sch.schedulePlans[jobConf.Id]; ok {
		log.Warnf("the job conf:%#v exist", jobConf)
		return
	}

	if jobConf.Status == JobStopStatus {

		log.Warnf("the job conf:%#v status is stop", jobConf)
		return
	}

	if schedule, err = cron.Parse(jobConf.Cron); err != nil {
		log.Errorf("the job conf:%#v cron is error exp ", jobConf)
		return
	}

	// build schedule plan
	plan := &SchedulePlan{
		Id:       jobConf.Id,
		Name:     jobConf.Name,
		Cron:     jobConf.Cron,
		Target:   jobConf.Target,
		Params:   jobConf.Params,
		Mobile:   jobConf.Mobile,
		Remark:   jobConf.Remark,
		schedule: schedule,
		NextTime: schedule.Next(time.Now()),
	}

	sch.schedulePlans[jobConf.Id] = plan

	log.Printf("the job conf:%#v create a new schedule plan:%#v", jobConf, plan)

}

// handle the job update event
func (sch *JobScheduler) handleJobUpdateEvent(event *JobChangeEvent) {

	var (
		err      error
		schedule cron.Schedule
	)
	sch.lk.Lock()
	defer sch.lk.Unlock()

	jobConf := event.Conf

	if _, ok := sch.schedulePlans[jobConf.Id]; !ok {
		log.Errorf("the job conf:%#v not  exist", jobConf)
		return
	}

	// stop must delete from the job schedule plan list
	if jobConf.Status == JobStopStatus {

		log.Warnf("the job conf:%#v status is stop must delete from the schedule plan ", jobConf)
		delete(sch.schedulePlans, jobConf.Id)
		return
	}

	if schedule, err = cron.Parse(jobConf.Cron); err != nil {
		log.Errorf("the job conf:%#v  parse the cron error:%#v", jobConf, err)
		return
	}

	// build schedule plan
	plan := &SchedulePlan{
		Id:       jobConf.Id,
		Name:     jobConf.Name,
		Cron:     jobConf.Cron,
		Target:   jobConf.Target,
		Params:   jobConf.Params,
		Mobile:   jobConf.Mobile,
		Remark:   jobConf.Remark,
		schedule: schedule,
		NextTime: schedule.Next(time.Now()),
	}

	// update the schedule plan
	sch.schedulePlans[jobConf.Id] = plan
	log.Printf("the job conf:%#v update a new schedule plan:%#v", jobConf, plan)
}

// handle the job delete event
func (sch *JobScheduler) handleJobDeleteEvent(event *JobChangeEvent) {

	var (
		plan *SchedulePlan
		ok   bool
	)
	sch.lk.Lock()
	defer sch.lk.Unlock()

	jobConf := event.Conf

	if plan, ok = sch.schedulePlans[jobConf.Id]; !ok {
		log.Printf("the job conf:%#v not  exist", jobConf)
		return
	}
	log.Warnf("the job conf:%#v delete a  schedule plan:%#v", jobConf, plan)
	delete(sch.schedulePlans, jobConf.Id)

}

// push a job change event
func (sch *JobScheduler) pushJobChangeEvent(event *JobChangeEvent) {

	sch.eventChan <- event
}
