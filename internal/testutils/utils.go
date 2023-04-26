package testutils

import (
	"database/sql/driver"
	"reflect"
	"time"
)

func MatchEntities(matcher interface{}, Obj interface{}) bool {
	o := reflect.ValueOf(Obj)
	m := reflect.ValueOf(matcher)
	if o.Kind() == reflect.Pointer {
		return MatchEntities(matcher, o.Elem().Interface())
	} else if m.Kind() == reflect.Pointer {
		return MatchEntities(m.Elem().Interface(), o)
	} else {
		mt := reflect.TypeOf(matcher)
		for i := 0; i < mt.NumField(); i++ {
			mtf := mt.Field(i)
			of := o.FieldByName(mtf.Name)
			mf := m.FieldByName(mtf.Name)
			if !mf.IsNil() && !of.Equal(mf.Elem()) {
				return false
			}
		}
		return true
	}
}

type StructMatcher struct {
	Obj interface{}
}

func (matcher StructMatcher) Matches(x interface{}) bool {
	return MatchEntities(matcher.Obj, x)
}

func (StructMatcher) String() string {
	return "StructMatcher"
}

type TimeGreaterThanNow struct {
	Time time.Time
}

func (matcher TimeGreaterThanNow) Matches(x interface{}) bool {
	return matcher.Time.Before(x.(time.Time))
}

func (TimeGreaterThanNow) String() string {
	return "TimeGreaterThanNow"
}

type Copyable interface {
	uint64 | uint | string | bool | time.Time
}

func Ptr[T Copyable](s T) *T {
	return &s
}

type Anything struct{}

func (Anything) Match(v driver.Value) bool {
	return true
}

func ReturnArg(arg interface{}) interface{} {
	return arg
}
