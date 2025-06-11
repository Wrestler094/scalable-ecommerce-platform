package healthcheck

type Manager interface {
	SetAlive(bool)
	SetReady(bool)
	IsAlive() bool
	IsReady() bool
}
