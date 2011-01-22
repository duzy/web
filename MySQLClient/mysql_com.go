package mysql

import (
        "strconv"
        "fmt"
)

type FieldType uint8
type FieldFlag uint
type FieldFlags uint

const ( // see enum enum_field_types in mysql_com.h
        FIELD_TYPE_DECIMAL      FieldType = 0 + iota
        FIELD_TYPE_TINY         // = 1
	FIELD_TYPE_SHORT        // = 2
        FIELD_TYPE_LONG         // = 3
	FIELD_TYPE_FLOAT        // = 4
        FIELD_TYPE_DOUBLE       // = 5
	FIELD_TYPE_NULL         // = 6
        FIELD_TYPE_TIMESTAMP    // = 7
	FIELD_TYPE_LONGLONG     // = 8
        FIELD_TYPE_INT24        // = 9
	FIELD_TYPE_DATE         // = 10
        FIELD_TYPE_TIME         // = 11
	FIELD_TYPE_DATETIME     // = 12
        FIELD_TYPE_YEAR         // = 13
	FIELD_TYPE_NEWDATE      // = 14
        FIELD_TYPE_VARCHAR      // = 15
	FIELD_TYPE_BIT          // = 16

        FIELD_TYPE_NEWDECIMAL   FieldType = 246
	FIELD_TYPE_ENUM         FieldType = 247
	FIELD_TYPE_SET          FieldType = 248
	FIELD_TYPE_TINY_BLOB    FieldType = 249
	FIELD_TYPE_MEDIUM_BLOB  FieldType = 250
	FIELD_TYPE_LONG_BLOB    FieldType = 251
	FIELD_TYPE_BLOB         FieldType = 252
	FIELD_TYPE_VAR_STRING   FieldType = 253
	FIELD_TYPE_STRING       FieldType = 254
	FIELD_TYPE_GEOMETRY     FieldType = 255
)

const ( // see *_FLAG in mysql_com.h
        FIELD_FLAG_NOT_NULL             FieldFlag = 1
        FIELD_FLAG_PRI_KEY              FieldFlag = 2
        FIELD_FLAG_UNIQUE_KEY           FieldFlag = 4
        FIELD_FLAG_MULTIPLE_KEY         FieldFlag = 8
        FIELD_FLAG_BLOB_FLAG            FieldFlag = 16
        FIELD_FLAG_UNSIGNED             FieldFlag = 32
        FIELD_FLAG_ZEROFILL             FieldFlag = 64
        FIELD_FLAG_BINARY               FieldFlag = 128

        FIELD_FLAG_ENUM                 FieldFlag = 256
        FIELD_FLAG_AUTO_INCREMENT       FieldFlag = 512
        FIELD_FLAG_TIMESTAMP            FieldFlag = 1024
        FIELD_FLAG_SET                  FieldFlag = 2048
        FIELD_FLAG_NO_DEFAULT_VALUE     FieldFlag = 4096
        FIELD_FLAG_ON_UPDATE_NOW        FieldFlag = 8192
        FIELD_FLAG_NUM                  FieldFlag = 32768
        FIELD_FLAG_PART_KEY             FieldFlag = 16384
        FIELD_FLAG_GROUP                FieldFlag = 32768 //same as FIELD_FLAG_NUM???
        FIELD_FLAG_UNIQUE               FieldFlag = 65536
        FIELD_FLAG_BINCMP               FieldFlag = 131072
        FIELD_FLAG_GET_FIXED_FIELDS     FieldFlag = (1 << 18)
        FIELD_FLAG_FIELD_IN_PART_FUNC   FieldFlag = (1 << 19)
        FIELD_FLAG_FIELD_IN_ADD_INDEX   FieldFlag = (1 << 20)
        FIELD_FLAG_FIELD_IS_RENAMED     FieldFlag = (1 << 21)
)

func (fs FieldFlags) Has(f FieldFlag) bool {
        return (uint(fs) & uint(f)) != 0
}

func (t FieldType) convert(s string, flags FieldFlags) (v interface{}) {
        switch t {
        default:
                v = s
        case FIELD_TYPE_TINY, FIELD_TYPE_SHORT, FIELD_TYPE_LONG:
                if flags.Has(FIELD_FLAG_UNSIGNED) {
                        v, _ = strconv.Atoui(s)
                } else {
                        v, _ = strconv.Atoi(s)
                }
	case FIELD_TYPE_FLOAT:
                v, _ = strconv.Atof32(s)
        case FIELD_TYPE_DOUBLE:
                v, _ = strconv.Atof64(s)
	case FIELD_TYPE_LONGLONG:
                if flags.Has(FIELD_FLAG_UNSIGNED) {
                        v, _ = strconv.Atoui64(s)
                } else {
                        v, _ = strconv.Atoi64(s)
                }
                /*
        case FIELD_TYPE_DECIMAL:
	case FIELD_TYPE_NULL:
        case FIELD_TYPE_TIMESTAMP:
        case FIELD_TYPE_INT24:
	case FIELD_TYPE_DATE:
        case FIELD_TYPE_TIME:
	case FIELD_TYPE_DATETIME:
        case FIELD_TYPE_YEAR:
	case FIELD_TYPE_NEWDATE:
        case FIELD_TYPE_VARCHAR:
	case FIELD_TYPE_BIT:
        case FIELD_TYPE_NEWDECIMAL:
	case FIELD_TYPE_ENUM:
	case FIELD_TYPE_SET:
	case FIELD_TYPE_TINY_BLOB:
	case FIELD_TYPE_MEDIUM_BLOB:
	case FIELD_TYPE_LONG_BLOB:
	case FIELD_TYPE_BLOB:
	case FIELD_TYPE_VAR_STRING:
	case FIELD_TYPE_STRING:
	case FIELD_TYPE_GEOMETRY:
                */
        }
        return
}

func (t FieldType) String() (s string) {
        switch t {
        case FIELD_TYPE_DECIMAL:        s = "FIELD_TYPE_DECIMAL"
        case FIELD_TYPE_TINY:           s = "FIELD_TYPE_TINY"
	case FIELD_TYPE_SHORT:          s = "FIELD_TYPE_SHORT"
        case FIELD_TYPE_LONG:           s = "FIELD_TYPE_LONG"
	case FIELD_TYPE_FLOAT:          s = "FIELD_TYPE_FLOAT"
        case FIELD_TYPE_DOUBLE:         s = "FIELD_TYPE_DOUBLE"
	case FIELD_TYPE_NULL:           s = "FIELD_TYPE_NULL"
        case FIELD_TYPE_TIMESTAMP:      s = "FIELD_TYPE_TIMESTAMP"
	case FIELD_TYPE_LONGLONG:       s = "FIELD_TYPE_LONGLONG"
        case FIELD_TYPE_INT24:          s = "FIELD_TYPE_INT24"
	case FIELD_TYPE_DATE:           s = "FIELD_TYPE_DATE"
        case FIELD_TYPE_TIME:           s = "FIELD_TYPE_TIME"
	case FIELD_TYPE_DATETIME:       s = "FIELD_TYPE_DATETIME"
        case FIELD_TYPE_YEAR:           s = "FIELD_TYPE_YEAR"
	case FIELD_TYPE_NEWDATE:        s = "FIELD_TYPE_NEWDATE"
        case FIELD_TYPE_VARCHAR:        s = "FIELD_TYPE_VARCHAR"
	case FIELD_TYPE_BIT:            s = "FIELD_TYPE_BIT"
        case FIELD_TYPE_NEWDECIMAL:     s = "FIELD_TYPE_NEWDECIMAL"
	case FIELD_TYPE_ENUM:           s = "FIELD_TYPE_ENUM"
	case FIELD_TYPE_SET:            s = "FIELD_TYPE_SET"
	case FIELD_TYPE_TINY_BLOB:      s = "FIELD_TYPE_TINY_BLOB"
	case FIELD_TYPE_MEDIUM_BLOB:    s = "FIELD_TYPE_MEDIUM_BLOB"
	case FIELD_TYPE_LONG_BLOB:      s = "FIELD_TYPE_LONG_BLOB"
	case FIELD_TYPE_BLOB:           s = "FIELD_TYPE_BLOB"
	case FIELD_TYPE_VAR_STRING:     s = "FIELD_TYPE_VAR_STRING"
	case FIELD_TYPE_STRING:         s = "FIELD_TYPE_STRING"
	case FIELD_TYPE_GEOMETRY:       s = "FIELD_TYPE_GEOMETRY"
        default: s = fmt.Sprintf("FieldType(%d)", uint8(t))
        }
        return
}
