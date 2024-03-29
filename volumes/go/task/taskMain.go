package task

import (
	"time"

	"github.com/go-co-op/gocron"
)

func InitTask() {
	cron := gocron.NewScheduler(time.Local)

	// 每天早上抓取每日股價更新股號
	cron.Every(1).Days().At("08:30").Do(RealTimeStockWrite)
	// 09-13 每五分鐘更新一次股價
	cron.Cron("*/5 09-13 * * *").Do(RealTimeStockPrice)

	//每60分抓取會長股
	cron.Every(60).Minutes().Do(GetDirectorStockSelection)

	//每天抓取TWSE檔案 Parse
	cron.Every(1).Days().At("14:30").Do(ParseCSV)

	cron.Every(1).Days().At("17:30").Do(ParseCSV)

	cron.Every(1).Days().At("20:30").Do(ParseCSV)

	cron.Every(1).Days().At("15:05").Do(ParseTpexCSV)

	cron.Every(1).Days().At("17:45").Do(ParseTpexCSV)

	cron.Every(1).Days().At("21:05").Do(ParseTpexCSV)

	// 開始工作排程
	cron.StartAsync()
}
