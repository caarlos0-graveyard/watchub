package watchub

type UserInfo struct {
	Login string
	Email string
}

type UserSvc interface {
	Info() (UserInfo, error)
}
