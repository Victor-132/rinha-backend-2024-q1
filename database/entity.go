package database

import "time"

type Client struct {
	Id      int
	Limit   int
	Balance int
}

type Transaction struct {
	ClientId    int
	Value       int
	Type        string
	Description string
	Date        time.Time
}
