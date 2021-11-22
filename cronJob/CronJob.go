package cronJob

import (
	"wzm/danmu3.0/service"
	"wzm/danmu3.0/util"

	"github.com/robfig/cron"
)

var Cron *cron.Cron

func CronJob() {
	if Cron == nil {
		Cron = cron.New()
	}

	//每天凌晨1点执行
	err := Cron.AddFunc("0 0 1 * * ?", service.ClicksStoreInDB)
	if err != nil {
		util.Logfile("[Error]", " cornJob error  "+err.Error())
	}
	Cron.Start()
	util.PrintLog("[Info]", " Created cron job")
}
