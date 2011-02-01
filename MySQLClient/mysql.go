package mysql

/*
 #include <stdlib.h>
 #include <string.h>
 #include <mysql.h>

 short int      _conv_int16(void *p) { return *((short int*)p); }
 int            _conv_int32(void *p) { return *((int*)p); }
 long long int  _conv_int64(void *p) { return *((long long int*)p); }
 float          _conv_float32(void *p) { return *((float*)p); }
 double         _conv_float64(void *p) { return *((double*)p); }

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

 MYSQL_BIND *_bind_new(int count) {
   MYSQL_BIND *binds = malloc(sizeof(MYSQL_BIND) * count);
   memset(binds, 0, sizeof(MYSQL_BIND) * count);
   return binds;
 }
 void _bind_delete(MYSQL_BIND *binds) { free(binds); }
 void _bind_set(MYSQL_BIND *binds, int i,
     enum enum_field_types type,
     void *buf, unsigned long buflen,
     void *len, void *nul, void *error)
 {
   binds[i].buffer_type = type;
   binds[i].buffer = buf;
   binds[i].buffer_length = buflen;
   binds[i].length = (unsigned long *)len;
   binds[i].is_null = (my_bool *)nul;
   binds[i].error = (my_bool *)error;
 }

 */
import "C"

import (
        "os"
        "sync"
        "unsafe"
        //"strconv"
        "fmt"
)

type Connection struct {
        h *C.MYSQL
        mtx *sync.Mutex
}

type Statement struct {
        h *C.MYSQL_STMT

        params []bind
        result []bind
}

type ResultSet struct {
        h *C.MYSQL_RES

        NumFields uint
        AffectedRows uint64

        InsertId uint64
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

type bind struct {
        buffer []byte;
	buffer_type FieldType;
        length int;
        is_null [1]byte;
        error [1]byte;
}

func (b *bind) makeBuffer(v interface{}) {
        switch a := v.(type) {
        case int:
                l := unsafe.Sizeof(int(0))
                b.buffer_type = FIELD_TYPE_LONG
                b.buffer = make([]byte, l)
                for i := uint(0); i < uint(l); i += 1 {
                        b.buffer[i] = uint8((a >> (i*8)) & 0xFF)
                }

        case string:
                b.buffer_type = FIELD_TYPE_STRING
                b.buffer = []byte(a)
        }
}

func (b *bind) makeBufferForField(field *C.MYSQL_FIELD) {
        l := int(field.length)
        b.buffer_type = FieldType(field._type)
        switch b.buffer_type {
        case FIELD_TYPE_DECIMAL:        l = unsafe.Sizeof(float64(0))
        case FIELD_TYPE_TINY:           l = unsafe.Sizeof(int8(0))
	case FIELD_TYPE_SHORT:          l = unsafe.Sizeof(int16(0))
        case FIELD_TYPE_LONG:           l = unsafe.Sizeof(int(0))
	case FIELD_TYPE_FLOAT:          l = unsafe.Sizeof(float32(0))
        case FIELD_TYPE_DOUBLE:         l = unsafe.Sizeof(float64(0))
	case FIELD_TYPE_NULL:           //l = unsafe.Sizeof()
        case FIELD_TYPE_TIMESTAMP:      //l = unsafe.Sizeof()
	case FIELD_TYPE_LONGLONG:       l = unsafe.Sizeof(int64(0))
        case FIELD_TYPE_INT24:          l = unsafe.Sizeof(int(0))
	case FIELD_TYPE_DATE:           //l = unsafe.Sizeof()
        case FIELD_TYPE_TIME:           //l = unsafe.Sizeof()
	case FIELD_TYPE_DATETIME:       //l = unsafe.Sizeof()
        case FIELD_TYPE_YEAR:           //l = unsafe.Sizeof()
	case FIELD_TYPE_NEWDATE:        //l = unsafe.Sizeof()
        case FIELD_TYPE_VARCHAR:        //l = unsafe.Sizeof()
	case FIELD_TYPE_BIT:            //l = unsafe.Sizeof()
        case FIELD_TYPE_NEWDECIMAL:     //l = unsafe.Sizeof()
	case FIELD_TYPE_ENUM:           //l = unsafe.Sizeof()
	case FIELD_TYPE_SET:            //l = unsafe.Sizeof()
	case FIELD_TYPE_TINY_BLOB:      //l = unsafe.Sizeof()
	case FIELD_TYPE_MEDIUM_BLOB:    //l = unsafe.Sizeof()
	case FIELD_TYPE_LONG_BLOB:      //l = unsafe.Sizeof()
	case FIELD_TYPE_BLOB:           //l = unsafe.Sizeof()
	case FIELD_TYPE_VAR_STRING:     //l = unsafe.Sizeof()
	case FIELD_TYPE_STRING:         //l = unsafe.Sizeof()
	case FIELD_TYPE_GEOMETRY:       //l = unsafe.Sizeof()
        }
        if 0 < l {
                b.buffer = make([]byte, l)
        }
}

func (b *bind) set(binds *C.MYSQL_BIND, i int) {
        bufptr := unsafe.Pointer(nil)
        if 0 < len(b.buffer) {
                bufptr = unsafe.Pointer(&b.buffer[0])
        }
        
        C._bind_set(
                binds, C.int(i),
                uint32(b.buffer_type),
                bufptr, C.ulong(len(b.buffer)),
                unsafe.Pointer(&b.length),
                unsafe.Pointer(&b.is_null),
                unsafe.Pointer(&b.error))
}

func (b *bind) makeValue() (v interface{}, ok bool) {
        //fmt.Printf("value: [%v] %v\n", b.buffer_type, b.buffer)
        if len(b.buffer) <= 0 {
                return
        }

        if b.is_null[0] == 1 {
                v, ok =  nil, true
                return
        }

        ptr := unsafe.Pointer(&b.buffer[0])

        switch b.buffer_type {
        case FIELD_TYPE_TINY:
                if b.length == 1 {
                        v, ok = uint8(b.buffer[0]), true
                }
	case FIELD_TYPE_SHORT:
                if b.length == 2 {
                        v, ok = int16(C._conv_int16(ptr)), true
                }
        case FIELD_TYPE_LONG:
                if b.length == 4 {
                        v, ok = int(C._conv_int32(ptr)), true
                }
	case FIELD_TYPE_LONGLONG:
                if b.length == 8 {
                        v, ok = int64(C._conv_int64(ptr)), true
                }
	case FIELD_TYPE_FLOAT:
                if b.length == 4 {
                        v, ok = float32(C._conv_float32(ptr)), true
                }
        case FIELD_TYPE_DECIMAL:        fallthrough
        case FIELD_TYPE_NEWDECIMAL:     fallthrough
        case FIELD_TYPE_DOUBLE:
                if b.length == 8 {
                        v, ok = float64(C._conv_float64(ptr)), true
                }
        case FIELD_TYPE_INT24:
                if b.length == 3 {
                        v, ok = int(C._conv_int32(ptr)), true
                }

	case FIELD_TYPE_VAR_STRING:     fallthrough
	case FIELD_TYPE_STRING:         fallthrough
        case FIELD_TYPE_VARCHAR:
                if 0 < b.length {
                        v, ok = string(b.buffer[0:b.length]), true
                }

                /*
	case FIELD_TYPE_TINY_BLOB:      fallthrough
	case FIELD_TYPE_MEDIUM_BLOB:    fallthrough
	case FIELD_TYPE_LONG_BLOB:      fallthrough
	case FIELD_TYPE_BLOB:           fallthrough
                 */

                // TODO: more types...
        }
        return
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
        InsertId: uint64(C.mysql_insert_id(conn.h)),
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

func (conn *Connection) Prepare(sql string) (stmt *Statement, err os.Error) {
        cs := C.CString(sql)
        conn.mtx.Lock()

        defer func() {
                conn.mtx.Unlock()
                C.free(unsafe.Pointer(cs))
        }()

        stmt = &Statement{}
        stmt.h = C.mysql_stmt_init(conn.h)
        if stmt.h == nil {
                err, stmt = conn.getLastError(), nil
                if err == nil {
                        err = os.NewError("mysql_stmt_init failed")
                }
                return
        }
        
        rc := C.mysql_stmt_prepare(stmt.h, cs, C.ulong(C.strlen(cs)))
        if rc != 0 {
                err = stmt.getLastError()
                if err == nil {
                        err = os.NewError("mysql_stmt_prepare failed")
                }

                stmt.Close()
                stmt = nil
                return
        }
        
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

func (rs *ResultSet) GetNumRows() (num uint64) {
        if rs.h != nil {
                num = uint64(C.mysql_num_rows(rs.h))
        }
        return
}

func (rs *ResultSet) RowSeek() (num uint64) {
        if rs.h != nil {
                num = uint64(C.mysql_num_rows(rs.h))
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

func (stmt *Statement) getLastError() os.Error {
        if err := C.mysql_stmt_error(stmt.h); *err != 0 {
                return os.NewError(C.GoString(err))
        }
        return nil
}

func (stmt *Statement) BindParams(params ...interface{}) (err os.Error) {
        if stmt.h == nil {
                err = os.NewError("Invalid statement!")
                return
        }

        count := int(C.mysql_stmt_param_count(stmt.h))
        if len(params) != count {
                err = os.NewError("wrong number of parameters")
                return
        }

        if count == 0 {
                return
        }

        binds := C._bind_new(C.int(count))
        defer func() {
                C._bind_delete(binds)
        }()

        stmt.params = make([]bind, count)
        
        for i, a := range params {
                stmt.params[i].makeBuffer(a)
                if stmt.params[i].buffer == nil {
                        err = os.NewError(fmt.Sprintf("Unsupported parameter type: ", a))
                        return
                }

                stmt.params[i].set(binds, i)

                rc := C.mysql_stmt_bind_param(stmt.h, binds)
                if rc != 0 {
                        err = stmt.getLastError()
                        return
                }
        }
        return
}

func (stmt *Statement) ParamCount() (n int) {
        if stmt.h != nil {
                n = int(C.mysql_stmt_param_count(stmt.h))
        }
        return
}

func (stmt *Statement) Execute() (err os.Error) {
        if stmt.h == nil {
                err = os.NewError("Invalid statement!")
                return
        }

        rc := C.mysql_stmt_execute(stmt.h)
        if rc != 0 {
                err = stmt.getLastError()
                return
        }

        meta := C.mysql_stmt_result_metadata(stmt.h)
        if meta != nil {
                if count := C.mysql_num_fields(meta); 0 < count {
                        binds := C._bind_new(C.int(count))
                        stmt.result = make([]bind, int(count))
                        for i := C.uint(0); i < count; i += 1 {
                                field := C.mysql_fetch_field_direct(meta, i)
                                stmt.result[i].makeBufferForField(field)
                                stmt.result[i].set(binds, int(i))
                        }

                        rc := C.mysql_stmt_bind_result(stmt.h, binds)
                        if rc != 0 {
                                err = stmt.getLastError()
                                return
                        }

                        C._bind_delete(binds)
                        C.mysql_free_result(meta)
                }
        }

        rc = C.mysql_stmt_store_result(stmt.h)
        if rc != 0 {
                err = stmt.getLastError()
                return
        }

        return
}

func (stmt *Statement) SQLState() (s string) {
        if stmt.h != nil {
                cs := C.mysql_stmt_sqlstate(stmt.h)
                s = C.GoString(cs)
        }
        return
}

func (stmt *Statement) FieldCount() (n uint) {
        if stmt.h != nil {
                n = uint(C.mysql_stmt_field_count(stmt.h))
        }
        return
}

func (stmt *Statement) NumRows() (n uint64) {
        if stmt.h != nil {
                n = uint64(C.mysql_stmt_num_rows(stmt.h))
        }
        return
}

func (stmt *Statement) AffectedRows() (n uint64) {
        if stmt.h != nil {
                n = uint64(C.mysql_stmt_affected_rows(stmt.h))
        }
        return
}

func (stmt *Statement) InsertId() (n uint64) {
        if stmt.h != nil {
                n = uint64(C.mysql_stmt_insert_id(stmt.h))
        }
        return
}

func (stmt *Statement) Fetch() (row []interface{}, err os.Error) {
        if stmt.h == nil {
                err = os.NewError("Invalid statement!")
                return
        }

        rc := C.mysql_stmt_fetch(stmt.h)

        //fmt.Printf("result: [%d], %v\n", rc, stmt.result)

        if rc != 0 {
                if rc == 100 { // MYSQL_NO_DATA
                        err = os.EOF
                        return
                } else if rc == 101 { // MYSQL_DATA_TRUNCATED
                        // ...
                } else {
                        err = stmt.getLastError()
                        return
                }
        }

        if stmt.result != nil {
                row = make([]interface{}, len(stmt.result))
                for i := 0; i < len(stmt.result); i += 1 {
                        row[i], _ = stmt.result[i].makeValue()
                }
        }
        
        return
}

func (stmt *Statement) Reset() (err os.Error) {
        if stmt.h == nil {
                err = os.NewError("Invalid statement!")
                return
        }

        rc := C.mysql_stmt_reset(stmt.h)
        if rc != 0 {
                err = stmt.getLastError()
                return
        }

        return
}

func (stmt *Statement) Close() (err os.Error) {
        if stmt.h == nil {
                err = os.NewError("Invalid statement!")
                return
        }

        rc := C.mysql_stmt_free_result(stmt.h)
        if rc != 0 {
                err = stmt.getLastError()
                // NOTE: don't return here
        }

        rc = C.mysql_stmt_close(stmt.h)
        stmt.h = nil
        stmt.params = nil
        stmt.result = nil

        if rc != 0 {
                err = stmt.getLastError()
                return
        }

        return
}

