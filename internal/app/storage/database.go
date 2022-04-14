package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	"github.com/GorunovAlx/shortening_long_url/internal/app/utils"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBStorage struct {
	dsn      string
	Postgres *pgxpool.Pool
}

func NewPGXPool(ctx context.Context, dsn string, logger pgx.Logger, logLevel pgx.LogLevel) (*pgxpool.Pool, error) {
	conf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	conf.ConnConfig.Logger = logger

	if logLevel != 0 {
		conf.ConnConfig.LogLevel = logLevel
	}

	conf.MaxConns = 20
	conf.MaxConnIdleTime = time.Second * 30
	conf.MaxConnLifetime = time.Minute * 2

	pool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, fmt.Errorf("pgx connection error: %w", err)
	}
	return pool, nil
}

// LogLevelFromEnv returns the pgx.LogLevel from the environment variable PGX_LOG_LEVEL.
// By default this is info (pgx.LogLevelInfo), which is good for development.
func LogLevelFromEnv() (pgx.LogLevel, error) {
	if level := configs.Cfg.PgxLogLevel; level != "" {
		l, err := pgx.LogLevelFromString(level)
		if err != nil {
			return pgx.LogLevelDebug, fmt.Errorf("pgx configuration: %w", err)
		}
		return l, nil
	}
	return pgx.LogLevelInfo, nil
}

// PGXStdLogger prints pgx logs to the standard logger.
// os.Stderr by default.
type PGXStdLogger struct{}

func (l *PGXStdLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	args := make([]interface{}, 0, len(data)+2) // making space for arguments + level + msg
	args = append(args, level, msg)
	for k, v := range data {
		args = append(args, fmt.Sprintf("%s=%v", k, v))
	}
	log.Println(args...)
}

func NewDBStorage() (*DBStorage, error) {
	pgxLogLevel, err := LogLevelFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	pgPool, err := NewPGXPool(context.Background(), configs.Cfg.DatabaseDSN, &PGXStdLogger{}, pgxLogLevel)
	if err != nil {
		log.Fatal(err)
	}

	storage := &DBStorage{
		dsn:      configs.Cfg.DatabaseDSN,
		Postgres: pgPool,
	}

	if err := storage.Init(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (dbs *DBStorage) Init() error {
	if err := dbs.CreateTable(); err != nil {
		return err
	}
	return nil
}

func (dbs *DBStorage) GetInitialLink(shortLink string) (string, error) {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return "", e
	}
	defer conn.Release()

	var iLink string
	var deleted bool
	err := conn.QueryRow(
		context.Background(),
		"select initial_link, COALESCE(deleted, false) from shortened_links where short_link=$1",
		shortLink,
	).Scan(&iLink, &deleted)
	if err != nil {
		return "", err
	}

	if deleted {
		err = utils.NewDeletedLinkError(shortLink)
		return "", err
	}

	return iLink, nil
}

func (dbs *DBStorage) WriteShortURL(shortURL *ShortURL) error {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return e
	}
	defer conn.Release()

	insertStatement := `
	INSERT INTO shortened_links (initial_link, short_link, user_id, date_of_create)
	VALUES ($1, $2, $3, $4) ON CONFLICT (initial_link) DO NOTHING;`

	commandTag, err := conn.Exec(
		context.Background(),
		insertStatement,
		shortURL.InitialLink,
		shortURL.ShortLink,
		shortURL.UserID,
		time.Now(),
	)

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		err = utils.NewInsertUniqueLinkError(shortURL.InitialLink)
		return err
	}
	return nil
}

func (dbs *DBStorage) GetAllShortURLByUser(userID uint32) ([]ShortURLByUser, error) {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return nil, e
	}
	defer conn.Release()

	var result []ShortURLByUser

	selectStatement := "select initial_link, short_link from shortened_links where user_id=$1"
	rows, err := conn.Query(context.Background(), selectStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s ShortURLByUser
		err = rows.Scan(&s.InitialLink, &s.ShortLink)
		if err != nil {
			return nil, err
		}
		s.ShortLink = configs.Cfg.BaseURL + "/" + s.ShortLink
		result = append(result, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (dbs *DBStorage) PingDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.Connect(ctx, dbs.dsn)
	if err != nil {
		return err
	}
	defer dbPool.Close()

	if dbPool != nil {
		return nil
	}

	return errors.New("ping attempt failed")
}

func (dbs *DBStorage) WriteListShortURL(links []ShortURLByUser) error {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return e
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for _, l := range links {
		if _, err = tx.Conn().Exec(
			context.Background(),
			"INSERT INTO shortened_links (initial_link, short_link, user_id, date_of_create) VALUES ($1, $2, $3, $4)",
			l.InitialLink,
			l.ShortLink,
			nil,
			time.Now(),
		); err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (dbs *DBStorage) CheckURLsCreatedByUser(links []string, id uint32) ([]string, error) {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return nil, e
	}
	defer conn.Release()

	var result []string

	selectStatement := `with temp as (
		select short_link 
		from public.shortened_links
		where user_id = $1)
		
		select links.short_link from
		(select unnest(ARRAY[$2::varchar[]]) short_link) links
		except
		select short_link
		from temp`
	rows, err := conn.Query(context.Background(), selectStatement, id, links)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			return nil, err
		}

		result = append(result, s)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (dbs *DBStorage) DeleteShortURLByUser(link string, id uint32) error {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return e
	}
	defer conn.Release()

	sqlStmt := `
	update shortened_links set deleted = true 
	where user_id = $1 and short_link = $2;`

	_, err := conn.Exec(
		context.Background(),
		sqlStmt,
		id,
		link,
	)

	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) CreateTable() error {
	conn, e := dbs.Postgres.Acquire(context.Background())
	if e != nil {
		return e
	}
	defer conn.Release()

	sqlCreateStmt := `
	create table if not exists public.shortened_links ( id bigserial constraint shortened_link_pk primary key,
	initial_link varchar(256) not null unique, short_link varchar(256) not null, user_id bigint,
	date_of_create date, deleted boolean ); alter table public.shortened_links owner to postgres;`

	_, err := conn.Exec(context.Background(), sqlCreateStmt)
	if err != nil {
		return err
	}

	return nil
}
