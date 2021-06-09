package logger

type LogLevel int

const (
	TraceLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelText = map[LogLevel]string{
	TraceLevel: "TRACE",
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

func LevelText(level LogLevel) string {
	return levelText[level]
}

type Logger interface {
	// Ping pings log server and return error if got problem
	Ping() error

	// SetLevel specifies base level to log
	SetLevel(level LogLevel)

	// Trace is a code smell if used in production, this should be used during development
	// to track bugs, but never committed to VCS
	Trace(args ...interface{})
	// Debug is about anything that happens in the program. This is mostly used during debugging
	// should be trimed down the number of debug statement before entering the production stage,
	// so that only the most meaningful entries are left, and can be activated during troubleshooting
	Debug(args ...interface{})
	// Info is about all actions that are user-driven, or system specific (ie regularly scheduled operations…)
	Info(args ...interface{})
	// Warn all events that could potentially become an error. For instance if one database call took more
	// than a predefined time, or if an in-memory cache is near capacity. This will allow proper automated
	// alerting, and during troubleshooting will allow to better understand how the system was behaving before the failure.
	Warn(args ...interface{})
	// Error is about error condition. That can be API calls that return errors or consolelogger error conditions.
	Error(args ...interface{})
	// Fatal means it's too bad, it’s doomsday. Use this very scarcely, this shouldn’t happen a lot.
	// Usually logging at this level signifies the end of the program. For instance, if a network daemon
	// can’t bind a network socket, log at this level and exit is the only sensible thing to do.
	Fatal(args ...interface{})
}
