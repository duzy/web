package web

import (
        "./_obj/web"
        "testing"
        "reflect"
)

func check(t *testing.T, a, b interface{}) {
        if !reflect.DeepEqual(a, b) {
                t.Error("[%v]!=[%v]", a, b)
        }
}

func TestLoadAppConfig(t *testing.T) {
        cfg, err := web.LoadAppConfig("test_app.json")
        if err != nil { t.Error(err) }

        check(t, cfg.Title, "test app via json")
        check(t, cfg.Model, "CGI")
        check(t, len(cfg.Databases), 2)
        check(t, len(cfg.Persisters), 3)

        d := cfg.Databases["dusell"]
        if d == nil { t.Error("no database 'dusell'") }
        check(t, d.Host, "localhost")
        check(t, d.User, "test")
        check(t, d.Password, "abc")
        check(t, d.Database, "dusell")

        d = cfg.Databases["dusell_2"]
        if d == nil { t.Error("no database 'dusell_2'") }
        check(t, d.Host, "localhost")
        check(t, d.User, "test")
        check(t, d.Password, "abc")
        check(t, d.Database, "dusell_2")

        p := cfg.Persisters["per_1"]
        if p == nil { t.Error("no persister 'per_1'") }
        if ok, v := p.IsFS(); ok {
                check(t, v.Location, "/tmp/web/sessions")
        } else { t.Error("per_1 is not FS persister") }

        p = cfg.Persisters["per_2"]
        if p == nil { t.Error("no persister 'per_2'") }
        if ok, v := p.IsDB(); ok {
                check(t, v.Host, "localhost")
                check(t, v.User, "test")
                check(t, v.Password, "abc")
                check(t, v.Database, "dusell_2")
        } else { t.Error("per_2 is not DB persister") }

        p = cfg.Persisters["per_3"]
        if p == nil { t.Error("no persister 'per_3'") }
        if ok, v := p.IsDB(); ok {
                check(t, v.Host, "localhost")
                check(t, v.User, "test")
                check(t, v.Password, "abc")
                check(t, v.Database, "dusell_2")
        } else { t.Error("per_3 is not DB persister") }

        if cfg.Database == nil { t.Error("database is nil") }
        if cfg.Persister == nil { t.Error("persister is nil") }

        check(t, cfg.Database.Database, "dusell")

        if ok, v := cfg.Persister.IsFS(); ok {
                check(t, v.Location, "/tmp/web/sessions")
        } else {
                t.Error("not FS persister")
        }
}
