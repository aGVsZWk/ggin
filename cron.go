package main

import (
	"log"

	"github.com/robfig/cron"

	"ggin/models"
)

func main() {
	log.Println("Starting...")
	c := cron.New()

	c.AddFunc("*/5 * * * * ?", func() {
		log.Println("Run models.CLeanAllTag...")
		models.CleanAllTag()
	})
	c.AddFunc("*/5 * * * * ?", func() {
		log.Println("Run models.CLeanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()
	log.Println("cron start")
	select {
	}
	//inspect(c.Entries())
	//c.Stop()
	//t1 := time.NewTimer(time.Second * 10)
	//for {
	//	select {
	//	case <- t1.C:
	//		t1.Reset(time.Second * 10)
	//
	//	}
	//}
}