package web

import (
        "os"
        //"io"
        "io/ioutil"
        "fmt"
        //"xml"
        "json"
        "strings"
)

type internalAppConfig struct {
        Title string
        Model string
        Persister interface{} // string or map[string]string
        Database interface{} // string or map[string]string
        Persisters map[string] struct {
                Type string

                // Fields for FS type
                Location string

                // Fields for DB type
                Named string // for Named database
                Host string
                User string
                Password string
                Database string
        }
        Databases map[string] struct {
                Host string
                User string
                Password string
                Database string
        }
}

type AppConfig struct {
        Title string
        Model string /* TODO: make it AppModel  */

        /**
         * Session persister driver.
         */
        Persister AppConfig_Persister

        /**
         * Default database config.
         */
        Database *AppConfig_Database

        /**
         * Named persisters.
         */
        Persisters map[string]AppConfig_Persister

        /**
         * Named databases.
         */
        Databases map[string]*AppConfig_Database
}

/**
 *  *AppConfig_PersisterFS or *AppConfig_PersisterDB
 */
type AppConfig_Persister interface{}

type AppConfig_PersisterFS struct {
        /**
         *  File location(must be a directory) for storing sessions.
         */
        Location string
}

type AppConfig_PersisterDB struct {
        AppConfig_Database
}

type AppConfig_Database struct {
        Host string
        User string
        Password string
        Database string
}

/**
 *  Load AppConfig from a JSON or XML file.
 */
func LoadAppConfig(fn string) (cfg *AppConfig, err os.Error) {
        switch {
        case strings.HasSuffix(fn, ".json"):
                cfg, err = loadAppConfigJSON(fn)
        case strings.HasSuffix(fn, ".xml"):
                cfg, err = loadAppConfigXML(fn)
        }
        return
}

func loadAppConfigJSON(fn string) (cfg *AppConfig, err os.Error) {
        f, err := os.Open(fn, os.O_RDONLY, 0)
        if err != nil { return }

        defer f.Close()

        data, err := ioutil.ReadAll(f)
        if err != nil { return }

        c := new(internalAppConfig)
        err = json.Unmarshal(data, c)
        if err != nil { return }

        cfg = new(AppConfig)
        cfg.Databases = make(map[string]*AppConfig_Database)
        for k, d := range c.Databases {
                cfg.Databases[k] = &AppConfig_Database{
                Host: d.Host,
                User: d.User,
                Password: d.Password,
                Database: d.Database,
                }
        }
        cfg.Persisters = make(map[string]AppConfig_Persister)
        for k, p := range c.Persisters {
                switch p.Type {
                case "FS":
                        cfg.Persisters[k] = &AppConfig_PersisterFS {
                        Location: p.Location,
                        }
                case "DB":
                        if p.Named != "" {
                                db := cfg.Databases[p.Named]
                                if db != nil {
                                        cfg.Persisters[k] = &AppConfig_PersisterDB { *db }
                                }
                        } else {
                                cfg.Persisters[k] = &AppConfig_PersisterDB{
                                        AppConfig_Database {
                                        Host: p.Host,
                                        User: p.User,
                                        Password: p.Password,
                                        Database: p.Database,
                                        },
                                }
                        }
                }
        }
        cfg.Title = c.Title
        cfg.Model = c.Model

        switch v := c.Database.(type) {
        case string:
                if cfg.Databases != nil {
                        cfg.Database = cfg.Databases[v]
                }
        case map[string]interface{}:
                cfg.Database = new(AppConfig_Database)
                for name, value := range v {
                        switch name {
                        case "host": cfg.Database.Host = value.(string)
                        case "user": cfg.Database.User = value.(string)
                        case "password": cfg.Database.Password = value.(string)
                        case "database": cfg.Database.Database = value.(string)
                        }
                }
        }

        switch v := c.Persister.(type) {
        case string:
                if cfg.Persisters != nil {
                        cfg.Persister = cfg.Persisters[v]
                }
        case map[string]interface{}:
                switch fmt.Sprintf("%v",v["type"]) {
                case "DB":
                        var db *AppConfig_Database
                        if v["named"] != nil {
                                if cfg.Databases == nil { break }
                                db = cfg.Databases[fmt.Sprintf("%v", v["named"])]
                        } else {
                                db = &AppConfig_Database{
                                Host: v["host"].(string),
                                User: v["user"].(string),
                                Password: v["password"].(string),
                                Database: v["database"].(string),
                                }
                        }
                        if db != nil {
                                cfg.Persister = &AppConfig_PersisterDB{ *db }
                        }
                case "FS":
                        cfg.Persister = &AppConfig_PersisterFS{
                        Location: fmt.Sprintf("%v", v["location"]),
                        }
                }
        }
        return
}

func loadAppConfigXML(fn string) (cfg *AppConfig, err os.Error) {
        // TODO: ...
        return
}
