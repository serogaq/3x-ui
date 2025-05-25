package job

import (
	"x-ui/logger"
	"x-ui/web/service"
)

type LoginStatus byte

const (
	LoginSuccess LoginStatus = 1
	LoginFail    LoginStatus = 0
)

type StatsNotifyJob struct {
	xrayService    service.XrayService
	tgbotService   service.Tgbot
	settingService service.SettingService
	serverService  service.ServerService
}

func NewStatsNotifyJob() *StatsNotifyJob {
	return new(StatsNotifyJob)
}

// Here run is a interface method of Job interface
func (j *StatsNotifyJob) Run() {
	if !j.xrayService.IsXrayRunning() {
		return
	}

	j.tgbotService.SendReport()

	runtime, err := j.settingService.GetTgbotRuntime()
	if err != nil {
		logger.Warning("StatsNotifyJob: get runtime failed", err)
	}

	if err == nil && runtime == "@daily" {
		j.serverService.ResetDailyTraffic("telegram report")
	}
}
