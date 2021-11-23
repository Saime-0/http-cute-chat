package its

type Someone uint8

const (
	Owner Someone = iota
	Admin
	Moder
	Member
)

func List(somes ...Someone) []Someone {
	return somes
}
