//nolint:revive
package repository

import (
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/domain"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/pkg/freecache"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Use singleton pattern for DB connection.
var (
	db      *sqlx.DB
	cacheDB *freecache.CacheDB
)

func init() {
	var err error
	if !common.Config.UseCacheDB {
		db, err = sqlx.Open("postgres", common.ConnectionDBString(common.Config.DB))
		if common.IsErr(err) {
			common.Logger.Fatal().Err(err)
		}
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
	} else {
		cacheDB = freecache.NewCacheDB(1024 * 1024 * 100)
	}
}

func GetEventRepository() domain.EventRepository {
	var eventRepository domain.EventRepository
	if common.Config.UseCacheDB {
		eventRepository = NewEventCacheRepository()
	} else {
		eventRepository = NewEventDBRepository()
	}
	return eventRepository
}
