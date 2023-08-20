package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// プロフィールのstruct
type profile struct {
	X string
	Y string
}

var member [6]profile

var (
	// メンバーリスト.
	MemberList = []string{
		"潮紗理菜",
		"影山優佳",
		"加藤史帆",
		"齊藤京子",
		"佐々木久美",
		"佐々木美玲",
		"高瀬愛奈",
		"高本彩花",
		"東村芽依",
		"金村美玖",
		"河田陽菜",
		"小坂菜緒",
		"富田鈴花",
		"丹生明里",
		"濱岸ひより",
		"松田好花",
		"宮田愛萌",
		"渡邉美穂",
		"上村ひなの",
		"髙橋未来虹",
		"森本茉莉",
		"山口陽世",
		"石塚瑶季",
		"岸帆夏",
		"小西夏菜実",
		"清水理央",
		"正源司陽子",
		"竹内希来里",
		"平尾帆夏",
		"平岡海月",
		"藤嶌果歩",
		"宮地すみれ",
		"山下葉留花",
		"渡辺莉奈",
	}
)

// メンバーとメンバーid.
var MemberToIdMap = map[string]string{
	"潮紗理菜":  "2",
	"影山優佳":  "4",
	"加藤史帆":  "5",
	"齊藤京子":  "6",
	"佐々木久美": "7",
	"佐々木美玲": "8",
	"高瀬愛奈":  "9",
	"高本彩花":  "10",
	"東村芽依":  "11",
	"金村美玖":  "12",
	"河田陽菜":  "13",
	"小坂菜緒":  "14",
	"富田鈴花":  "15",
	"丹生明里":  "16",
	"濱岸ひより": "17",
	"松田好花":  "18",
	"宮田愛萌":  "19",
	"渡邉美穂":  "20",
	"上村ひなの": "21",
	"髙橋未来虹": "22",
	"森本茉莉":  "23",
	"山口陽世":  "24",
	"石塚瑶季":  "25",
	"岸帆夏":   "26",
	"小西夏菜実": "27",
	"清水理央":  "28",
	"正源司陽子": "29",
	"竹内希来里": "30",
	"平尾帆夏":  "31",
	"平岡海月":  "32",
	"藤嶌果歩":  "33",
	"宮地すみれ": "34",
	"山下葉留花": "35",
	"渡辺莉奈":  "36",
	"ポカ":    "000",
}

func main() {
	var name string
	fmt.Scan(&name)                                                                              //任意のメンバーを入力
	url := "https://www.hinatazaka46.com/s/official/artist/" + MemberToIdMap[name] + "?ima=0000" // 任意のメンバーのURL

	// HTTPリクエストを送信してHTMLを取得
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// goqueryを使ってHTMLを解析
	article, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc := article.Find(".l-contents")

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
