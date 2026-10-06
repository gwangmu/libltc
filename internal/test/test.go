package test

import (
	"reflect"
	"runtime"
	"testing"
	"path/filepath"
)

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

func Assert[Pred IPredicate, T comparable](t *testing.T, got T, want T) {
	_, fname, fno, ok := runtime.Caller(2)
	if ok {
		fname = filepath.Base(fname)
	} else {
		fname = "?"
		fno = 0
	}

	var pred Pred
	if !pred.Test(got, want) {
		t.Fatalf("Assertion failed: %v %s %v (%s:%d)", got, pred.String(), want, fname, fno)
	} else {
		t.Logf("Assertion succeeded: %v %s %v (%s:%d)", got, pred.String(), want, fname, fno)
	}
}

func AssertEQ[T comparable](t *testing.T, got T, want T) {
	Assert[EQ](t, got, want)
}

func AssertNE[T comparable](t *testing.T, got T, want T) {
	Assert[NE](t, got, want)
}
