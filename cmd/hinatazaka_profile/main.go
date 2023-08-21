package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"

	"github.com/PuerkitoBio/goquery"
)

// プロフィールのstruct
type profile struct {
	entry string // 生年月日、星座などのプロフィール項目
	value string //具体的な値
}

// InputNameはコマンド引数により名前を出力
func InputName() string {
	var name string

	flag.StringVar(&name, "name", "hinata", "名前を入力してください")
	flag.Parse()

	return name
}

// GetProfileSelectionはメンバーごとのプロフィールが記載されたセレクションを出力
func GetProfileSelection(name string) *goquery.Selection {
	//入力がポカである場合
	if model.MemberToIdMap[name] == "000" {
		fmt.Println(name)
		fmt.Println("生年月日:2019年12月25日, 星座:山羊座, 身長:???, 出身地:???, 血液型:???")
		os.Exit(0)
	}

	//入力がメンバー名でない場合
	if !model.IsMember(name) {
		fmt.Println("日向坂46に存在しないメンバーです。")
		os.Exit(0)
	}

	url := "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[name] + "?ima=0000" // 任意のメンバーのURL
	document, _ := scrape.GetDocumentFromURL(url)
	selection := document.Find(".l-contents")

	return selection
}

// ScrapeProfileはセレクションからスクレイピングしたプロフィールを出力
func ScrapeProfile(selection *goquery.Selection) [6]profile {
	var member [6]profile
	cntv := 0
	cnte := 0

	//セレクタを使って要素を抽出
	selection.Find(".c-member__info-td__text").Each(func(index int, element *goquery.Selection) {
		v := strings.TrimSpace(element.Text())
		member[cntv].value = v
		cntv += 1
	})

	selection.Find(".c-member__info-td__name").Each(func(index int, element *goquery.Selection) {
		e := strings.TrimSpace(element.Text())
		member[cnte].entry = e
		cnte += 1
	})

	return member
}

// OutputProfileはプロフィールを標準形で出力
func OutputProfile(name string, member [6]profile) {
	fmt.Println(name) //メンバーの名前を出力

	//プロフィールの項目と値をカンマで区切り出力
	for index, prof := range member {
		fmt.Printf("%s:%s", prof.entry, prof.value)

		if index == 4 {
			os.Exit(0) //不要な項目を含むため途中で終了
		}

		fmt.Printf(", ")
	}
}

func main() {
	name := InputName()
	selection := GetProfileSelection(name)
	member := ScrapeProfile(selection)

	OutputProfile(name, member)
}
