package common

import (
	"flag"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// DBConfig database config.
type DBConfig struct {
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
	Database string `mapstructure:"DATABASE"`
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	SSLMode  string `mapstructure:"SSL_MODE"`
}

// ServerConfig server config.
type ServerConfig struct {
	Host            string `mapstructure:"HOST"`
	Port            int    `mapstructure:"PORT"`
	Debug           bool   `mapstructure:"DEBUG"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	ShutdownTimeout int    `mapstructure:"SHUTDOWN_TIMEOUT_SECOND"`
}

// AppConfig app config.
type AppConfig struct {
	Server ServerConfig `mapstructure:"APP"`
	DB     DBConfig     `mapstructure:"DB"`
}

// Config project config.
var Config AppConfig

func setDefaults() {
	viper.SetDefault("DB.USERNAME", "admin")
	viper.SetDefault("DB.PASSWORD", "password")
	viper.SetDefault("DB.DATABASE", "calendar-service")
	viper.SetDefault("DB.HOST", "127.0.0.1")
	viper.SetDefault("DB.PORT", 5455)
	viper.SetDefault("DB.SSL_MODE", "disable")

	// Set default values for the ServerConfig.
	viper.SetDefault("APP.HOST", "127.0.0.1")
	viper.SetDefault("APP.PORT", 8080)
	viper.SetDefault("APP.DEBUG", true)
	viper.SetDefault("APP.LOG_LEVEL", "info")
	viper.SetDefault("APP.SHUTDOWN_TIMEOUT_SECOND", 30)
}

func init() {
	var (
		path string
		err  error
	)
	flag.StringVar(&path, "config", "", "Path to configuration file")
	flag.Parse()
	log := logger.InitLogger(zerolog.GlobalLevel())
	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); IsErr(err) {
			log.Fatal().Msgf("Error reading config file, %s", err)
		}
		err = viper.Unmarshal(&Config)
	} else {
		setDefaults()
	}
	if IsErr(err) {
		log.Fatal().Msgf("Unable to decode into struct, %v", err)
	}
	// Override config with environment variables if they exist.
	viper.AutomaticEnv()
	err = viper.Unmarshal(&Config)
	if IsErr(err) {
		log.Fatal().Msgf("Unable to decode into struct, %v", err)
	}
}
