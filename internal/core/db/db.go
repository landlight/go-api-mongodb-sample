package db

import (
	"context"
	"go-api-mongodb-sample/internal/core/config"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

//---------------------------------------------------------

// Key to use when setting the db session.
type ctxDatabaseKey int

// RequestIDKey is the key that holds th unique request ID in a request context.
const DBSessionKey ctxDatabaseKey = 0

var (
	filesCollection = "files"
)

type Collection struct {
	Name    string
	Indices []mgo.Index
}

//---------------------------------------------------------

type DBConn struct {
	DialInfo *mgo.DialInfo
}

func (db *DBConn) DialDB() *DBSession {
	if db.DialInfo == nil {
		return &DBSession{}
	}

	session, err := mgo.DialWithInfo(db.DialInfo)
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)

	return &DBSession{Session: session}
}

//---------------------------------------------------------

type DBSession struct {
	Session *mgo.Session
}

func (s *DBSession) Copy() *DBSession {
	if s.Session == nil {
		return s
	}

	return &DBSession{Session: s.Session.Copy()}
}

func (s *DBSession) Close() {
	if s.Session == nil {
		return
	}
	s.Session.Close()
}

func (s *DBSession) DB(name string) *DBInstance {
	if s.Session == nil {
		return &DBInstance{}
	}

	return &DBInstance{Database: s.Session.DB(name)}
}

//---------------------------------------------------------

type DBInstance struct {
	Database *mgo.Database
}

//---------------------------------------------------------

// DialMongoMiddleware -> This function create handler that connect to mongo and put into context.
// Note: it is not a best practice to put database session into context but the main reason is
// only to automatic open/close database connection without developer to remember to do by themselves.
func DialMongoMiddleware(dbConn *DBConn, dbName string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		session := dbConn.DialDB()

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s := session.Copy()
			defer s.Close()

			// create context and add mgo session into it
			ctx := r.Context()

			// get database session then add to context
			ctx = context.WithValue(ctx, DBSessionKey, s.DB(dbName))

			// put copied session into new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getMongoSessionFromContext(ctx context.Context) (*DBInstance, error) {
	if db, ok := ctx.Value(DBSessionKey).(*DBInstance); ok {
		return db, nil
	}
	logrus.Errorf("[getMongoSessionFromContext] get mongo session with DBSessionKey=%d error: %s", DBSessionKey, config.RR.Internal.DBSessionNotFound.Error())
	return nil, config.RR.Internal.DBSessionNotFound
}

func getDBInstance(ctx context.Context) (*mgo.Database, error) {
	dbInstance, err := getMongoSessionFromContext(ctx)
	if err != nil {
		logrus.Errorf("[getDBInstance] get database instance error: %s", err)
		return nil, config.RR.Internal.DBSessionNotFound
	}
	return dbInstance.Database, nil
}

func GetCollectionFiles(ctx context.Context) (*mgo.Collection, error) {
	db, err := getDBInstance(ctx)
	if err != nil {
		logrus.Errorf("[GetCollectionFiles] get database error: %s", err)
		return nil, err
	}

	return db.C(filesCollection), nil
}

func Indices() []Collection {
	collections := []Collection{
		{
			Name: filesCollection,
			Indices: []mgo.Index{
				{
					Key:        []string{"bucket"},
					Background: true,
				},
			},
		},
	}
	return collections
}
