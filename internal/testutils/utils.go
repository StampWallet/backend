package testutils

import (
	"bytes"
	"io"
	"os"
	"path"
	"reflect"
	"time"

	"database/sql/driver"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/StampWallet/backend/internal/database"
	"github.com/gin-gonic/gin"
)

// Recursively compares matcher with obj. Only keys present in matcher are compared
// Mostly broken. Should only be used in tests.
// Note: for exact equality reflect.DeepEqual should be used instead
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

type TimeJustBeforeNow struct {
}

func (matcher TimeJustBeforeNow) Matches(x interface{}) bool {
	return time.Now().After(x.(time.Time)) && time.Now().Add(-5*time.Minute).Before(x.(time.Time))
}

func (TimeJustBeforeNow) String() string {
	return "TimeJustBeforeNow"
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

func TestFileReader(filename string) io.Reader {
	buf := new(bytes.Buffer)
	mw := multipart.NewWriter(buf)
	w, _ := mw.CreateFormFile("file", "test")

	file, _ := os.Open(filename)
	io.Copy(w, file)

	file.Close()
	mw.Close()
	return buf
}

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

func (tc *TestContextBuilder) AddQueryParam(paramName string, paramValue string) *TestContextBuilder {
	query := tc.Context.Request.URL.Query()
	query.Add(paramName, paramValue)
	tc.Context.Request.URL.RawQuery = query.Encode()
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

// Overwrites request body
// Lifted from https://github.com/gin-gonic/gin/blob/master/context_test.go > TestContextFormFile
func (tc *TestContextBuilder) SetTestFile(filename string) *TestContextBuilder {
	r := TestFileReader(filename)
	tc.Context.Request.Body = io.NopCloser(r)
	return tc.SetHeader("Content-Type", "multipart/form-data")
}

func (tc *TestContextBuilder) AttachTestPng() *TestContextBuilder {
	return tc.SetTestFile("resources/test.png")
}

func (tc *TestContextBuilder) AttachTestJpeg() *TestContextBuilder {
	return tc.SetTestFile("resources/test.jpeg")
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

func (tc *TestContextBuilder) SetParam(key string, value string) *TestContextBuilder {
	tc.Context.AddParam(key, value)
	return tc
}

func ExtractResponse[T any](w *httptest.ResponseRecorder) (*T, int, error) {
	bodyBytes := w.Body.Bytes()
	bodyPtr := new(T)
	err := json.Unmarshal(bodyBytes, bodyPtr)
	return bodyPtr, w.Code, err
}

func TimeJustAroundNow(x time.Time) bool {
	return time.Now().After(x) && time.Now().Add(-5*time.Minute).Before(x)
}
