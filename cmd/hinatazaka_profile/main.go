package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"

	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"

	"github.com/PuerkitoBio/goquery"
)

// プロフィールのstruct
type profile struct {
	X string
	Y string
}

var member [6]profile

func main() {
	var name string
	fmt.Scan(&name) //任意のメンバーを入力

	if !slices.Contains(model.MemberList, name) {
		fmt.Println("人名でない文字列もしくは日向坂46に存在しないメンバーです。")
		os.Exit(0)
	}

	url := "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[name] + "?ima=0000" // 任意のメンバーのURL

	document, _ := scrape.GetDocumentFromURL(url)

	doc := document.Find(".l-contents")

	// セレクタを使って要素を抽出
	cnty := 0
	doc.Find(".c-member__info-td__text").Each(func(index int, element *goquery.Selection) {
		y := strings.TrimSpace(element.Text())
		member[cnty].Y = y
		cnty += 1
	})

	cntx := 0
	doc.Find(".c-member__info-td__name").Each(func(index int, element *goquery.Selection) {
		x := strings.TrimSpace(element.Text())
		member[cntx].X = x
		cntx += 1
	})

	fmt.Println(name)
	fmt.Println(member[0:5])
}
