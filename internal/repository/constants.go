package repository

import "os"

const DBName = "challange"
const DBURL = "127.0.0.1:3306"
const DBProtocol = "tcp"

func DBUser() string {
	return os.Getenv("DBUSER")
}

func DBPass() string {
	return os.Getenv("DBPASS")
}
