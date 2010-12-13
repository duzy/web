package web

import (
        "os"
        //"xml"
        "json"
        "strings"
)

type AppConfig struct {
        Title string
        Model string /* TODO: make it AppModel  */

        /**
         * Session persister driver.
         */
        Persister *AppConfig_Persister

        /**
         * Default database config.
         */
        Database *AppConfig_Database

        /**
         * Named persisters.
         */
        Persisters map[string]*AppConfig_Persister

        /**
         * Named databases.
         */
        Databases map[string]*AppConfig_Database
}

type AppConfig_Persister struct {
        /**
         *  *AppConfig_PersisterFS or *AppConfig_Database
         */
        base interface{}
}

func (p *AppConfig_Persister) IsFS() (ok bool, v *AppConfig_PersisterFS) {
        v, ok = p.base.(*AppConfig_PersisterFS)
        return
}

func (p *AppConfig_Persister) IsDB() (ok bool, v *AppConfig_Database) {
        v, ok = p.base.(*AppConfig_Database)
        return
}

type AppConfig_PersisterFS struct {
        /**
         *  File location(must be a directory) for storing sessions.
         */
        Location string
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

/**
 eg1: { "type": "FS", "location": "/tmp/web/sessions" }
 eg2: { "type": "DB", "named":"dusell_2" }
 eg3: { "type": "DB",
        "host": "localhost",
        "user": "test",
        "password": "abc",
        "database": "dusell_2"
      }
 */
func parseJSONDecodedPersister(v interface{}) (p *AppConfig_Persister) {
        if m, ok := v.(map[string]interface{}); ok {
                if m["type"] == nil { goto finish }
                if typ, ok := m["type"].(string); ok {
                        switch typ {
                        case "DB":
                                p = new(AppConfig_Persister)
                                if m["named"] != nil {
                                        if s, ok := m["named"].(string); ok {
                                                // a string value will be replaced by a '*AppConfig_Database'
                                                p.base = s
                                        }
                                } else {
                                        p.base = parseJSONDecodedDatabase(v)
                                        // TODO: check "p.base != nil" ?
                                }
                        case "FS":
                                fs := new(AppConfig_PersisterFS)
                                if s, ok := m["location"].(string); ok {
                                        fs.Location = s
                                } else {
                                        // TODO: error: location is not string
                                }
                                p = new(AppConfig_Persister)
                                p.base = fs
                        }//switch(persister type)
                } else {
                        // TODO: error: 'type' is not a string value
                }
        }
finish:
        return
}

/**
 eg: { "host": "localhost",
       "user": "test",
       "password": "abc",
       "database": "dusell"
     }
 */
func parseJSONDecodedDatabase(v interface{}) (p *AppConfig_Database) {
        if m, ok := v.(map[string]interface{}); ok {
                p = new(AppConfig_Database)
                for name, value := range m {
                        switch name {
                        case "host": p.Host = value.(string)
                        case "user": p.User = value.(string)
                        case "password": p.Password = value.(string)
                        case "database": p.Database = value.(string)
                        }
                }
        }
        return
}

func loadAppConfigJSON(fn string) (cfg *AppConfig, err os.Error) {
        f, err := os.Open(fn, os.O_RDONLY, 0)
        if err != nil {
                // TODO: error handling
                goto finish
        }

        defer f.Close()

        m := make(map[string]interface{})
        
        dec := json.NewDecoder(f)
        if err = dec.Decode(&m); err != nil {
                // TODO: error handling
                goto finish
        }

        cfg = new(AppConfig)

        var i interface{}

        i = m["databases"]
        if i != nil {
                if v, ok := i.(map[string]interface{}); ok {
                        cfg.Databases = make(map[string]*AppConfig_Database)
                        for k, d := range v {
                                cfg.Databases[k] = parseJSONDecodedDatabase(d)
                        }
                }
        }

        i = m["persisters"]
        if i != nil {
                if v, ok := i.(map[string]interface{}); ok {
                        cfg.Persisters = make(map[string]*AppConfig_Persister)
                        for k, ip := range v {
                                p := parseJSONDecodedPersister(ip)
                                if s, ok := p.base.(string); ok {
                                        // TODO: check cfg.Databases != nil
                                        // convert named database persister
                                        p.base = cfg.Databases[s]
                                }
                                cfg.Persisters[k] = p
                        }
                }
        }

        

        i = m["title"]
        if s, ok := i.(string); ok { cfg.Title = s }

        i = m["model"]
        if s, ok := i.(string); ok { cfg.Model = s }

        i = m["persister"]
        switch v := i.(type) {
        case string:
                if cfg.Persisters != nil {
                        cfg.Persister = cfg.Persisters[v]
                }
        case map[string]interface{}:
                cfg.Persister = parseJSONDecodedPersister(v)
                if s, ok := cfg.Persister.base.(string); ok {
                        // TODO: check cfg.Databases != nil
                        cfg.Persister.base = cfg.Databases[s]
                }
        }

        i = m["database"]
        switch v := i.(type) {
        case string:
                if cfg.Databases != nil {
                        cfg.Database = cfg.Databases[v]
                }
        case map[string]interface{}:
                cfg.Database = parseJSONDecodedDatabase(i)
        }

finish:
        return
}

func loadAppConfigXML(fn string) (cfg *AppConfig, err os.Error) {
        // TODO: ...
        return
}
