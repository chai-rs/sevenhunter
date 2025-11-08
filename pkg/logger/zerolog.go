package logx

import "github.com/rs/zerolog"

func ConsoleWriter() zerolog.Logger {
	return zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
}
