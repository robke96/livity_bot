package functions

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
)

func SleepCron(s *discordgo.Session) {
	lt, _ := time.LoadLocation("Europe/Vilnius")
	c := cron.New(cron.WithLocation(lt))

	sleepTime, err := strconv.Atoi(os.Getenv("HOUR_TO_SLEEP"))
	if err != nil {
		fmt.Println("Can't find or convert sleepTime to int")
	}

	cronTime := fmt.Sprintf("0 0 %v * * *", sleepTime)

	c.AddFunc(cronTime, func() {
		fmt.Println("it is time")

		channelId := os.Getenv("CHANNEL_ID")
		channel, err := s.State.Channel(channelId)

		if err != nil {
			fmt.Println("Cant grab channel")
		}

		checker(channel)
	})

	c.Start()
}

func checker(ch *discordgo.Channel) {
	retryInterval := 15 * time.Minute // If players found on server recheck after 15 min

	re, err := regexp.Compile("0/20")
	if err != nil {
		fmt.Println("Error compiling regex")
		return
	}

	wakeUpHour, err := strconv.Atoi(os.Getenv("MORNING_WAKEUP_HOUR"))
	if err != nil {
		fmt.Println("Cant convert wakeUpHour to int")
	}

	for {
		isPlayers := re.MatchString(ch.Name)

		if err != nil {
			fmt.Println("Error while matching regex", err)
		}

		if isPlayers {
			currentTime := time.Now()
			targetTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), wakeUpHour, 0, 0, 0, currentTime.Location())

			diff := targetTime.Sub(currentTime).Seconds()
			if diff < 0 {
				targetTime = targetTime.AddDate(0, 0, 1)
				diff = targetTime.Sub(currentTime).Seconds()
			}

			diffToString := strconv.FormatFloat(diff, 'f', 0, 64)

			cmd := exec.Command("sudo", "rtcwake", "-m", "mem", "-s", diffToString)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error exec command", err)
				return
			}

			break
		} else {
			// checking again later
			time.Sleep(retryInterval)
		}
	}
}
