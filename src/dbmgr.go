package web

import (
        "os"
        //"fmt"
)

type dbrecord struct {
        Database
        cfg *DatabaseConfig
        useCount int
}

type DBManager struct {
        dbs []*dbrecord
}

var dbmanager = &DBManager{
        make([]*dbrecord, 0),
}

func GetDBManager() *DBManager { return dbmanager }

func (dbm *DBManager) findDatabase(cfg *DatabaseConfig) (rec *dbrecord) {
        for _, r := range dbm.dbs {
                if r.cfg.Host == cfg.Host &&
                        r.cfg.User == cfg.User &&
                        r.cfg.Password == cfg.Password &&
                        r.cfg.Database == cfg.Database {
                        rec = r
                        break
                }
        }
        return
}

func (dbm *DBManager) GetDatabase(cfg *DatabaseConfig) (db Database, err os.Error) {
        rec := dbm.findDatabase(cfg)
        if rec == nil {
                rec = &dbrecord{ NewDatabase(), cfg, 1 }
                err = rec.Connect(cfg.Host, cfg.User, cfg.Password, cfg.Database)
                if err != nil { /* TODO: error */ goto finish }
                dbm.dbs = append( dbm.dbs, rec )
        } else {
                rec.useCount += 1
        }
        db = Database(rec)
finish:
        return
}

func (dbm *DBManager) CloseAll() /* TODO: return status code */ {
        for _, r := range dbm.dbs {
                r.Database.Close()

                // reset useCount, in case external references exists
                r.useCount = 0
        }
        dbm.dbs = make([]*dbrecord, 0)
}

func (rec *dbrecord) Close() (err os.Error) {
        //fmt.Printf("close: dbrecord%v\n", rec)
        rec.useCount -= 1
        if rec.useCount == 0 {
                //rec.Database.Close()
                // TODO: remove the record from dbmanager.dbs
                //fmt.Printf("closed: %v\n", rec)
        }
        return
}

