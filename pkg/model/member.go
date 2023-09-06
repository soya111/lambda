package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// MemberList is a list of all members.
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

// MemberToIdMap is a map of member name to member ID.
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

// MemberToGenerationMap is a map of member name to member generation.
var MemberToGenerationMap = map[string]string{
	"潮紗理菜":  "1",
	"影山優佳":  "1",
	"加藤史帆":  "1",
	"齊藤京子":  "1",
	"佐々木久美": "1",
	"佐々木美玲": "1",
	"高瀬愛奈":  "1",
	"高本彩花":  "1",
	"東村芽依":  "1",
	"金村美玖":  "2",
	"河田陽菜":  "2",
	"小坂菜緒":  "2",
	"富田鈴花":  "2",
	"丹生明里":  "2",
	"濱岸ひより": "2",
	"松田好花":  "2",
	"宮田愛萌":  "2",
	"渡邉美穂":  "2",
	"上村ひなの": "3",
	"髙橋未来虹": "3",
	"森本茉莉":  "3",
	"山口陽世":  "3",
	"石塚瑶季":  "4",
	"岸帆夏":   "4",
	"小西夏菜実": "4",
	"清水理央":  "4",
	"正源司陽子": "4",
	"竹内希来里": "4",
	"平尾帆夏":  "4",
	"平岡海月":  "4",
	"藤嶌果歩":  "4",
	"宮地すみれ": "4",
	"山下葉留花": "4",
	"渡辺莉奈":  "4",
	"ポカ":    "?",
}

var ErrNonExistentMember = errors.New("日向坂46に存在しないメンバーです。")

// NormalizeName normalizes a member name.
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "")
	return name
}

// IsMember returns true if the given text is a member name.
func IsMember(text string) bool {
	_, exists := MemberToIdMap[text]
	return exists
}

// GetMemberId returns the member ID for the given member name.
func GetMemberId(memberName string) (string, error) {
	memberName = NormalizeName(memberName)
	number, exists := MemberToIdMap[memberName]
	if !exists {
		return "", fmt.Errorf("member not found: %s", memberName)
	}
	return number, nil
}
