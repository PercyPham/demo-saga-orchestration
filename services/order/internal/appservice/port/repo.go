package port

type Repo interface {
	Ping() error
}
