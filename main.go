package main

import (
	"./routes"
)

func main() {
	// This repository is based on this article https://tutorialedge.net/golang/creating-restful-api-with-golang/
	//TODO save articles in Postgres https://www.enterprisedb.com/postgres-tutorials/postgresql-and-golang-tutorial
	//article about go and postgres https://golangdocs.com/golang-postgresql-example
	routes.HandleRequests()
}
