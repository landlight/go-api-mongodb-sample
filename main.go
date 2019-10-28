package main

import (
	"go-api-mongodb-sample/internal/core/config"
	"go-api-mongodb-sample/internal/core/db"
	"go-api-mongodb-sample/internal/handler/routes"
	"go-api-mongodb-sample/internal/handler/servers"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	if err := config.InitConfig("configs"); err != nil {
		panic(err)
	}
	logrus.Infof("Initital 'Configuration'. %+v", config.CF)

	if err := config.InitReturnResult("configs"); err != nil {
		panic(err)
	}
	logrus.Infof("Initital 'Return Result''. %+v", config.CF)

	logrus.SetFormatter(&logrus.TextFormatter{})

	dbcon := &db.DBConn{
		DialInfo: config.CF.MongoDB.DialInfo.DBName,
	}

	db.EnsureIndex(dbcon)

	r := routes.NewRouter(dbcon)
	srv := servers.NewServer(strconv.Itoa(config.CF.Port), r)
	srv.ListenAndServeWithGracefulShutdown()
}
