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
	member, err := profile.ScrapeProfile(name)

	if errors.Is(err, profile.ErrNoUrl) {
		fmt.Println(profile.CreateProfileMessage(member))
		return
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(profile.CreateProfileMessage(member))
}
