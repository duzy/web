package web

import (
        "testing"
        "reflect"
        "fmt"
)

func check(t *testing.T, a, b interface{}) {
        if !reflect.DeepEqual(a, b) {
                t.Errorf("%v != %v", a, b)
        }
}

func TestLoadAppConfig(t *testing.T) {
        cfg, err := LoadAppConfig("test_app.json")
        if err != nil { t.Error(err); return }
        if cfg == nil { t.Error("no AppConfig loaded"); return }

        check(t, cfg.Title, "test app via json")
        check(t, cfg.Model, "CGI")
        check(t, len(cfg.Databases), 2)
        check(t, len(cfg.Persisters), 3)

        d := cfg.Databases["dusell"]
        if d == nil { t.Error("no database 'dusell'") }
        check(t, d.Host, "localhost")
        check(t, d.User, "dusellco_test")
        check(t, d.Password, "abc")
        check(t, d.Database, "dusellco_dusell")

        d = cfg.Databases["dusell_2"]
        if d == nil { t.Error("no database 'dusell_2'") }
        check(t, d.Host, "localhost")
        check(t, d.User, "test")
        check(t, d.Password, "abc")
        check(t, d.Database, "dusell_2")

        var per interface{}
        per = cfg.Persisters["per_1"]
        if per == nil { t.Error("no persister 'per_1'") }
        if fmt.Sprintf("%v",per) != "&{/tmp/web-test/sessions}" {
                t.Error("per_1 is unexpected:", per)
        }
        if v, ok := per.(*PersisterConfigFS); ok {
                check(t, v.Location, "/tmp/web-test/sessions")
        } else { t.Error("per_1 is not FS persister:", per) }

        per = cfg.Persisters["per_2"]
        if per == nil { t.Error("no persister 'per_2'") }
        if fmt.Sprintf("%v",per) != "&{{localhost test abc dusell_2}}" {
                t.Error("per_2 is unexpected:", per)
        }
        if v, ok := per.(*PersisterConfigDB); ok {
                check(t, v.Host, "localhost")
                check(t, v.User, "test")
                check(t, v.Password, "abc")
                check(t, v.Database, "dusell_2")
        } else { t.Error("per_2 is not DB persister:", per) }

        per = cfg.Persisters["per_3"]
        if per == nil { t.Error("no persister 'per_3'") }
        if fmt.Sprintf("%v",per) != "&{{localhost test abc dusell_2}}" {
                t.Error("per_3 is unexpected:", per)
        }
        if v, ok := per.(*PersisterConfigDB); ok {
                check(t, v.Host, "localhost")
                check(t, v.User, "test")
                check(t, v.Password, "abc")
                check(t, v.Database, "dusell_2")
        } else { t.Error("per_3 is not DB persister:", per) }

        if cfg.Database == nil { t.Error("database is nil") }
        if cfg.Persister == nil { t.Error("persister is nil") }

        check(t, cfg.Database.Database, "dusellco_dusell")

        if v, ok := cfg.Persister.(*PersisterConfigFS); ok {
                check(t, v.Location, "/tmp/web-test/sessions")
        } else {
                t.Error("not FS persister:", cfg.Persister)
        }
}
