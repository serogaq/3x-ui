package job

import (
	"strconv"
	"time"

	"x-ui/web/service"

	"github.com/shirou/gopsutil/v4/mem"
)

type CheckMemJob struct {
	tgbotService   service.Tgbot
	settingService service.SettingService
	serverService  service.ServerService
}

func NewCheckMemJob() *CheckMemJob {
	return new(CheckMemJob)
}

func (j *CheckMemJob) round(num float64) int {
    return int(num + math.Copysign(0.5, num))
}

func (j *CheckMemJob) toFixed(num float64, precision int) float64 {
    output := math.Pow(10, float64(precision))
    return float64(j.round(num * output)) / output
}

// Here run is a interface method of Job interface
func (j *CheckMemJob) Run() {
	threshold, _ := j.settingService.GetTgMem()
	needRestart, _ := j.settingService.GetRestartAtMemThreshold()

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		logger.Error("CheckMemJob -- get virtual memory failed:", err)
	} else {
		currentMem := memInfo.Used
		totalMem := memInfo.Total
		percentMem := int(j.toFixed(currentMem / totalMem * 100, 0))

		if percentMem >= int(threshold) && bool(needRestart) == true {
			msg := j.tgbotService.I18nBot("tgbot.messages.memThreshold", "Threshold=="+strconv.Itoa(threshold))
			j.tgbotService.SendMsgToTgbotAdmins(msg)

			err := a.serverService.RestartXrayService()
			if err != nil {
				logger.Error("CheckMemJob -- RestartXrayService failed:", err)
			} else {
				logger.Info("CheckMemJob -- RestartXrayService success")
			}
		}
	}
}
