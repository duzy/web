package eBay

import (
        "os"
        "reflect"
)

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
