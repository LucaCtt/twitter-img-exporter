package main

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lucactt/twitter-img-exporter/twitter"
	"github.com/spf13/viper"
)

func main() {
	Run()
}

// Run starts the app CLI.
func Run() {
	client, err := twitter.NewClient(viper.GetString("twitter.key"),
		viper.GetString("twitter.secret"))
	if err != nil {
		fmt.Print(fmt.Errorf("Create Twitter client failed: %w", err))
	}

	user := inputUser()
	src := &TwitterImgSrc{client, user}
	dst := &DirImgDst{inputDir()}

	imgs, err := src.Read()
	if err != nil {
		fmt.Print(fmt.Errorf("Get images failed: %w", err))
	}

	err = dst.Write(imgs)
	if err != nil {
		fmt.Print(fmt.Errorf("Write images failed: %w", err))
	}
}

func inputUser() string {
	name := ""
	prompt := &survey.Input{
		Message: "Twitter screen name: ",
	}
	survey.AskOne(prompt, &name)

	return name
}

func inputDir() string {
	dir := ""
	prompt := &survey.Input{
		Message: "Output dir: ",
	}
	survey.AskOne(prompt, &dir)

	return dir
}
