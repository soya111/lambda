package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"

	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"

	"github.com/PuerkitoBio/goquery"
)

// プロフィールのstruct
type profile struct {
	birthday   string
	age        string
	sign       string
	height     string
	birthplace string
	bloodtype  string
	imageUrl   string
}

var (
	ErrNonExistentMember = errors.New("日向坂46に存在しないメンバーです。")
	ErrNoUrl             = errors.New("ポカは日向坂46の一員ですが、URLが存在しません。")
)

// ポカのプロフィール
var pokaProfile = &profile{
	"2019年12月25日",
	calcAge(time.Date(2019, 12, 25, 0, 0, 0, 0, time.Local), time.Now()),
	"やぎ座",
	"???",
	"???",
	"???",
	"https://cdn.hinatazaka46.com/images/14/8e6/b044f0e534295d2d91700d8613270/1000_1000_102400.jpg",
}

var name string

func init() {
	flag.StringVar(&name, "name", "hinata", "名前を入力してください")
}

// inputNameはコマンド引数により名前を取得
func inputName() string {
	flag.Parse()
	return name
}

// getProfileSelectionはメンバーごとのプロフィールが記載されたセレクションを取得
func getProfileSelection(name string) (*goquery.Selection, error) {
	//入力がメンバー名でない場合
	if !model.IsMember(name) {
		return nil, ErrNonExistentMember
	}

	//入力がポカである場合
	if model.MemberToIdMap[name] == "000" {
		return nil, ErrNoUrl
	}

	url := "https://www.hinatazaka46.com/s/official/artist/" + model.MemberToIdMap[name] + "?ima=0000" // 任意のメンバーのURL
	document, _ := scrape.GetDocumentFromURL(url)
	selection := document.Find(".l-contents")

	return selection, nil
}

// newProfileは新しいprofileをつくるコンストラクタ
func newProfile(birthday, sign, height, birthplace, bloodtype, imageUrl string) (*profile, error) {
	member := new(profile)

	member.birthday = birthday
	normalizedBirthday, err := normalizeDate(member.birthday)
	if err != nil {
		member.age = "???"
	} else {
		member.age = calcAge(normalizedBirthday, time.Now())
	}
	member.sign = sign
	member.height = height
	member.birthplace = birthplace
	member.bloodtype = bloodtype
	member.imageUrl = imageUrl

	return member, err
}

// scrapeProfileはセレクションからスクレイピングしたプロフィールを取得
func scrapeProfile(selection *goquery.Selection) *profile {
	texts := make(map[int]string)
	//セレクタを使って要素を抽出
	selection.Find(".c-member__info-td__text").Each(func(index int, element *goquery.Selection) {
		text := strings.TrimSpace(element.Text())
		texts[index] = text
	})

	var srcs []string
	selection = selection.Find(".c-member__thumb.c-member__thumb__large")
	selection.Find("img").Each(func(index int, element *goquery.Selection) {
		src, exists := element.Attr("src")
		if exists {
			srcs = append(srcs, src)
		}
	})

	member, _ := newProfile(texts[0], texts[1], texts[2], texts[3], texts[4], srcs[0])
	return member
}

// normalizeDateは"YYYY年MM月DD日"を標準化したtime.Time型で出力
func normalizeDate(date string) (time.Time, error) {
	layout := "2006年1月2日"

	return time.Parse(layout, date)
}

// calcAgeは生年月日から年齢を取得
func calcAge(birthday time.Time, now time.Time) string {
	//今日の年月日を取得
	thisYear, thisMonth, day := now.Date()

	//年から年齢を計算
	age := thisYear - birthday.Year()

	// 誕生日を迎えていない場合はageを「−1」する
	if thisMonth < birthday.Month() || (thisMonth == birthday.Month() && day < birthday.Day()) {
		age -= 1
	}

	return strconv.Itoa(age)
}

// isTodayBirthdayは今日が誕生日の場合にtrueを返す
func isTodayBirthday(birthday time.Time) bool {
	//今日の年月日を取得
	now := time.Now()
	_, thisMonth, day := now.Date()

	//今日が誕生日の場合にtrueを返す
	return thisMonth == birthday.Month() && day == birthday.Day()
}

// outputProfileはプロフィールを標準形で出力
func outputProfile(name string, member *profile) {
	fmt.Println(name) //メンバーの名前を出力
	fmt.Printf("生年月日:%s, 年齢:%s歳, 星座:%s, 身長:%s, 出身地:%s, 血液型:%s\n%s<---画像はここをクリック\n", member.birthday, member.age, member.sign, member.height, member.birthplace, member.bloodtype, member.imageUrl)

	normalizedBirthday, err := normalizeDate(member.birthday)

	if err != nil {
		fmt.Println(err)
	}

	if isTodayBirthday(normalizedBirthday) {
		fmt.Println("Happy birthday!!")
	}
}

func main() {
	name := inputName()
	selection, err := getProfileSelection(name)

	if err != nil {
		if errors.Is(err, ErrNoUrl) {
			outputProfile(name, pokaProfile)
		} else {
			fmt.Println(err)
		}
		return
	}

	member := scrapeProfile(selection)

	outputProfile(name, member)
}
