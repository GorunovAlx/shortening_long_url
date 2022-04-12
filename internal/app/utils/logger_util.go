package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
)

var Logger zerolog.Logger

func LoggerInit() {
	logfile, err := os.OpenFile("server_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal().Err(err).Msgf("Cannot start %v", logfile)
	}
	defer logfile.Close()

	Logger := zerolog.New(logfile).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.Level(configs.Cfg.ZerologLevel))

	Logger.Info().Msg("start server")
}
