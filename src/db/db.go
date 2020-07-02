package db

import (
	"community-cloud/config"
	"community-cloud/logging"
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func (h *GroupHandler) getMaster() (master *sql.DB) {
	h.DbMaster.rwlock.RLock()
	if h.DbMaster.Status {
		master = h.DbMaster.Db
	} else {
		master = nil
	}
	h.DbMaster.rwlock.RUnlock()
	return
}

func (h *GroupHandler) getSlave() (slave *sql.DB) {
	h.DbSlave.rwlock.RLock()
	if h.DbSlave.Status {
		slave = h.DbSlave.Db
	} else {
		slave = nil
	}
	h.DbSlave.rwlock.RUnlock()
	return
}

func pingDBWithCancel(db *sql.DB) error {
	ctx, cancel := context.WithCancel(context.Background())
	err := db.PingContext(ctx)
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	return err
}

func pingDB(db *sql.DB) error {
	return db.Ping()
}

func connDB(handler *Handler) {
	handler.rwlock.Lock()
	if handler.Conf != "" {
		db, err := sql.Open("mysql", handler.Conf)
		logging.Logger.Infof("db.connDb Open %s %s: %v", handler.Name, handler.Conf, err)
		if err == nil {
			if handler.MaxOpenConn > 0 {
				db.SetMaxOpenConns(handler.MaxOpenConn)
			} else {
				db.SetMaxOpenConns(20)
			}
			if handler.MaxIdleConn > 0 {
				db.SetMaxIdleConns(handler.MaxIdleConn)
			} else {
				db.SetMaxIdleConns(10)
			}
			if handler.MaxLifetime > 0 {
				db.SetConnMaxLifetime(time.Duration(handler.MaxLifetime) * time.Second)
			} else {
				db.SetConnMaxLifetime(60 * 60 * time.Second)
			}
			handler.Db = db
			err := pingDB(db)
			logging.Logger.Infof("db.connDb Ping %s: %v", handler.Name, err)
			if err == nil {
				handler.Status = true
			}
		}
	}
	handler.rwlock.Unlock()
}

func initGroupDB(dbconf config.DatabaseConf) *GroupHandler {

	ghandler := &GroupHandler{
		DbMaster: &Handler{Name: "master",
			Conf:        dbconf.MysqlMasterconf,
			Status:      false,
			MaxOpenConn: dbconf.MaxOpenConn,
			MaxIdleConn: dbconf.MaxIdleConn,
			MaxLifetime: dbconf.MaxLifetime,
		},

		DbSlave: &Handler{Name: "slave",
			Conf:        dbconf.MysqlSlaveconf,
			Status:      false,
			MaxOpenConn: dbconf.MaxOpenConn,
			MaxIdleConn: dbconf.MaxIdleConn,
			MaxLifetime: dbconf.MaxLifetime,
		},
	}

	connDB(ghandler.DbMaster)
	connDB(ghandler.DbSlave)

	return ghandler
}

func InitAllDB(conf config.BaseConfig, dbInstence map[string]*GroupHandler) {
	logging.Logger.Info("db.InitAllDB begin")

	DBs := conf.GetDBs()
	for name, dbconf := range DBs {

		//init mysql:second para is username:password@protocol(address)/dbname?param=value
		if dbconf.Enable != 0 {
			groupHandler := initGroupDB(dbconf)
			groupHandler.Name = name
			go Monitor(groupHandler)
			dbInstence[name] = groupHandler
		}
	}

	logging.Logger.Info("db.InitAllDB end")
}

func InitDB(name, mysqlConf string) *GroupHandler {
	logging.Logger.Info("db.InitAllDB begin")

	//init mysql:second para is username:password@protocol(address)/dbname?param=value
	groupHandler := initGroupDB(config.DatabaseConf{MysqlMasterconf: mysqlConf, MysqlSlaveconf: mysqlConf, Enable: 1})
	groupHandler.Name = name

	logging.Logger.Info("db.InitAllDB end")
	return groupHandler
}
