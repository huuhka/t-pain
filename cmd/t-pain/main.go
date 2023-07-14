package main

import "t-pain/pkg/tgbot"

func main() {
	err := tgbot.Run()
	if err != nil {
		panic(err)
	}
}