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

// エラーメッセージのマップ
var errormap = map[string]error{
	"0": fmt.Errorf("日向坂46に存在しないメンバーです。"),
	"1": fmt.Errorf("ポカは日向坂46の一員ですが、URLが存在しません。"),
}

// ポカのプロフィール
var pokaprofile = [6]profile{
	{"生年月日", "2019年12月25日"},
	{"星座", "山羊座"},
	{"身長", "???"},
	{"出身地", "???"},
	{"血液型", "???"},
}

var errorcode string // エラーコード

// inputNameはコマンド引数により名前を取得
func inputName() string {
	var name string

	flag.StringVar(&name, "name", "hinata", "名前を入力してください")
	flag.Parse()

	return name
}

// getProfileSelectionはメンバーごとのプロフィールが記載されたセレクションを取得
func getProfileSelection(name string) (*goquery.Selection, error) {
	//入力がメンバー名でない場合
	if !model.IsMember(name) {
		errorcode = "0"
		return nil, errormap[errorcode]
	}

	//入力がポカである場合
	if model.MemberToIdMap[name] == "000" {
		errorcode = "1"
		return nil, errormap[errorcode]
	}

	url := "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[name] + "?ima=0000" // 任意のメンバーのURL
	document, _ := scrape.GetDocumentFromURL(url)
	selection := document.Find(".l-contents")

	return selection, nil
}

// scrapeProfileはセレクションからスクレイピングしたプロフィールを取得
func scrapeProfile(selection *goquery.Selection) [6]profile {
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

// outputProfileはプロフィールを標準形で出力
func outputProfile(name string, member [6]profile) {
	fmt.Println(name) //メンバーの名前を出力

	//プロフィールの項目と値をカンマで区切り出力
	var profile []string
	for _, prof := range member[:5] {
		profile = append(profile, prof.entry+":"+prof.value)
	}

	fmt.Println(strings.Join(profile, ", "))
}

func main() {
	name := inputName()
	selection, err := getProfileSelection(name)

	if err != nil {
		// ポカである場合のみポカのプロフィールを出力
		if errorcode == "1" {
			outputProfile(name, pokaprofile)
			return
		} else {
			fmt.Println(err)
			return
		}
	}

	member := scrapeProfile(selection)

	outputProfile(name, member)
}
