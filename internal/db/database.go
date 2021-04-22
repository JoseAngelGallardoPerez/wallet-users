package db

import (
	"fmt"

	"github.com/Confialink/wallet-pkg-env_config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// CreateConnection creates a new db connection
func CreateConnection(dbconf *env_config.Db) (*gorm.DB, error) {

	// Get database details from environment variables
	host := dbconf.Host
	port := dbconf.Port
	user := dbconf.User
	DBName := dbconf.Schema
	password := dbconf.Password

	// initialize a new db connection
	return gorm.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True", // username:password@protocol(host:port)/dbname?param=value
			user, password, host, port, DBName,
		),
	)
}
