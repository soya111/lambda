package profile

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"notify/pkg/infrastructure/scrape"
	"notify/pkg/model"

	"github.com/PuerkitoBio/goquery"
)

// プロフィールのstruct
type Profile struct {
	birthday   string
	age        string
	sign       string
	height     string
	birthplace string
	bloodtype  string
	ImageUrl   string
}

var (
	ErrNonExistentMember = errors.New("日向坂46に存在しないメンバーです。")
	ErrNoUrl             = errors.New("ポカは日向坂46の一員ですが、URLが存在しません。")
)

// ポカのプロフィール
var PokaProfile = &Profile{
	"2019年12月25日",
	calcAge(time.Date(2019, 12, 25, 0, 0, 0, 0, time.Local), time.Now()),
	"やぎ座",
	"???",
	"???",
	"???",
	"https://cdn.hinatazaka46.com/images/14/8e6/b044f0e534295d2d91700d8613270/1000_1000_102400.jpg",
}

// GetProfileSelectionはメンバーごとのプロフィールが記載されたセレクションを取得
func GetProfileSelection(name string) (*goquery.Selection, error) {
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
func newProfile(birthday, sign, height, birthplace, bloodtype, imageUrl string) (*Profile, error) {
	member := new(Profile)

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
	member.ImageUrl = imageUrl

	return member, err
}

// ScrapeProfileはセレクションからスクレイピングしたプロフィールを取得
func ScrapeProfile(selection *goquery.Selection) *Profile {
	texts := make(map[int]string)
	//セレクタを使って要素を抽出
	selection.Find(".c-member__info-td__text").Each(func(index int, element *goquery.Selection) {
		text := strings.TrimSpace(element.Text())
		texts[index] = text
	})

	selection = selection.Find(".c-member__thumb.c-member__thumb__large")
	element := selection.Find("img").First()
	src, _ := element.Attr("src")

	member, _ := newProfile(texts[0], texts[1], texts[2], texts[3], texts[4], src)
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

// CreateProfileMessageはプロフィールメッセージを生成
func CreateProfileMessage(name string, member *Profile) string {
	message := fmt.Sprintf("%s\n生年月日:%s\n年齢:%s歳\n星座:%s\n身長:%s\n出身地:%s\n血液型:%s\n", name, member.birthday, member.age, member.sign, member.height, member.birthplace, member.bloodtype)
	return message
}
