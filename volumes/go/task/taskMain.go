package task

import (
	"time"

	"github.com/go-co-op/gocron"
)

func InitTask() {
	cron := gocron.NewScheduler(time.Local)

	//每60分抓取會長股
	cron.Every(60).Minutes().Do(GetDirectorStockSelection)

	//每天抓取TWSE檔案 Parse
	cron.Every(1).Days().At("14:30").Do(ParseCSV)
	cron.Every(1).Days().At("15:05").Do(ParseTpexCSV)

	// 開始工作排程
	cron.StartAsync()
}
