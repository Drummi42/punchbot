package main

import (
	"fmt"
	"github.com/drummi42/punchbot/bot"
	"github.com/drummi42/punchbot/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
