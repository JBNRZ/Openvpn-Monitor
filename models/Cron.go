package models

import (
	"github.com/robfig/cron/v3"
)

func InitCron() *cron.Cron {
	c := cron.New(cron.WithSeconds(), cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)))
	if _, err := c.AddFunc("* * * * * *", cronUpdate); err != nil {
		Logger.Fatalln(err)
	}
	return c
}

func cronUpdate() {
	clients := GetDetail()
	var updated []string
	if len(clients) == 0 {
		updated = []string{"none"}
	} else {
		for _, cli := range clients {
			ok, user := Update(cli)
			if !ok {
				Logger.Warning("Failed to update ", cli.Name)
				continue
			}
			updated = append(updated, user.Name)
		}
		if len(updated) == 0 {
			updated = []string{"none"}
		}
	}
	Offline(updated)
	Sum()
}
