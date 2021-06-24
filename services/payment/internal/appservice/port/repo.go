package port

import "services.shared/saga"

type Repo interface {
	Ping() error

	saga.Repo
}
