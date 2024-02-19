package common

import "fmt"

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
