package job

import "x-ui/web/service"

type ResetDailyTrafficJob struct {
	server service.ServerService
}

func NewResetDailyTrafficJob() *ResetDailyTrafficJob {
	return new(ResetDailyTrafficJob)
}

func (j *ResetDailyTrafficJob) Run() {
	j.server.ResetDailyTraffic("cron job")
}
