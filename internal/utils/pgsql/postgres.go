package pgsql

import "fmt"

func BuildDns(host, port, sslMode, user, password, name string) string {
	return fmt.Sprintf(
		"host=%s port=%s sslmode=%s user=%s password=%s dbname=%s",
		host, port, sslMode, user, password, name,
	)
}
