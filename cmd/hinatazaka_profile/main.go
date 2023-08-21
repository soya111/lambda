package main

import (
	"flag"
	"fmt"
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

var member [6]profile

func main() {
	var name string

	flag.StringVar(&name, "name", "hinata", "名前を入力してください")
	flag.Parse()

	if model.MemberToIdMap[name] == "000" {
		fmt.Println(name)
		fmt.Println("生年月日:2019年12月25日, 星座:山羊座, 身長:???, 出身地:???, 血液型:???")
		return
	}

	if !model.IsMember(name) {
		fmt.Println("人名でない文字列もしくは日向坂46に存在しないメンバーです。")
		return
	}

	url := "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[name] + "?ima=0000" // 任意のメンバーのURL

	document, _ := scrape.GetDocumentFromURL(url)

	doc := document.Find(".l-contents")

	// セレクタを使って要素を抽出
	cntv := 0
	doc.Find(".c-member__info-td__text").Each(func(index int, element *goquery.Selection) {
		v := strings.TrimSpace(element.Text())
		member[cntv].value = v
		cntv += 1
	})

	cnte := 0
	doc.Find(".c-member__info-td__name").Each(func(index int, element *goquery.Selection) {
		e := strings.TrimSpace(element.Text())
		member[cnte].entry = e
		cnte += 1
	})

	fmt.Println(name)

	//プロフィールの項目と値をカンマで区切り出力
	for index, prof := range member {
		fmt.Printf("%s:%s", prof.entry, prof.value)
		if index == 4 {
			return //不要な項目を含むため途中で終了
		}
		fmt.Printf(", ")
	}
}
