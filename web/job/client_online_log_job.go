package job

import (
	"time"
	"github.com/mhsanaei/3x-ui/v2/web/service"
)

type ClientOnlineLogJob struct {
	settingService service.SettingService
	inboundService service.InboundService
}

func NewClientOnlineLogJob() *ClientOnlineLogJob {
	return new(ClientOnlineLogJob)
}

func (j *ClientOnlineLogJob) Run() {
	limit, err := j.settingService.GetClientConnLog()
	if err != nil || limit <= 0 {
		return
	}
	now := time.Now()
	clients := j.inboundService.GetOnlineClients()
	for _, email := range clients {
		j.inboundService.AddClientOnlineLog(email, now, limit)
	}
}
