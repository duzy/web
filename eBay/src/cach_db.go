package eBay

import (
        "os"
        "../_obj/web"
)

type dbCache struct {
        db web.Database
}

// NewDBCache accepts parameters in this fixed order: host, user,
// password, database.
func NewDBCache(params ...interface{}) (c Cacher, err os.Error) {
        a := []interface{}(params)
        var host, user, password, database string
        host, user, password, database = a[0].(string), a[1].(string), a[2].(string), a[3].(string)
        db := web.NewDatabase(host, user, password, database)
        dbc := &dbCache{ db }
        c = Cacher(dbc);
        return
}
