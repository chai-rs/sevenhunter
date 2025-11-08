package logx

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(debug ...bool) {
	if len(debug) == 0 || !debug[0] {
		log.Logger = log.Logger.Level(zerolog.InfoLevel).With().Timestamp().Caller().Logger()
	} else {
		log.Logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Caller().Logger()
		log.Logger = log.Logger.Level(zerolog.DebugLevel)
	}
}

func Debug() *zerolog.Event {
	return log.Debug()
}

func Info() *zerolog.Event {
	return log.Info()
}

func Warn() *zerolog.Event {
	return log.Warn()
}

func Error() *zerolog.Event {
	return log.Error()
}

func Panic() *zerolog.Event {
	return log.Panic()
}

func Fatal() *zerolog.Event {
	return log.Fatal()
}
