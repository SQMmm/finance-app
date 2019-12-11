package container

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"github.com/sqmmm/finance-app/internal/config"
	"net"
	"strconv"
	"time"
)

func buildMySQLClient(cfg *config.Config) (write *sql.DB, read *sql.DB, err error) {
	if write, err = buildMysqlClientConn(&cfg.MySQL.Write, true); err != nil {
		return nil, nil, err
	}
	if read, err = buildMysqlClientConn(&cfg.MySQL.Read, false); err != nil {
		return nil, nil, err
	}

	return write, read, nil
}

func buildMysqlClientConn(cfg *config.MysqlConnect, onlyWrite bool) (*sql.DB, error) {
	if cfg.Port == 0 {
		cfg.Port = 3306
	}
	mcfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		User:                 cfg.Username,
		Passwd:               cfg.Password,
		DBName:               cfg.Database,
		RejectReadOnly:       onlyWrite,
		ParseTime:            true,
		AllowNativePasswords: true,
		Collation:            "utf8mb4_unicode_ci",
	}
	conn, err := connectCfg(mcfg)
	if err != nil {
		return nil, errors.Wrap(err, "mysql connection initialization error")
	}
	conn.SetConnMaxLifetime(10 * time.Second)

	return conn, nil
}

// ConnectCfg establish new connection to mysql database using config
func connectCfg(cfg *mysql.Config) (*sql.DB, error) {
	if cfg.Params == nil {
		cfg.Params = make(map[string]string)
	}
	// Setting default parameters
	if _, ok := cfg.Params["parseTime"]; !ok {
		cfg.Params["parseTime"] = "true"
	}
	if cfg.Net == "" {
		cfg.Net = "tcp"
	}
	if cfg.Loc == nil {
		cfg.Loc = time.Local
	}

	db, err := sql.Open("mysql", (*mysql.Config)(cfg).FormatDSN())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
