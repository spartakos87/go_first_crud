package main

import (
	"FirstRestGo/routes"
)


func main() {
	//TODO save articles in Postgres https://www.enterprisedb.com/postgres-tutorials/postgresql-and-golang-tutorial
	routes.HandleRequests()
}
