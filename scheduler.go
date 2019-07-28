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
	go sch.loopSchedule()

	return
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
	sch.lk.Lock()
	defer sch.lk.Unlock()
	sch.createJobPlan(event)

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
		log.Warnf("the job conf:%#v not  exist", jobConf)
		log.Warnf("the job conf:%#v change job create event", jobConf)
		sch.createJobPlan(&JobChangeEvent{
			Type: JobCreateChangeEvent,
			Conf: jobConf,
		})
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
		Group:    jobConf.Group,
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

func (sch *JobScheduler) createJobPlan(event *JobChangeEvent) {

	var (
		err      error
		schedule cron.Schedule
	)

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
		Group:    jobConf.Group,
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

// push a job change event
func (sch *JobScheduler) pushJobChangeEvent(event *JobChangeEvent) {

	sch.eventChan <- event
}

// loop schedule job
func (sch *JobScheduler) loopSchedule() {

	timer := time.NewTimer(time.Second)

	for {

		select {

		case <-timer.C:

		case event := <-sch.eventChan:

			sch.handleJobChangeEvent(event)
		}

		durationTime := sch.trySchedule()
		log.Infof("the durationTime :%d", durationTime)
		timer.Reset(durationTime)
	}

}

// try schedule the job
func (sch *JobScheduler) trySchedule() time.Duration {

	var (
		first bool
	)
	if len(sch.schedulePlans) == 0 {

		return time.Second
	}

	now := time.Now()
	leastTime := new(time.Time)
	first = true
	for _, plan := range sch.schedulePlans {

		scheduleTime := plan.NextTime
		if scheduleTime.Before(now) && sch.node.state == NodeLeaderState {
			log.Infof("schedule execute the plan:%#v", plan)

			snapshot := &JobSnapshot{
				Id:         plan.Id + GenerateSerialNo(),
				JobId:      plan.Id,
				Name:       plan.Name,
				Group:      plan.Group,
				Cron:       plan.Cron,
				Target:     plan.Target,
				Params:     plan.Params,
				Mobile:     plan.Mobile,
				Remark:     plan.Remark,
				CreateTime: ToDateString(now),
			}
			sch.node.exec.pushSnapshot(snapshot)
		}
		nextTime := plan.schedule.Next(now)
		plan.NextTime = nextTime
		plan.BeforeTime = scheduleTime

		// first
		if first {
			first = false
			leastTime = &nextTime
		}

		// check least time after next schedule  time
		if leastTime.After(nextTime) {

			leastTime = &nextTime
		}

	}

	if leastTime.Before(now) {

		return time.Second
	}

	return leastTime.Sub(now)

}
