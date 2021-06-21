package saga

import "services.shared/saga/msg"

type CommandHandler func(command msg.Command) error
