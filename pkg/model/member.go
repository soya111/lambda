package model

import "fmt"

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
		"ポカ",
		"四期生リレー",
	}
)

var MemberToNumberMap = map[string]string{
	"潮紗理菜":   "2",
	"影山優佳":   "4",
	"加藤史帆":   "5",
	"齊藤京子":   "6",
	"佐々木久美":  "7",
	"佐々木美玲":  "8",
	"高瀬愛奈":   "9",
	"高本彩花":   "10",
	"東村芽依":   "11",
	"金村美玖":   "12",
	"河田陽菜":   "13",
	"小坂菜緒":   "14",
	"富田鈴花":   "15",
	"丹生明里":   "16",
	"濱岸ひより":  "17",
	"松田好花":   "18",
	"宮田愛萌":   "19",
	"渡邉美穂":   "20",
	"上村ひなの":  "21",
	"髙橋未来虹":  "22",
	"森本茉莉":   "23",
	"山口陽世":   "24",
	"ポカ":     "000",
	"四期生リレー": "2000",
}

func IsMember(text string) bool {
	_, exists := MemberToNumberMap[text]
	return exists
}

func GetMemberNumber(memberName string) (string, error) {
	number, exists := MemberToNumberMap[memberName]
	if !exists {
		return "", fmt.Errorf("member not found: %s", memberName)
	}
	return number, nil
}
