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
	ticker := time.NewTicker(40 * time.Second) // how often should check status
	defer ticker.Stop()

	for range ticker.C {
		channel, err := s.Channel(channelId)

		go func() {
			if err != nil {
				log.Fatalln("Can't fetch channel")
				return
			}

			result, err := mcstatus.GetJavaStatus(host, uint16(convertedPort))
			if err != nil {
				fmt.Printf("Can't get server status: %v", err)
				return
			}

			if !result.Online {
				if channel.Name == offlineName {
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
				onlineName := "🟢 Žaidžia: " + fmt.Sprintf("%d/20", result.Players.Online)

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
