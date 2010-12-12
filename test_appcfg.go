package main

import (
        "./_obj/web"
        "fmt"
        "reflect"
)

func check(a, b interface{}) {
        if !reflect.DeepEqual(a, b) {
                panic(fmt.Sprintf("[%v] != [%v]", a, b))
        }
}

func main() {
        var cfg *web.AppConfig
        cfg, err := web.LoadAppConfig("test_app.json")
        if err != nil {
                fmt.Printf("error: %s", err);
                goto finish
        }

        check(cfg.Title, "test app via json")
        check(cfg.Model, "CGI")
        check(len(cfg.Databases), 2)
        check(len(cfg.Persisters), 3)

        d := cfg.Databases["dusell"]
        if d == nil { panic("no database 'dusell'") }
        check(d.Host, "localhost")
        check(d.User, "test")
        check(d.Password, "abc")
        check(d.Database, "dusell")

        d = cfg.Databases["dusell_2"]
        if d == nil { panic("no database 'dusell_2'") }
        check(d.Host, "localhost")
        check(d.User, "test")
        check(d.Password, "abc")
        check(d.Database, "dusell_2")

        p := cfg.Persisters["per_1"]
        if p == nil { panic("no persister 'per_1'") }
        if ok, v := p.IsFS(); ok {
                check(v.Location, "/tmp/web/sessions")
        } else { panic("per_1 is not FS persister") }

        p = cfg.Persisters["per_2"]
        if p == nil { panic("no persister 'per_2'") }
        if ok, v := p.IsDB(); ok {
                check(v.Host, "localhost")
                check(v.User, "test")
                check(v.Password, "abc")
                check(v.Database, "dusell_2")
        } else { panic("per_1 is not FS persister") }

        p = cfg.Persisters["per_3"]
        if p == nil { panic("no persister 'per_3'") }
        if ok, v := p.IsDB(); ok {
                check(v.Host, "localhost")
                check(v.User, "test")
                check(v.Password, "abc")
                check(v.Database, "dusell_2")
        } else { panic("per_1 is not FS persister") }

        if cfg.Database == nil { panic("database is nil") }
        if cfg.Persister == nil { panic("persister is nil") }

        check(cfg.Database.Database, "dusell")

        if ok, v := cfg.Persister.IsFS(); ok {
                check(v.Location, "/tmp/web/sessions")
        } else {
                panic("not FS persister")
        }

finish:
}
