package main

import (
	"os"
	"testing"

	"github.com/kataras/iris/httptest"
)

// $ go test -v
func TestNewApp(t *testing.T) {
	os.Remove("./foo_test.db")

	appConfig := AppConfig{DB: createDb("foo_test.db")}
	app := newIrisApp(&appConfig)
	e := httptest.New(t, app)

	postData := `[{"value": 12, "lat": 55.1, "lng": 22.321}]`

	formData := make(map[string]string)
	formData["reading"] = postData

	e.POST("/readings/abcd123").WithBasicAuth("user", "password").WithForm(formData).Expect().Status(httptest.StatusOK).Body().Equal("ok")
	e.POST("/readings/abcd123").WithBasicAuth("user", "zen").WithForm(formData).Expect().Status(httptest.StatusUnauthorized)

	e.GET("/readings/abcd123").Expect().Status(httptest.StatusOK).Body().Equal("ok")
}
