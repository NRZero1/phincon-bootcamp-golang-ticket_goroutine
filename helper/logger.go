package helper

import "github.com/rs/zerolog/log"

func Info(message string) {
	log.Info().Msg(message)
}

func Debug(message string) {
	log.Debug().Msg(message)
}

func Trace(message string) {
	log.Trace().Msg(message)
}

func Error(message string) {
	log.Error().Msg(message)
}

func Fatal(message string) {
	log.Fatal().Msg(message)
}