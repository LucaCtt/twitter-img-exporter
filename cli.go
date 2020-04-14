package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lucactt/twitter-img-exporter/twitter"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/twitter_img_exporter/")
	viper.AddConfigPath("$HOME/.twitter_img_exporter/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	Run()
}

// Run starts the app CLI.
func Run() {
	checkAuth()

	client, err := twitter.NewClient(viper.GetString("twitter.key"),
		viper.GetString("twitter.secret"))
	if err != nil {
		fmt.Print(fmt.Errorf("Create Twitter client failed: %w", err))
	}

	user := input("Twitter screen name: ")
	src := &TwitterImgSrc{client, user}
	dst := &DirImgDst{input("Output dir: ")}

	imgs, err := src.Read()
	if err != nil {
		fmt.Print(fmt.Errorf("Get images failed: %w", err))
	}

	err = dst.Write(imgs)
	if err != nil {
		fmt.Print(fmt.Errorf("Write images failed: %w", err))
	}
}

func checkAuth() {
	if viper.GetString("twitter.key") == "" {
		key := input("Twitter API Key: ")
		viper.Set("twitter.key", key)
		viper.WriteConfig()
	}

	if viper.GetString("twitter.secret") == "" {
		key := input("Twitter API secret: ")
		viper.Set("twitter.secret", key)
		viper.WriteConfig()
	}
}

func input(p string) string {
	in := ""
	prompt := &survey.Input{
		Message: p,
	}
	survey.AskOne(prompt, &in)

	return in
}
