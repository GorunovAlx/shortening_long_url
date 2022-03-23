package storage

import (
	"context"
	"errors"
	"time"

	"github.com/GorunovAlx/shortening_long_url/internal/app/configs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DBStorage struct {
	lastInsertId int
	dsn          string
}

func NewDBStorage() (*DBStorage, error) {
	storage := &DBStorage{
		dsn:          configs.Cfg.DatabaseDSN,
		lastInsertId: 0,
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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dbs.connectDB()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var iLink string
	e := conn.QueryRow(
		ctx,
		"select initial_link from shortened_links where short_link=$1",
		shortLink,
	).Scan(&iLink)
	if e != nil {
		return "", e
	}
	defer conn.Close()

	return iLink, nil
}

func (dbs *DBStorage) WriteShortURL(shortURL *ShortURL) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	conn, err := dbs.connectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	var iLink string
	e := conn.QueryRow(
		ctx,
		"select initial_link from shortened_links where initial_link=$1",
		shortURL.InitialLink,
	).Scan(&iLink)
	if e != pgx.ErrNoRows {
		return e
	}
	defer conn.Close()

	if iLink == shortURL.InitialLink {
		return nil
	}

	insertStatement := `
	INSERT INTO shortened_links (id, initial_link, short_link, user_id, date_of_create)
	VALUES ($1, $2, $3, $4, $5)`

	dbs.lastInsertId += 1
	commandTag, err := conn.Exec(
		context.Background(),
		insertStatement,
		dbs.lastInsertId,
		shortURL.InitialLink,
		shortURL.ShortLink,
		shortURL.UserID,
		time.Now(),
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no row inserted")
	}
	return nil
}

func (dbs *DBStorage) GetAllShortURLByUser(userID uint32) ([]ShortURLByUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dbs.connectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var result []ShortURLByUser

	selectStatement := "select initial_link, short_link from shortened_links where user_id=$1"
	rows, err := conn.Query(ctx, selectStatement, userID)
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

func (dbs *DBStorage) CreateTable() error {
	createTableSQL := "create table if not exists public.shortened_links" +
		"( id integer not null constraint shortened_link_pk primary key, initial_link varchar(256) not null," +
		"short_link varchar(256) not null, user_id bigint not null, date_of_create date not null); alter table public.shortened_links owner to postgres;"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := dbs.connectDB()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec(ctx, createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DBStorage) connectDB() (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbs.dsn)
	if err != nil {
		return nil, err
	}
	//config.Logger = log15adapter.NewLogger(log.New("module", "pgx"))

	conn, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
