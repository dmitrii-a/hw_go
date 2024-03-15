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
	Host              string `mapstructure:"HOST"`
	Port              int    `mapstructure:"PORT"`
	GrpcHost          string `mapstructure:"GRPC_HOST"`
	GrpcPort          int    `mapstructure:"GRPC_PORT"`
	GrpcGWHost        string `mapstructure:"GRPC_GW_HOST"`
	GrpcGWPort        int    `mapstructure:"GRPC_GW_PORT"`
	Debug             bool   `mapstructure:"DEBUG"`
	LogLevel          string `mapstructure:"LOG_LEVEL"`
	ShutdownTimeout   int    `mapstructure:"SHUTDOWN_TIMEOUT_SECOND"`
	ReadHeaderTimeout int    `mapstructure:"READ_HEADER_TIMEOUT_SECOND"`
	ReadTimeout       int    `mapstructure:"READ_TIMEOUT_SECOND"`
}

type SchedulerConfig struct {
	EventLifetime     int `mapstructure:"EVENT_LIFETIME_SECOND"`
	PublishPeriodTime int `mapstructure:"PUBLISH_PERIOD_TIME_SECOND"`
}

type RabbitConfig struct {
	Host     string `mapstructure:"HOST"`
	Port     int    `mapstructure:"PORT"`
	Username string `mapstructure:"USERNAME"`
	Password string `mapstructure:"PASSWORD"`
}

// AppConfig app config.
type AppConfig struct {
	Server     ServerConfig    `mapstructure:"APP"`
	Scheduler  SchedulerConfig `mapstructure:"SCHEDULER"`
	DB         DBConfig        `mapstructure:"DB"`
	RabbitMQ   RabbitConfig    `mapstructure:"RABBITMQ"`
	UseCacheDB bool            `mapstructure:"USE_CACHE_DB"`
}

// Config project config.
var Config AppConfig

func setDefaults() {
	viper.SetDefault("USE_CACHE_DB", false)

	viper.SetDefault("DB.USERNAME", "admin")
	viper.SetDefault("DB.PASSWORD", "password")
	viper.SetDefault("DB.DATABASE", "calendar-service")
	viper.SetDefault("DB.HOST", "127.0.0.1")
	viper.SetDefault("DB.PORT", 5455)
	viper.SetDefault("DB.SSL_MODE", "disable")

	viper.SetDefault("APP.HOST", "127.0.0.1")
	viper.SetDefault("APP.PORT", 8080)
	viper.SetDefault("APP.GRPC_HOST", "127.0.0.1")
	viper.SetDefault("APP.GRPC_PORT", 50051)
	viper.SetDefault("APP.GRPC_GW_HOST", "127.0.0.1")
	viper.SetDefault("APP.GRPC_GW_PORT", 3000)
	viper.SetDefault("APP.DEBUG", true)
	viper.SetDefault("APP.LOG_LEVEL", "info")
	viper.SetDefault("APP.SHUTDOWN_TIMEOUT_SECOND", 30)
	viper.SetDefault("APP.READ_HEADER_TIMEOUT_SECOND", 10)
	viper.SetDefault("APP.READ_TIMEOUT_SECOND", 10)

	viper.SetDefault("RABBITMQ.HOST", "127.0.0.1")
	viper.SetDefault("RABBITMQ.PORT", 5675)
	viper.SetDefault("RABBITMQ.USERNAME", "admin")
	viper.SetDefault("RABBITMQ.PASSWORD", "password")

	viper.SetDefault("SCHEDULER.EVENT_LIFETIME_SECOND", 60*60*24*365)
	viper.SetDefault("SCHEDULER.PUBLISH_PERIOD_TIME_SECOND", 10)
}

func init() {
	var err error
	log := logger.InitLogger(zerolog.GlobalLevel())
	setDefaults()
	if IsErr(err) {
		log.Fatal().Msgf("Unable to decode into struct, %v", err)
	}
	viper.AutomaticEnv()
	err = viper.Unmarshal(&Config)
	if IsErr(err) {
		log.Fatal().Msgf("Unable to decode into struct, %v", err)
	}
}

// SetConfigFileSettings applying settings from a configuration file.
func (config *AppConfig) SetConfigFileSettings(path string) {
	flag.Parse()
	log := logger.InitLogger(zerolog.GlobalLevel())
	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); IsErr(err) {
			log.Fatal().Msgf("Error reading config file, %s", err)
		}
		err := viper.Unmarshal(&Config)
		if IsErr(err) {
			log.Fatal().Msgf("Unable to decode into struct, %v", err)
		}
	}
	// Override config with environment variables if they exist.
	viper.AutomaticEnv()
	err := setLogger(Config.Server.LogLevel)
	if IsErr(err) {
		log.Fatal().Msgf("Error parsing loglevel, %s", err)
	}
}
