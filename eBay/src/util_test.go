package eBay

import (
        "testing"
        "reflect"
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
