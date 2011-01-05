package eBay

import (
        "testing"
        "reflect"
        //"fmt"
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
}

func TestUtilAssignValue(t *testing.T) {
}

func TestUtilCopyFields(t *testing.T) {
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

                //err := CopyFields(s1, s2)
                err := RoughAssign(s1, s2)
                if err != nil { t.Errorf("CopyFields: %v", err); return }
                if s1.A != 1 { t.Errorf("CopyFields:A: %v", s1); return }
                if s1.B != 0.1 { t.Errorf("CopyFields:B: %v", s1); return }
                if s1.C != "foo" { t.Errorf("CopyFields:C: %v", s1); return }
                if s1.D != true { t.Errorf("CopyFields:D: %v", s1); return }
                if s1.E != "foo" { t.Errorf("CopyFields:E: %v", s1); return }
                if s1.S.A != 100 { t.Errorf("CopyFields:S.A: %v", s1); return }
                if s1.S.B != "foo" { t.Errorf("CopyFields:S.B: %v", s1); return }
                if s1.S.C != true { t.Errorf("CopyFields:S.C: %v", s1); return }
                if s1.S.Z != "foobar" { t.Errorf("CopyFields:S.Z: %v", s1); return }
                if s1.Z != 1000 { t.Errorf("CopyFields:Z: %v", s1); return }
        }
}
