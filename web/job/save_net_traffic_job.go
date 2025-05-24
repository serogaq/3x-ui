package job

import "x-ui/web/service"

type SaveNetTrafficJob struct {
	server service.ServerService
}

func NewSaveNetTrafficJob() *SaveNetTrafficJob {
	return new(SaveNetTrafficJob)
}

func (j *SaveNetTrafficJob) Run() {
	j.server.SaveNetTraffic()
}
