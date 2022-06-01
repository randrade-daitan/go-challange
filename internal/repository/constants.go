package repository

import "os"

const DBProtocol = "tcp"

func DBName() string {
	return os.Getenv("DBNAME")
}

func DBURL() string {
	return os.Getenv("DBURL")
}

func DBUser() string {
	return os.Getenv("DBUSER")
}

func DBPass() string {
	return os.Getenv("DBPASS")
}
