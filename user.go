package watchub

type UserInfo struct {
	Login string
	Email string
}

type UsersSvc interface {
	Info(token string) (UserInfo, error)
}
