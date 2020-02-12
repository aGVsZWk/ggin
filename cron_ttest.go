package main

import (
	"github.com/robfig/cron"
	"fmt"
)

func main() {

	i := 0

	c := cron.New()

	spec := "*/5 * * * * ?"

	c.AddFunc(spec, func() {

		i++

		fmt.Println("cron running:", i)

	})

	c.Start()

	select {}

}
