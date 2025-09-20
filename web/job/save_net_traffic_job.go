package job

import "github.com/mhsanaei/3x-ui/v2/web/service"

type SaveNetTrafficJob struct {
	server service.ServerService
}

func NewSaveNetTrafficJob() *SaveNetTrafficJob {
	return new(SaveNetTrafficJob)
}

func (j *SaveNetTrafficJob) Run() {
	j.server.SaveNetTraffic()
}
