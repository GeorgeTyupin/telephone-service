//go:build integration

package tests

import (
	"context"
	"crypto/rand"
	"io/fs"
	"log/slog"
	"net/http/httptest"
	"time"

	api "telephone-service/codegen/oas"
	"telephone-service/internal/config"
	"telephone-service/internal/container"
	"telephone-service/internal/db/postgres/migrator"
	"telephone-service/internal/repository/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AcceptanceSuite struct {
	suite.Suite

	DatabaseURL string
	Migrations  fs.FS

	Ctx        context.Context
	DB         *pgxpool.Pool
	Repository *postgres.Repository
	APIClient  *api.Client

	adminDB      *pgxpool.Pool
	templateName string
}

func (s *AcceptanceSuite) SetupSuite() {
	t := s.T()
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()

	var err error
	s.adminDB, err = pgxpool.New(ctx, s.DatabaseURL)
	require.NoError(t, err)
	t.Cleanup(s.adminDB.Close)

	s.templateName = "template_" + rand.Text()
	_, err = s.adminDB.Exec(ctx, "CREATE DATABASE "+pgx.Identifier{s.templateName}.Sanitize())
	require.NoError(t, err)
	name := s.templateName
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := s.adminDB.Exec(cleanupCtx, "DROP DATABASE "+pgx.Identifier{name}.Sanitize()+" WITH (FORCE)")
		require.NoError(t, err)
	})

	cfg, err := pgxpool.ParseConfig(s.DatabaseURL)
	require.NoError(t, err)
	cfg.ConnConfig.Database = s.templateName
	templateDB, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	defer templateDB.Close()

	require.NoError(t, migrator.Run(ctx, templateDB, s.Migrations))
}

func (s *AcceptanceSuite) SetupTest() {
	t := s.T()
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	t.Cleanup(cancel)
	s.Ctx = ctx

	databaseName := "test_" + rand.Text()
	_, err := s.adminDB.Exec(ctx, "CREATE DATABASE "+pgx.Identifier{databaseName}.Sanitize()+
		" TEMPLATE "+pgx.Identifier{s.templateName}.Sanitize())
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanupCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_, err := s.adminDB.Exec(cleanupCtx, "DROP DATABASE "+pgx.Identifier{databaseName}.Sanitize()+" WITH (FORCE)")
		require.NoError(t, err)
	})

	cfg, err := pgxpool.ParseConfig(s.DatabaseURL)
	require.NoError(t, err)
	cfg.ConnConfig.Database = databaseName
	cfg.MaxConns = 4
	s.DB, err = pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(s.DB.Close)

	logger := slog.New(slog.NewTextHandler(t.Output(), nil))
	application := container.NewContainerWithDependencies(config.Config{
		Server: config.Server{MaxBodyBytes: 1 << 20},
	}, logger, s.DB)
	s.Repository = application.ResolvePostgresRepository()
	handler, err := application.BuildHandler()
	require.NoError(t, err)

	server := httptest.NewTestServer(t, handler)
	client := server.Client()
	client.Timeout = 10 * time.Second
	s.APIClient, err = api.NewClient(server.URL, api.WithClient(client))
	require.NoError(t, err)
}
