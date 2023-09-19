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
	Name       string
	Birthday   string
	Age        string
	Sign       string
	Height     string
	Birthplace string
	Bloodtype  string
	ImageUrl   string
}

var ErrNoUrl = errors.New("ポカは日向坂46の一員ですが、URLが存在しません。")

// ポカのプロフィール
var PokaProfile = &Profile{
	"ポカ",
	"2019年12月25日",
	calcAge(time.Date(2019, 12, 25, 0, 0, 0, 0, time.Local), time.Now()),
	"やぎ座",
	"???",
	"???",
	"???",
	"https://cdn.hinatazaka46.com/images/14/8e6/b044f0e534295d2d91700d8613270/1000_1000_102400.jpg",
}

// getProfileSelectionはメンバーごとのプロフィールが記載されたセレクションを取得
func getProfileSelection(name string) (*goquery.Selection, error) {
	// 入力が卒業メンバーである場合
	if model.IsGrad(name) {
		return nil, model.ErrGraduatedMember
	}
	//入力がメンバー名でない場合
	if !model.IsMember(name) {
		return nil, model.ErrNonExistentMember
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
func newProfile(name, birthday, sign, height, birthplace, bloodtype, imageUrl string) *Profile {
	member := new(Profile)

	member.Birthday = birthday
	normalizedBirthday, err := normalizeDate(member.Birthday)
	if err != nil {
		member.Age = "???"
	} else {
		member.Age = calcAge(normalizedBirthday, time.Now())
	}

	member.Name = name
	member.Sign = sign
	member.Height = height
	member.Birthplace = birthplace
	member.Bloodtype = bloodtype
	member.ImageUrl = imageUrl

	return member
}

// ScrapeProfileはセレクションからスクレイピングしたプロフィールを取得
func ScrapeProfile(name string) (*Profile, error) {
	selection, err := getProfileSelection(name)

	if errors.Is(err, ErrNoUrl) {
		return PokaProfile, nil
	}
	if err != nil {
		return nil, err
	}

	texts := make(map[int]string)
	//セレクタを使って要素を抽出
	selection.Find(".c-member__info-td__text").Each(func(index int, element *goquery.Selection) {
		text := strings.TrimSpace(element.Text())
		texts[index] = text
	})

	selection = selection.Find(".c-member__thumb.c-member__thumb__large")
	element := selection.Find("img").First()
	src, _ := element.Attr("src")

	member := newProfile(name, texts[0], texts[1], texts[2], texts[3], texts[4], src)
	return member, nil
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
func CreateProfileMessage(member *Profile) string {
	message := fmt.Sprintf("%s\n生年月日:%s\n年齢:%s歳\n星座:%s\n身長:%s\n出身地:%s\n血液型:%s", member.Name, member.Birthday, member.Age, member.Sign, member.Height, member.Birthplace, member.Bloodtype)
	return message
}
