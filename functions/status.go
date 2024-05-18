package functions

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mcstatus-io/go-mcstatus"
)

const offlineName = "🔴 Serveris Išjungtas!"

func ChangeStatus(s *discordgo.Session) {
	host := os.Getenv("SERVER_IP")
	port := os.Getenv("PORT")
	channelId := os.Getenv("CHANNEL_ID")

	if host == "" || port == "" {
		log.Fatalln("Can't find host ip or port.")
	}

	convertedPort, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		fmt.Printf("problem at converting port string to uint64: %s", err)
	}

	if channelId == "" {
		log.Fatalln("Can't find channelId in .env")
	}

	// intervalas
	ticker := time.NewTicker(5 * time.Second) // how often should check status
	defer ticker.Stop()

	for range ticker.C {
		go func() {
			channel, err := s.State.Channel(channelId)
			if err != nil {
				fmt.Println("Can't fetch channel")
			}

			result, err := mcstatus.GetJavaStatus(host, uint16(convertedPort))
			if err != nil {
				fmt.Printf("Can't get server status: %v", err)
			}

			onlineName := "🟢 Žaidžia: " + fmt.Sprintf("%d/20", result.Players.Online)

			if !result.Online {
				if channel.Name == offlineName {
					fmt.Println("ignoruoju1")
					return
				}

				_, err = s.ChannelEdit(channelId, &discordgo.ChannelEdit{
					Name: offlineName,
				})

				if err != nil {
					fmt.Println("Error editing channel")
					return
				}
			} else {
				if channel.Name == onlineName {
					return
				}

				_, err = s.ChannelEdit(channelId, &discordgo.ChannelEdit{
					Name: onlineName,
				})

				if err != nil {
					fmt.Println("error at editing online channel name")
				}
			}
		}()
	}
}
