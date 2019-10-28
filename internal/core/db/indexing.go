package db

import (
	"go-api-mongodb-sample/internal/core/config"
	"github.com/sirupsen/logrus"
)

func EnsureIndex(dbConn *DBConn) {
	session := dbConn.DialDB()
	defer session.Close()

	dbInstance := session.DB(config.CF.MongoDB.Schema.DBName)
	ensure(dbInstance, Indices())
}

func ensure(db *DBInstance, collections []Collection) {
	for _, col := range collections {
		for _, index := range col.Indices {
			if err := db.Database.C(col.Name).EnsureIndex(index); err != nil {
				//when you fix index and then different options
				logrus.Warnf("cannot create indext %+v for collection %s, err: %+v", index.Key, col.Name, err)
				logrus.Infof("drop Index: %+v for collection: %s", index.Key, col.Name)
				if err := db.Database.C(col.Name).DropIndexName(index.Name); err != nil {
					logrus.Errorf("cannot DropIndex %+v for collection %s, err: %+v", index.Key, col.Name, err)
				}

				if err := db.Database.C(col.Name).EnsureIndex(index); err != nil {
					// Panic because some function and logic code MUST depend on index, if not,
					// the application will be failed or run unexpected
					logrus.Panicf("cannot create indext %+v for collection %s, err: %+v", index.Key, col.Name, err)
				}
				logrus.Infof("try again EnsureIndex: %s", index.Name)
			}
		}
	}
}
