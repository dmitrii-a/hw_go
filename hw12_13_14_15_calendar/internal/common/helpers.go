package common

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"
)

// IsErr helper function to return bool if err != nil.
func IsErr(err error) bool {
	return err != nil
}

// ConnectionDBString return connection string to DB.
func ConnectionDBString(c DBConfig) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode)
}

func GetServerAddr(host string, port int) string {
	return fmt.Sprintf("%v:%v", host, port)
}

func GetConfigPathFromArg() string {
	var path string
	flag.StringVar(&path, "config", "", "Path to configuration file")
	flag.Parse()
	return path
}

func GetNotifyCancelCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(
		context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)
	return ctx, cancel
}
