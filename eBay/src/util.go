package eBay

import (
        "os"
        "fmt"
        "strconv"
        "strings"
        "reflect"
        //"runtime"
        "./_obj/web"
)

type NamedValue struct {
        Name string
        Value interface{}
}

func GetValues(a []NamedValue) (res []interface{}) {
        if a == nil {
                return
        }

        res = make([]interface{}, len(a))
        for i, v := range a {
                res[i] = v.Value
        }
        return
}

// AssignFields assign left-hand-side fields to right-hand-side fields of
// same names.
func AssignFields_(lhs, rhs interface{}) (err os.Error) {
        var rv *reflect.StructValue
        if v, ok := reflect.NewValue(rhs).(*reflect.PtrValue).Elem().(*reflect.StructValue); ok {
                rv = v
        } else {
                err = os.NewError("args is not a *StructValue")
        }

        if err != nil { return }

        lv := reflect.NewValue(lhs).(*reflect.PtrValue).Elem()
        lt := lv.Type().(*reflect.StructType)
        for i := 0; i < lt.NumField(); i += 1 {
                var lft reflect.StructField

                if lft = lt.Field(i); lft.Anonymous {
                        // FIXME: bind base-struct fields?
                        continue
                }

                if rfv := rv.FieldByName(lft.Name); rfv != nil {
                        lfv := lv.(*reflect.StructValue).FieldByIndex(lft.Index)
                        lfv.SetValue(rfv)
                }
        }
        return
}

func AssignFields(lhs, rhs interface{}) (err os.Error) {
        err = MapFields(lhs, rhs, func(l, r reflect.Value)(nxt bool) {
                l.SetValue(r)
                return true
        })
        return
}

func ForEachField(s interface{}, f func(t *reflect.StructField, v reflect.Value)(nxt bool)) (err os.Error) {
        var ok bool
        var v *reflect.StructValue
        switch p := reflect.NewValue(s).(type) {
        case *reflect.StructValue: v, ok = p, true
        case *reflect.PtrValue: v, ok = p.Elem().(*reflect.StructValue)
        }

        if !ok { err = os.NewError("interface is not *reflect.StructValue"); return }

        t, ok := v.Type().(*reflect.StructType)
        for i := 0; i < t.NumField(); i += 1 {
                ft := t.Field(i)
                fv := v.FieldByIndex(ft.Index)

                if ft.Anonymous {
                        //ForEachField(fv.Interface(), f)
                        continue
                }

                if !f(&ft, fv) { break }
        }
        return
}

// MapFields invoke 'f' for each field of 'lhs' occurs in 'rhs'
func MapFields(lhs, rhs interface{}, f func(lf, rf reflect.Value)(nxt bool)) (err os.Error) {
        lv, ok := reflect.NewValue(lhs).(*reflect.PtrValue).Elem().(*reflect.StructValue)
        if !ok { err = os.NewError("lhs is not *reflect.StructValue"); return }
        lt, ok := lv.Type().(*reflect.StructType)

        rv, ok := reflect.NewValue(rhs).(*reflect.PtrValue).Elem().(*reflect.StructValue)
        if !ok { err = os.NewError("rhs is not *reflect.StructValue"); return }
        rt, ok := rv.Type().(*reflect.StructType)

        for i := 0; i < lt.NumField(); i += 1 {
                var t1, t2 reflect.StructField
                if t1 = lt.Field(i); t1.Anonymous {
                        // FIXME: ...
                        continue
                }

                if t2, ok = rt.FieldByName(t1.Name); !ok {
                        // TODO: handle with Anonymous?
                        continue
                }

                //if rf := rv.FieldByName(t1.Name); rf != nil {
                if rf := rv.FieldByIndex(t2.Index); rf != nil {
                        v := lv.FieldByIndex(t1.Index)
                        if !f(v, rf) { break }
                }
        }
        return
}

func ConvertValue(k reflect.Kind, v reflect.Value) (ov reflect.Value) {
        if k == v.Type().Kind() { ov = v; return }

        if a, ok := v.(reflect.ArrayOrSliceValue); ok {
                if 0 < a.Len() { v = a.Elem(0) } else { return }
        }

        if k == v.Type().Kind() { ov = v; return }

        //s := v.Interface().(string) // TODO: arbitray type
        s := fmt.Sprintf("%v", v.Interface())
        switch k {
        case reflect.Bool:      if o, e := strconv.Atob(s); e == nil { ov = reflect.NewValue(o)
                } else if s == "0" { ov = reflect.NewValue(false)
                } else if s == "1" { ov = reflect.NewValue(true) }
        case reflect.Int:       if o, e := strconv.Atoi(s); e == nil { ov = reflect.NewValue(o) }
        case reflect.Float:     if o, e := strconv.Atof(s); e == nil { ov = reflect.NewValue(o) }
        case reflect.String:    ov = reflect.NewValue(s)
        default:
                fmt.Printf("todo: convert: (%s) %v -> (%s)\n", v.Type().Kind(), v.Interface(), k)
        }
        return
}

func RoughAssignValue(lhs, rhs reflect.Value) (err os.Error) {
        //if p, ok := lhs.(*reflect.InterfaceValue); ok { lhs = p.Elem(); }
        //if p, ok := rhs.(*reflect.InterfaceValue); ok { rhs = p.Elem(); }
        if p, ok := lhs.(*reflect.PtrValue); ok { lhs = p.Elem(); }
        if p, ok := rhs.(*reflect.PtrValue); ok { rhs = p.Elem(); }

        switch lv := lhs.(type) {
        case *reflect.StructValue:
                // Make sure rhs is also StructValue
                //      eg: []struct{} -> struct{}
                rhs = ConvertValue(reflect.Struct, rhs)
                if rhs == nil {
                        err = os.NewError("rhs is not *reflect.StructValue")
                }

                if rv, ok := rhs.(*reflect.StructValue); ok {
                        lt := lv.Type().(*reflect.StructType)
                        //rt := rv.Type().(*reflect.StructType)
                        for i := 0; i < lt.NumField(); i += 1 {
                                ft := lt.Field(i)
                                fv := lv.FieldByIndex(ft.Index)
                                if v := rv.FieldByName(ft.Name); v != nil {
                                        err = RoughAssignValue(fv, v)
                                }
                        }
                } else {
                        err = os.NewError("rhs is not *reflect.StructValue")
                }
        case reflect.ArrayOrSliceValue:
                if rv, ok := rhs.(reflect.ArrayOrSliceValue); ok {
                        //lv = reflect.MakeSlice(lv.Type().(*reflect.SliceType), rv.Len(), rv.Len())
                        for i := 0 ; i < lv.Len() && i < rv.Len(); i += 1 {
                                err = RoughAssignValue(lv.Elem(i), rv.Elem(i))
                        }
                } else {
                        //lv = reflect.MakeSlice(lv.Type().(*reflect.SliceType), 1, 1)
                        if 0 < lv.Len() {
                                err = RoughAssignValue(lv.Elem(0), rv)
                        }
                }
        default:
                if v := ConvertValue(lhs.Type().Kind(), rhs); v != nil {
                        //fmt.Printf("assign: (%s) = (%s) %v\n", lhs.Type().Kind(), rhs.Type().Kind(), rhs.Interface())
                        lhs.SetValue(v)
                } else {
                        fmt.Printf("todo: assign: (%s) = (%s) %v\n", lhs.Type().Kind(), rhs.Type().Kind(), rhs.Interface())
                }
        }
        return
}

func RoughAssign(lhs, rhs interface{}) (err os.Error) {
        defer func() {
                if r := recover(); r != nil {
                        switch e := r.(type) {
                        case os.Error:
                                err = e
                        case string:
                                err = os.NewError(e)
                        default:
                                panic(r)
                        }
                }
        }()

        lv, rv := reflect.NewValue(lhs), reflect.NewValue(rhs)
        return RoughAssignValue(lv, rv)
}

func FieldsToArray(s interface{}) (a []NamedValue, err os.Error) {
        return fieldsToArray(s, false)
}

func FieldsToArrayFlat(s interface{}) (a []NamedValue, err os.Error) {
        return fieldsToArray(s, true)
}

func fieldsToArray(s interface{}, flat bool) (a []NamedValue, err os.Error) {
        v := reflect.NewValue(s)
        if p, ok := v.(*reflect.PtrValue); ok { v = p.Elem() }

        sv, ok := v.(*reflect.StructValue)
        if !ok {
                err = os.NewError("not a *StructValue")
                return
        }
        
        st := sv.Type().(*reflect.StructType)
        if flat {
                a = make([]NamedValue, 0, 2*st.NumField())
                for i := 0; i < st.NumField(); i += 1 {
                        ft := st.Field(i)
                        fv := sv.Field(i) //FieldByIndex(ft.Index)
                        if fsv, ok := fv.(*reflect.StructValue); ok {
                                fa, err := fieldsToArray(fsv.Interface(), true)
                                if err != nil { return }

                                for i := 0; i < len(fa); i += 1 {
                                        fa[i].Name = ft.Name + "." + fa[i].Name
                                }

                                //fmt.Printf("struct: %s = %v\n", ft.Name, fa) //fv.Interface())
                                a = append(a, fa...)
                        } else {
                                a = append(a, NamedValue{ ft.Name, fv.Interface() })
                        }
                }
                //fmt.Printf("%v\n", a);
        } else {
                a = make([]NamedValue, st.NumField())
                for i := 0; i < st.NumField(); i += 1 {
                        ft := st.Field(i)
                        fv := sv.Field(i)
                        a[i] = NamedValue{ ft.Name, fv.Interface() }
                }
        }

        return
}

func RoughAssignQueryResult(iv interface {}, qr web.QueryResult) (err os.Error) {
        v := reflect.NewValue(iv)
        if p, ok := v.(*reflect.PtrValue); ok {
                v = p.Elem()
        } else {
                err = os.NewError("RoughAssignQueryResult: can't assign QueryResult to non ptr value")
                return
        }

        if sv, ok := v.(*reflect.StructValue); ok {
                count := int(qr.GetFieldCount())
                if count <= 0 {
                        os.NewError("RoughAssignQueryResult: no fields in QueryResult")
                        return
                }

                names := make([]string, count)
                for i := 0; i < count; i += 1 {
                        names[i] = qr.GetFieldName(i)
                }

                var row []interface{}
                row, err = qr.FetchRow()
                if err != nil { return }
                if row == nil {
                        err = os.NewError("RoughAssignQueryResult: no rows in QueryResult")
                        return
                }

                err = roughAssignQueryResultRow(sv, names, row)
        } else {
                err = os.NewError("RoughAssignQueryResult: can't assign QueryResult to non struct value")
        }

        return
}

func roughAssignQueryResultRow(sv *reflect.StructValue, names []string, row []interface{}) (err os.Error) {
        for i := 0; i < len(names); i += 1 {
                name := names[i]

                fsv := sv
                for {
                        if p := strings.Index(name, "$"); 0 < p {
                                fv := fsv.FieldByName(name[0:p])
                                if fv == nil {
                                        err = os.NewError("RoughAssignQueryResult: no field by name '"+name[0:p]+"'")
                                        return // or 'continue' to ignore this?
                                }

                                if v, ok := fv.(*reflect.StructValue); ok {
                                        name = name[p+1:len(name)]
                                        fsv = v
                                } else {
                                        err = os.NewError("RoughAssignQueryResult: field is not struct value")
                                        return
                                }
                        } else { break }
                }//for

                fv := fsv.FieldByName(name)
                if fv == nil {
                        err = os.NewError(fmt.Sprintf("RoughAssignQueryResult: no filed by name '%s'", names[i]))
                        return
                }

                err = RoughAssignValue(fv, reflect.NewValue(row[i]))
                if err != nil {
                        return
                }
        }
        //fmt.Printf("%v\n", iv)
        return
}

