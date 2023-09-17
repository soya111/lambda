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

	// GradList is a list of graduated members.
	GradList = []string{
		"井口眞緒",
		"柿崎芽実",
		"影山優佳",
		"長濱ねる",
		"宮田愛萌",
		"渡邉美穂",
	}
)

var (
	// MemberToIdMap is a map of member name to member ID.
	MemberToIdMap = map[string]string{
		"潮紗理菜":  "2",
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

	// gradToIdMap is a map of graduated member name to member ID.
	gradToIdMap = map[string]string{
		"井口眞緒": "1",
		"柿崎芽実": "3",
		"影山優佳": "4",
		"長濱ねる": "",
		"宮田愛萌": "19",
		"渡邉美穂": "20",
	}
)

var (
	// MemberToGenerationMap is a map of member name to member generation.
	MemberToGenerationMap = map[string]string{
		"潮紗理菜":  "1",
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

	// GradToGenerationMap is a map of graduated member name to member generation.
	GradToGenerationMap = map[string]string{
		"井口眞緒": "1",
		"柿崎芽実": "1",
		"影山優佳": "1",
		"長濱ねる": "1",
		"宮田愛萌": "2",
		"渡邉美穂": "2",
	}
)

// MemberToNicknameMap is a map of member name to nickname.
var MemberToNicknameMap = map[string][]string{
	"潮紗理菜":  {"潮くん", "なっちょ", "サリマカシー", "うしし"},
	"加藤史帆":  {"かとし", "しし", "としちゃん", "天使"},
	"齊藤京子":  {"きょんこ", "きょうこにょう"},
	"佐々木久美": {"くみてん", "ささく", "きくちゃん", "キャプテン"},
	"佐々木美玲": {"みーぱん", "ささみ"},
	"高瀬愛奈":  {"まなふぃ", "まなふい"},
	"高本彩花":  {"たけもと", "おたけ", "あやちぇり", "あや"},
	"東村芽依":  {"めいめい", "めいちご", "やんちゃる", "ちゃる"},
	"金村美玖":  {"おすし", "ミクティー", "みーきゅん"},
	"河田陽菜":  {"かわだ", "かわださん", "おひな"},
	"小坂菜緒":  {"こさかな", "こしゃ"},
	"富田鈴花":  {"すーじー"},
	"丹生明里":  {"にぶちゃん", "タルタルチキン"},
	"濱岸ひより": {"ひよたん"},
	"松田好花":  {"このちゃん", "だーこの"},
	"上村ひなの": {"ひなのなの"},
	"髙橋未来虹": {"みくにん", "みくにちゃん"},
	"森本茉莉":  {"まりもと", "天才", "まりぃ", "あいつ"},
	"山口陽世":  {"ぱる", "はるよちゃん"},
	"石塚瑶季":  {"たまちゃん"},
	"岸帆夏":   {"岸君", "きしほの", "きしほ"},
	"小西夏菜実": {"こにしん", "524773"},
	"清水理央":  {"りおたむ", "ずりお"},
	"正源司陽子": {"げんちゃん", "しょげこ"},
	"竹内希来里": {"きらりんちょ", "きらりん"},
	"平尾帆夏":  {"ひらほー", "ひらほ"},
	"平岡海月":  {"みっちゃん", "くらげ"},
	"藤嶌果歩":  {"かほりん", "かほりんこうりん"},
	"宮地すみれ": {"すみレジェンド", "レジェ", "すみこ"},
	"山下葉留花": {"はるはる"},
	"渡辺莉奈":  {"りなし", "べりな"},
	"ポカ":    {},
}

var (
	ErrNonExistentMember = errors.New("日向坂46に存在しないメンバーです。")
	ErrGraduatedMember   = errors.New("日向坂46の卒業メンバーです。")
)

// NormalizeName normalizes a member name.
func NormalizeName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, " ", "")
	return name
}

// Create a reverse map for MemberToNicknameMap for faster lookup.
var nicknameToMemberMap = make(map[string]string)

func init() {
	for member, nicknames := range MemberToNicknameMap {
		for _, nickname := range nicknames {
			nicknameToMemberMap[nickname] = member
		}
	}
}

// TranslateNicknametoMember returns the member translated from nickname, or returns the argument if a nickname does not exist.
func TranslateNicknametoMember(nickname string) string {
	member, exists := nicknameToMemberMap[nickname]
	if !exists {
		return nickname
	}
	return member
}

// IsMember returns true if the given text is a member name.
func IsMember(text string) bool {
	_, exists := MemberToIdMap[text]
	return exists
}

// IsGrad returns true if the given text is a guraduated member name.
func IsGrad(text string) bool {
	_, exists := gradToIdMap[text]
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
