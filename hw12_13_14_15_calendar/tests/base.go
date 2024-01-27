package tests

import (
	"context"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/common"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// BaseDBTestSuite is a base tests suite for tests that use a database.
type BaseDBTestSuite struct {
	suite.Suite
	PGContainer *postgres.PostgresContainer
	DB          *sqlx.DB
}

func (s *BaseDBTestSuite) SetupSuite() {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.4"),
		postgres.WithDatabase(common.Config.DB.Database),
		postgres.WithUsername(common.Config.DB.Username),
		postgres.WithPassword(common.Config.DB.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if common.IsErr(err) {
		panic("Failed to start container: " + err.Error())
	}
	s.PGContainer = pgContainer

	port, err := s.PGContainer.MappedPort(ctx, "5432/tcp")
	if common.IsErr(err) {
		panic("Failed to get mapped port: " + err.Error())
	}
	common.Config.DB.Port = port.Int()
	host, err := s.PGContainer.Container.Host(ctx)
	if common.IsErr(err) {
		panic("Failed to get mapped host: " + err.Error())
	}
	common.Config.DB.Host = host

	dbURL := common.ConnectionDBString(common.Config.DB)
	s.DB, err = sqlx.Open("postgres", dbURL)
	if common.IsErr(err) {
		panic("Failed to connect to the database: " + err.Error())
	}
	_, b, _, _ := runtime.Caller(0) //nolint:dogsled
	basePath := filepath.Dir(b)

	if err := goose.Up(s.DB.DB, filepath.Join(basePath, "../migrations")); common.IsErr(err) {
		panic("Failed to apply migrations: " + err.Error())
	}
}

func (s *BaseDBTestSuite) TearDownSuite() {
	if err := s.PGContainer.Terminate(context.Background()); common.IsErr(err) {
		panic("Error terminating container: " + err.Error())
	}
	if s.DB != nil {
		if err := s.DB.Close(); common.IsErr(err) {
			panic("Error closing the database connection: " + err.Error())
		}
	}
}
