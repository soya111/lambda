package model

import (
	"fmt"
	"strings"
)

var (
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
		"ポカ",
	}
)

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

func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "")
	return name
}

func IsMember(text string) bool {
	_, exists := MemberToIdMap[text]
	return exists
}

func GetMemberId(memberName string) (string, error) {
	memberName = NormalizeName(memberName)
	number, exists := MemberToIdMap[memberName]
	if !exists {
		return "", fmt.Errorf("member not found: %s", memberName)
	}
	return number, nil
}
