package main

import (
	"flag"
	"fmt"

	"zephyr/pkg/profile"
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

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(profile.CreateProfileMessage(member))
}
