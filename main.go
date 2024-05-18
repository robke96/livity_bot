package main

import (
	"fmt"
	"livitybot/functions"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// get .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error at loading .env: %s", err)
	}

	// grab token
	token := os.Getenv("TOKEN")

	// creating dc session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating discord session: %s", err)
	}

	dg.Identify.Intents = discordgo.IntentsAll

	dg.AddHandler(ready)

	// opening websocket
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection: %s", err)
	}

	// Manual program close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// close session
	dg.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Botas online - CTRL-C to exit")

	go functions.ChangeStatus(s)
	go functions.SleepCron(s)
}
