package mysql

/*
 #include <stdlib.h>
 #include <mysql.h>

 char *_field_value_at(MYSQL_ROW row, int i) { return row[i]; }

 MYSQL_FIELD *_field_at(MYSQL_FIELD *f, int i) { return f + i; }
 char *_field_name(MYSQL_FIELD * f) { return f->name; }
 char *_field_org_name(MYSQL_FIELD * f) { return f->org_name; }
 char *_field_table(MYSQL_FIELD * f) { return f->table; }
 char *_field_org_table(MYSQL_FIELD * f) { return f->org_table; }
 char *_field_db(MYSQL_FIELD * f) { return f->db; }
 char *_field_catalog(MYSQL_FIELD * f) { return f->catalog; }
 char *_field_def(MYSQL_FIELD * f) { return f->def; }
 unsigned long _field_length(MYSQL_FIELD * f) { return f->length; }
 unsigned long _field_max_length(MYSQL_FIELD * f) { return f->max_length; }
 unsigned int _field_flags(MYSQL_FIELD * f) { return f->flags; }
 unsigned int _field_decimals(MYSQL_FIELD * f) { return f->decimals; }
 unsigned int _field_charsetnr(MYSQL_FIELD * f) { return f->charsetnr; }
 unsigned char _field_type(MYSQL_FIELD * f) { return (unsigned char)(f->type); }

 
 */
import "C"

import (
        "os"
        "sync"
        "unsafe"
        //"strconv"
        //"fmt"
)

type Connection struct {
        h *C.MYSQL
        mtx *sync.Mutex
}

type Statement struct {
        h *C.MYSQL_STMT
}

type ResultSet struct {
        h *C.MYSQL_RES

        NumFields uint
        AffectedRows uint64
}

type Field struct {
        Name string
        OriginalName string
        Table string
        OriginalTable string
        Database string
        Catalog string
        Default string
        Length uint
        MaxLength uint
        Flags FieldFlags //uint // Div flags
        Decimals uint
        Charset uint
        Type FieldType //uint8
}

func newError(msg string) (err os.Error) {
        prefix := "" // TODO: file:linum
        err = os.NewError(prefix + msg)
        return
}

func (conn *Connection) getLastError() os.Error {
        if err := C.mysql_error(conn.h); *err != 0 {
                return os.NewError(C.GoString(err))
        }
        return nil
}

func Open(host, user, pass, db string) (rconn *Connection, err os.Error) {
        conn := &Connection{};

        conn.h = C.mysql_init(nil);
        if conn.h == nil {
                err = newError("Can't init MySQL connection.");
                return
        } else {
                var socket *C.char
                port := 3306
                a1, a2, a3, a4, a5, a6 := C.CString(host), C.CString(user), C.CString(pass), C.CString(db), C.uint(port), socket
                rh := C.mysql_real_connect(conn.h, a1, a2, a3, a4, a5, a6, 0)
                if a1 != nil { C.free(unsafe.Pointer(a1)) }
                if a2 != nil { C.free(unsafe.Pointer(a2)) }
                if a3 != nil { C.free(unsafe.Pointer(a3)) }
                if a4 != nil { C.free(unsafe.Pointer(a4)) }

                err = conn.getLastError()
                if err != nil || rh != conn.h {
                        C.mysql_close(conn.h)
                        return
                }
        }

        if err == nil {
                // Create new connection mutex
                conn.mtx = new(sync.Mutex);
                rconn = conn
        }
        return
}

func (conn *Connection) Close() os.Error {
        C.mysql_close(conn.h)
        conn.h = nil
        //conn.mtx = nil
        return nil
}

func (conn *Connection) Query(sql string) (rs *ResultSet, err os.Error) {
        cs := C.CString(sql)
        conn.mtx.Lock()
        
        defer func() {
                conn.mtx.Unlock()
                C.free(unsafe.Pointer(cs))
        }()

        rc := C.mysql_query(conn.h, cs)
        if rc != 0 {
                err = conn.getLastError()
                return
        }

        rs = &ResultSet{
        h: nil,
        NumFields: uint(C.mysql_field_count(conn.h)),
        AffectedRows: uint64(C.mysql_affected_rows(conn.h)),
        }

        //res := C.mysql_use_result(conn.h)
        res := C.mysql_store_result(conn.h)
        if res == nil {
                if err = conn.getLastError(); err != nil {
                        return
                }
        } else {
                rs.h = res
        }

        // TODO: must call C.mysql_free_result(unsafe.Pointer(res))

        return
}

func (rs *ResultSet) Free() (err os.Error) {
        if rs.h != nil {
                //C.mysql_free_result(unsafe.Pointer(rs.h))
                C.mysql_free_result(rs.h)
                rs.h = nil
        }
        return
}

func (rs *ResultSet) GetNumFields() (num uint) {
        if rs.h != nil {
                num = uint(C.mysql_num_fields(rs.h))
        }
        return
}

func (rs *ResultSet) FetchFields() (fields []Field, err os.Error) {
        count := int(rs.GetNumFields())
        if count <= 0 {
                return
        }

        var fs *C.MYSQL_FIELD
        fs = C.mysql_fetch_fields(rs.h)
        if fs == nil {
                err = os.NewError("Can't fetch fields.")
                return
        }

        fields = make([]Field, count)
        for i := 0; i < count; i += 1 {
                f := C._field_at(fs, C.int(i))
                fields[i] = Field{
                Name: C.GoString(C._field_name(f)),
                OriginalName: C.GoString(C._field_org_name(f)),
                Table: C.GoString(C._field_table(f)),
                OriginalTable: C.GoString(C._field_org_table(f)),
                Database: C.GoString(C._field_db(f)),
                Catalog: C.GoString(C._field_catalog(f)),
                Default: C.GoString(C._field_def(f)),
                Length: uint(C._field_length(f)),
                MaxLength: uint(C._field_max_length(f)),
                Flags: FieldFlags(C._field_flags(f)),
                Decimals: uint(C._field_decimals(f)),
                Charset: uint(C._field_charsetnr(f)),
                Type: FieldType(C._field_type(f)),
                }
        }
        return
}

func (rs *ResultSet) FetchRow() (row []string, err os.Error) {
        if rs.h == nil {
                err = os.NewError("Invalid result set.")
                return
        }

        r := C.mysql_fetch_row(rs.h)
        if r == nil {
                return
        } else {
                count := int(C.mysql_num_fields(rs.h)) // ASSERT(count)?
                row = make([]string, count)
                for i:=0; i<count; i+=1 {
                        row[i] = C.GoString(C._field_value_at(r, C.int(i)))
                }
        }
        return
}

func (rs *ResultSet) ConvertRow(r []string) (row []interface{}, err os.Error) {
        count := int(rs.GetNumFields())
        if count <= 0 {
                //err = os.NewError("No fields.")
                return
        }
        if len(r) != count {
                err = os.NewError("Fields not matched.")
                return
        }

        fields := C.mysql_fetch_fields(rs.h)
        if fields == nil {
                err = os.NewError("Can't fetch fields.")
                return
        }

        row = make([]interface{}, count)
        for i := 0; i < count; i += 1 {
                f := C._field_at(fields, C.int(i))
                t := FieldType(C._field_type(f))
                flags := FieldFlags(C._field_flags(f))
                row[i] = t.convert(r[i], flags)
        }
        return
}
