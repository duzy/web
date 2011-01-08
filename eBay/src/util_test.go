package eBay

import (
        "testing"
        "reflect"
        "fmt"
)

func TestUtilMapFields(t *testing.T) {
        {
                lhs := &struct{
                        A string
                        B int
                        C bool
                }{}

                rhs := &struct{
                        A string
                        B int
                }{ "rhs.A", 2 }

                rhs2 := &struct{
                        A string
                        C bool
                }{ "rhs2.A", true }

                err := MapFields(lhs, rhs, func(lf, rf reflect.Value)(nxt bool){
                        lf.SetValue(rf)
                        return true
                })

                if err != nil { t.Error(err); return }
                if lhs.A != "rhs.A" { t.Errorf("lhs.A: %v", lhs.A); return }
                if lhs.B != 2 { t.Errorf("lhs.B: %v", lhs.B); return }
                if lhs.C != false { t.Errorf("lhs.C: %v", lhs.C); return }

                err = AssignFields(lhs, rhs2)
                if err != nil { t.Error(err); return }
                if lhs.A != "rhs2.A" { t.Errorf("lhs.A: %v", lhs.A); return }
                if lhs.B != 2 { t.Errorf("lhs.B: %v", lhs.B); return }
                if lhs.C != true { t.Errorf("lhs.C: %v", lhs.C); return }

                i := 5
                rhs3 := &i
                err = AssignFields(lhs, rhs3)
                if err == nil { t.Error("accept int->struct assignment(mistake)"); return }
        }
}

func TestUtilConvertValue(t *testing.T) {
        // Scalar to scalar convertion:

        if v := ConvertValue(reflect.Int, reflect.NewValue("100")); v != nil {
                a, ok := v.Interface().(int)
                if !ok || a != 100 { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }

        if v := ConvertValue(reflect.Float, reflect.NewValue("0.1")); v != nil {
                a, ok := v.Interface().(float)
                if !ok || a != 0.1 { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> float: %v", v); return }

        if v := ConvertValue(reflect.Bool, reflect.NewValue("true")); v != nil {
                a, ok := v.Interface().(bool)
                if !ok || a != true { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> float: %v", v); return }

        if v := ConvertValue(reflect.String, reflect.NewValue("foo")); v != nil {
                a, ok := v.Interface().(string)
                if !ok || a != "foo" { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> string: %v", v); return }

        // Array to scalar convertion:

        if v := ConvertValue(reflect.Int, reflect.NewValue([]string{"100"})); v != nil {
                a, ok := v.Interface().(int)
                if !ok || a != 100 { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }

        if v := ConvertValue(reflect.Float, reflect.NewValue([]string{"0.1"})); v != nil {
                a, ok := v.Interface().(float)
                if !ok || a != 0.1 { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> float: %v", v); return }

        if v := ConvertValue(reflect.Bool, reflect.NewValue([]string{"true"})); v != nil {
                a, ok := v.Interface().(bool)
                if !ok || a != true { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> float: %v", v); return }

        if v := ConvertValue(reflect.String, reflect.NewValue([]string{"foo"})); v != nil {
                a, ok := v.Interface().(string)
                if !ok || a != "foo" { t.Errorf("ConvertValue: failed: string -> int: %v", v); return }
        } else { t.Errorf("ConvertValue: failed: string -> string: %v", v); return }
}

func TestUtilRoughAssign(t *testing.T) {
        {
                s1 := &struct {
                        A int
                        B float
                        C string
                        D bool
                        E string // array -> scalar
                        S struct{
                                A int
                                B string
                                C bool
                                Z string // un-touched
                        }
                        Z int // un-touched
                }{ Z:1000 }
                s1.S.Z = "foobar"

                s2 := &struct {
                        A string
                        B string
                        C string
                        D string
                        E []string
                        S struct{
                                A string
                                B string
                                C string
                        }
                }{
                A: "1",
                B: "0.1",
                C: "foo",
                D: "true",
                E: []string{ "foo", "bar" },
                //S: 
                }
                s2.S.A = "100"
                s2.S.B = "foo"
                s2.S.C = "true"

                err := RoughAssign(s1, s2)
                if err != nil { t.Errorf("RoughAssign: %v", err); return }
                if s1.A != 1 { t.Errorf("RoughAssign:A: %v", s1); return }
                if s1.B != 0.1 { t.Errorf("RoughAssign:B: %v", s1); return }
                if s1.C != "foo" { t.Errorf("RoughAssign:C: %v", s1); return }
                if s1.D != true { t.Errorf("RoughAssign:D: %v", s1); return }
                if s1.E != "foo" { t.Errorf("RoughAssign:E: %v", s1); return }
                if s1.S.A != 100 { t.Errorf("RoughAssign:S.A: %v", s1); return }
                if s1.S.B != "foo" { t.Errorf("RoughAssign:S.B: %v", s1); return }
                if s1.S.C != true { t.Errorf("RoughAssign:S.C: %v", s1); return }
                if s1.S.Z != "foobar" { t.Errorf("RoughAssign:S.Z: %v", s1); return }
                if s1.Z != 1000 { t.Errorf("RoughAssign:Z: %v", s1); return }
        }
        {
                s0 := &struct {
                        A struct {
                                Aa struct{ Aaa int }
                        }
                }{}

                s1 := &struct {
                        A struct {
                                Aa [1]struct{ Aaa int }
                        }
                }{}

                s2 := &struct {
                        A [1]struct{
                                Aa [1]struct { Aaa string }
                        }
                }{}
                s2.A[0].Aa[0].Aaa = "11"

                err := RoughAssign(s0, s2)
                if err != nil { t.Errorf("RoughAssign: %v", err); return }
                if s0.A.Aa.Aaa != 11 { t.Errorf("RoughAssign: %v", s0); return }

                err = RoughAssign(s1, s2)
                if err != nil { t.Errorf("RoughAssign: %v", err); return }
                if s1.A.Aa[0].Aaa != 11 { t.Errorf("RoughAssign: %v", s1); return }
        }
}

func TestUtilFieldsToArray(t *testing.T) {
        s := &struct{
                A int
                B string
                C struct{
                        C1 string
                        C2 bool
                        C3 struct {
                                C3_a string
                                C3_b string
                        }
                }
        }{}

        s.A = 100
        s.B = "foo"
        s.C.C1 = "foo"
        s.C.C2 = true
        s.C.C3.C3_a = "foo"
        s.C.C3.C3_b = "foo"
        
        {
                a, err := FieldsToArray(s)
                if err != nil {
                        t.Errorf("FieldsToArrayFlat: %v", err)
                        return
                }

                //fmt.Printf("%v\n", a)
                if a[0].Name != "A" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[1].Name != "B" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[2].Name != "C" { t.Errorf("FieldsToArrayFlat: %v", a); return }

                if a[0].Value != 100 { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[1].Value != "foo" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if fmt.Sprint(a[2].Value) != "{foo true {foo foo}}" {
                        t.Errorf("FieldsToArrayFlat: %v", fmt.Sprint(a[2].Value));
                        return
                }
        }

        {
                a, err := FieldsToArrayFlat(s)
                if err != nil {
                        t.Errorf("FieldsToArrayFlat: %v", err)
                        return
                }

                //fmt.Printf("%v\n", a)
                if a[0].Name != "A" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[1].Name != "B" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[2].Name != "C.C1" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[3].Name != "C.C2" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[4].Name != "C.C3.C3_a" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[5].Name != "C.C3.C3_b" { t.Errorf("FieldsToArrayFlat: %v", a); return }

                if a[0].Value != 100 { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[1].Value != "foo" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[2].Value != "foo" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[3].Value != true { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[4].Value != "foo" { t.Errorf("FieldsToArrayFlat: %v", a); return }
                if a[5].Value != "foo" { t.Errorf("FieldsToArrayFlat: %v", a); return }
        }
}
