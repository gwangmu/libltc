package test

import (
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"path/filepath"
	
	"libltc/warning"
)

var curT *testing.T

func Initialize(t *testing.T, cwd string) {
	curT = t
	t.Chdir(cwd)
}

type IPredicate interface {
	Test(any, any) bool
	String() string
}

type EQ struct {}
func (_ EQ) Test(got any, want any) bool {
    vgot := reflect.ValueOf(got)
    if vgot.Comparable() {
        vwant := reflect.ValueOf(want)
        return vgot.Equal(vwant)
    }
	return false
}
func (_ EQ) String() string {
	return "=="
}

type NE struct {}
func (_ NE) Test(got any, want any) bool {
    vgot := reflect.ValueOf(got)
    if vgot.Comparable() {
        vwant := reflect.ValueOf(want)
        return !vgot.Equal(vwant)
    }
	return false
}
func (_ NE) String() string {
	return "!="
}

type Contains struct {}
func (_ Contains) Test(got any, want any) bool {
	if sgot, ok := got.(string); ok {
		if swant, ok := want.(string); ok {
			return strings.Contains(sgot, swant)
		}
	}
	return false
}
func (_ Contains) String() string {
	return "⊇"
}

type WarnContains struct {}
func (_ WarnContains) Test(got any, want any) bool {
	if wgot, ok := got.(warning.Warnings); ok {
		if swant, ok := want.(string); ok {
			re := regexp.MustCompile(swant)
			found := false
			for _, w := range wgot {
				if re.MatchString(w.GetDesc()) {
					found = true
					break
				}
			}
			return found
		}
	}
	return false
}
func (_ WarnContains) String() string {
	return "⊇"
}

func Assert[Pred IPredicate](got any, want any) {
	_, fname, fno, ok := runtime.Caller(2)
	if ok {
		fname = filepath.Base(fname)
	} else {
		fname = "?"
		fno = 0
	}

	fcg := "#v"
	if _, ok := got.(string); ok {
		fcg = "q"
	} else if reflect.ValueOf(got).Kind() == reflect.Ptr {
		fcg = "p"
	}

	fcw := "#v"
	if _, ok := want.(string); ok {
		fcw = "q"
	} else if reflect.ValueOf(want).Kind() == reflect.Ptr {
		fcw = "p"
	}

	var pred Pred
	if !pred.Test(got, want) {
		curT.Fatalf("Assertion failed: %"+fcg+" %s %"+fcw+" (%s:%d)", got, pred.String(), want, fname, fno)
	} else {
		curT.Logf("Assertion succeeded: %"+fcg+" %s %"+fcw+" (%s:%d)", got, pred.String(), want, fname, fno)
	}
}

func AssertEQ[T comparable](got T, want T) {
	Assert[EQ](got, want)
}

func AssertNE[T comparable](got T, want T) {
	Assert[NE](got, want)
}

func AssertContains(got string, want string) {
	Assert[Contains](got, want)
}

func AssertWarnContains(gotw warning.Warnings, want string) {
	Assert[WarnContains](gotw, want)
}

func PrintAllWarnings(warns warning.Warnings) {
	for _, w := range warns {
		curT.Log(" - WARN: " + w.GetDesc())
	}
}
