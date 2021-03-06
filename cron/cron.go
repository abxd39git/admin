package cron

import (
	"admin/app/models"
	"admin/utils"
	"github.com/robfig/cron"
	"time"
)

func InitCron() {
	if utils.Cfg.MustBool("cron", "run", false) {
		c := cron.New()
		c.AddFunc("0 0 0 * * *", doTokensDailySheet)   // 凌晨0点
		c.AddFunc("0 0 3 * * *", doTransferDailySheet) // 凌晨3点
		c.Start()
	}
}

// 划转日汇总
func doTransferDailySheet() {
	today := time.Now().Format(utils.LAYOUT_DATE)

	// 开始汇总
	new(models.TransferDailySheet).DoDailySheet(today)
}

// 币种数量汇总
func doTokensDailySheet() {
	today := time.Now().Format(utils.LAYOUT_DATE)

	// 开始汇总
	new(models.TokensDailySheet).DoDailySheet(today)
}
