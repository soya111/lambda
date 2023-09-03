package main

import (
	"errors"
	"flag"
	"fmt"

	"notify/pkg/profile"
)

var name string

func init() {
	flag.StringVar(&name, "name", "hinata", "名前を入力してください")
}

// inputNameはコマンド引数により名前を取得
func inputName() string {
	flag.Parse()
	return name
}

func main() {
	name := inputName()
	selection, err := profile.GetProfileSelection(name)

	if err != nil {
		if errors.Is(err, profile.ErrNoUrl) {
			fmt.Print(profile.CreateProfileMessage(name, profile.PokaProfile))
		} else {
			fmt.Println(err)
		}
		return
	}

	member := profile.ScrapeProfile(selection)

	fmt.Print(profile.CreateProfileMessage(name, member))
}
