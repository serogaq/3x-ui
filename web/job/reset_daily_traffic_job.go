package job

import "github.com/mhsanaei/3x-ui/v2/web/service"

type ResetDailyTrafficJob struct {
	server service.ServerService
}

func NewResetDailyTrafficJob() *ResetDailyTrafficJob {
	return new(ResetDailyTrafficJob)
}

func (j *ResetDailyTrafficJob) Run() {
	j.server.ResetDailyTraffic("cron job")
}
