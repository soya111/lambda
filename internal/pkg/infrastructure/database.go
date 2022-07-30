package infrastructure

type Database interface {
	GetDestination(memberName string) ([]string, error)
}
