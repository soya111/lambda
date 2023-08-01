package webhook

var (
	memberList = []string{
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

func isMember(text string) bool {
	for _, v := range memberList {
		if text == v {
			return true
		}
	}
	return false
}
