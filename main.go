package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/spf13/viper"
)

type Config struct {
	Email   string
	Pass    string
	Journal string
}

var logger log.Logger

func main() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("$HOME/.config/vrtoday")
	viper.AddConfigPath("/etc/vrtoday")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("failed to parse config file: %v", err)
	}

	conf := Config{}
	viper.Unmarshal(&conf)

	date := time.Now().Local()
	datePlayer := date
	if date.Hour() < 19 {
		datePlayer = date.Add(-24 * time.Hour)
	}

	vrtDate := datePlayer.Format("20060102")
	url := fmt.Sprintf(
		"https://www.vrt.be/vrtnu/a-z/het-journaal/2021/het-journaal-het-journaal-%s-%s",
		conf.Journal,
		vrtDate,
	)

	options := fmt.Sprintf("--ytdl-raw-options=username=%s,password=\"%s\"", conf.Email, conf.Pass)

	cmd := exec.Command("mpv",
		options,
		url,
	)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to play news: %v", err)
	}
}
