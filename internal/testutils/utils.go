package testutils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"reflect"
	"time"

	"github.com/StampWallet/backend/internal/database"
	"github.com/gin-gonic/gin"
)

func MatchEntities(matcher interface{}, Obj interface{}) bool {
	o := reflect.ValueOf(Obj)
	m := reflect.ValueOf(matcher)
	if o.Kind() == reflect.Pointer {
		if o.IsNil() {
			return false
		} else {
			return MatchEntities(matcher, o.Elem().Interface())
		}
	} else if m.Kind() == reflect.Pointer {
		return MatchEntities(m.Elem().Interface(), o)
	} else {
		mt := reflect.TypeOf(matcher)
		for i := 0; i < mt.NumField(); i++ {
			mtf := mt.Field(i)
			of := o.FieldByName(mtf.Name)
			mf := m.FieldByName(mtf.Name)
			if (mf.Kind() == reflect.Pointer || mf.Kind() == reflect.Interface) && !mf.IsNil() && !of.Equal(mf.Elem()) {
				return false
			} else if (mf.Kind() == reflect.Array || mf.Kind() == reflect.Slice) && !reflect.DeepEqual(of, mf) {
				return false
			} else if mf.Kind() == reflect.Struct && !MatchEntities(of, mf) {
				return false
			} else if !of.Equal(mf) {
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
	uint64 | uint | string | bool | time.Time | database.GPSCoordinates
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

type TestContextBuilder struct{ Context *gin.Context }

func NewTestContextBuilder(w *httptest.ResponseRecorder) *TestContextBuilder {
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request, _ = http.NewRequest("", "", nil)
	return &TestContextBuilder{Context: ctx}
}

func TestContextCopy(c *gin.Context) *TestContextBuilder {
	return &TestContextBuilder{Context: c.Copy()}
}

func (tc *TestContextBuilder) SetUser(u *database.User) *TestContextBuilder {
	tc.Context.Set("user", u)
	return tc
}

func (tc *TestContextBuilder) SetDefaultUser() *TestContextBuilder {
	return tc.SetUser(GetDefaultUser())
}

func (tc *TestContextBuilder) SetUrl(argUrl *url.URL) *TestContextBuilder {
	tc.Context.Request.URL = argUrl
	return tc
}

func (tc *TestContextBuilder) SetDefaultUrl() *TestContextBuilder {
	url, _ := url.Parse("localhost")
	return tc.SetUrl(url)
}

func (tc *TestContextBuilder) SetEndpoint(endpointPath string) *TestContextBuilder {
	rUrl := tc.Context.Request.URL
	rUrl.Path = path.Join(rUrl.Path, endpointPath)
	return tc
}

func (tc *TestContextBuilder) SetMethod(method string) *TestContextBuilder {
	tc.Context.Request.Method = method
	return tc
}

func (tc *TestContextBuilder) SetHeader(headerKey string, headerValue string) *TestContextBuilder {
	tc.Context.Request.Header.Set(headerKey, headerValue)
	return tc
}

func (tc *TestContextBuilder) SetToken(token string) *TestContextBuilder {
	return tc.SetHeader("Authorization", "Bearer "+token)
}

func (tc *TestContextBuilder) SetDefaultToken() *TestContextBuilder {
	return tc.SetToken("012346789:ZWVnaDhhZWg4bGVpbDJhaXBlaW5nZWViNWFpU2hlaGUK")
}

func (tc *TestContextBuilder) SetBody(jsonBytes []byte) *TestContextBuilder {
	tc.Context.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
	return tc
}

func ExtractResponse[T any](w *httptest.ResponseRecorder) (*T, int, error) {
	bodyBytes := w.Body.Bytes()
	bodyPtr := new(T)
	err := json.Unmarshal(bodyBytes, bodyPtr)
	return bodyPtr, w.Code, err
}
