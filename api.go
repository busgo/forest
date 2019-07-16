package forest

import (
	"github.com/labstack/echo"
	"github.com/robfig/cron"
	"net/http"
)

type JobAPi struct {
	node *JobNode
	echo *echo.Echo
}

func NewJobAPi(node *JobNode) (api *JobAPi) {

	api = &JobAPi{
		node: node,
	}
	e := echo.New()
	e.POST("/job/add", api.AddJob)
	e.POST("/job/edit", api.editJob)
	e.POST("/job/delete", api.deleteJob)
	e.POST("/job/list", api.jobList)
	e.POST("/group/add", api.addGroup)
	e.POST("/group/list", api.groupList)
	e.POST("/node/list", api.nodeList)
	go e.Logger.Fatal(e.Start(node.apiAddress))
	api.echo = e
	return
}

// add a new job
func (api *JobAPi) AddJob(context echo.Context) (err error) {

	var (
		message string
	)
	jobConf := new(JobConf)
	if err = context.Bind(jobConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if jobConf.Name == "" {
		message = "任务名称不能为空"
		goto ERROR
	}
	if jobConf.Group == "" {
		message = "任务分组不能为空"
		goto ERROR
	}

	if jobConf.Cron == "" {
		message = "任务Cron表达式不能为空"
		goto ERROR
	}

	if _, err = cron.Parse(jobConf.Cron); err != nil {
		message = "非法的Cron表达式"
		goto ERROR
	}

	if jobConf.Target == "" {
		message = "任务Target不能为空"
		goto ERROR
	}

	if jobConf.Status == 0 {
		message = "任务状态不能为空"
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: jobConf, Message: "创建成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// edit a job
func (api *JobAPi) editJob(context echo.Context) (err error) {

	var (
		message string
	)
	jobConf := new(JobConf)
	if err = context.Bind(jobConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if jobConf.Id == "" {
		message = "此任务记录不存在"
		goto ERROR
	}
	if jobConf.Name == "" {
		message = "任务名称不能为空"
		goto ERROR
	}
	if jobConf.Group == "" {
		message = "任务分组不能为空"
		goto ERROR
	}

	if jobConf.Cron == "" {
		message = "任务Cron表达式不能为空"
		goto ERROR
	}

	if _, err = cron.Parse(jobConf.Cron); err != nil {
		message = "非法的Cron表达式"
		goto ERROR
	}

	if jobConf.Target == "" {
		message = "任务Target不能为空"
		goto ERROR
	}

	if jobConf.Status == 0 {
		message = "任务状态不能为空"
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: jobConf, Message: "修改成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// job  list
func (api *JobAPi) jobList(context echo.Context) (err error) {

	return context.JSON(http.StatusOK, Result{Code: 0, Message: "查询成成功"})

}

// delete a job
func (api *JobAPi) deleteJob(context echo.Context) (err error) {

	var (
		message string
	)
	job := new(JobConf)
	if err = context.Bind(job); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if job.Id == "" {
		message = "此任务记录不存在"
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: job, Message: "删除成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// add a job group
func (api *JobAPi) addGroup(context echo.Context) (err error) {

	var (
		message string
	)
	groupConf := new(GroupConf)
	if err = context.Bind(groupConf); err != nil {

		message = "请求参数不能为空"
		goto ERROR
	}

	if groupConf.Name == "" {
		message = "任务集群名称不能为空"
		goto ERROR
	}

	if groupConf.Remark == "" {
		message = "任务集群描述"
		goto ERROR
	}

	return context.JSON(http.StatusOK, Result{Code: 0, Data: groupConf, Message: "添加成功"})

ERROR:
	return context.JSON(http.StatusOK, Result{Code: -1, Message: message})
}

// job group list
func (api *JobAPi) groupList(context echo.Context) (err error) {

	return context.JSON(http.StatusOK, Result{Code: 0, Message: "查询成成功"})

}

// job node list
func (api *JobAPi) nodeList(context echo.Context) (err error) {

	return context.JSON(http.StatusOK, Result{Code: 0, Message: "查询成成功"})

}
