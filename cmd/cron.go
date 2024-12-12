package main

import (
	"gin_example/models"
	log "gin_example/pkg/logging"
	"time"

	"github.com/robfig/cron"
)

func task() {
	log.InfoF("Starting Cron Task")

	// 创建一个空白的Cron job runner
	c := cron.New()
	// ‘*’表示匹配字段的所所有值
	c.AddFunc("* * * * * *", func() {
		log.Info("Run models.CleanAllTag")
		models.CleanAllTag()
	})
	c.AddFunc("* * * * * *", func() {
		log.Info("Run models.CleanAllArticle")
		models.CleanAllArticle()
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		<-t1.C
		t1.Reset(time.Second * 10)
	}
}
