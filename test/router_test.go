package test

import (
	evoli "github.com/evolidev/evoli/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBasic(t *testing.T) {
	router := evoli.NewRouter()

	methodTable := make(map[string]func(path string, f func() string))
	methodTable[http.MethodGet] = router.Get
	methodTable[http.MethodPost] = router.Post
	methodTable[http.MethodPut] = router.Put
	methodTable[http.MethodPatch] = router.Patch
	methodTable[http.MethodDelete] = router.Delete
	methodTable[http.MethodHead] = router.Head
	methodTable[http.MethodOptions] = router.Options

	for method, fn := range methodTable {
		t.Run("Basic "+method+" route should have the returned string of the callback in the body", func(t *testing.T) {
			path := "/" + strings.ToLower(method)

			fn(path, func() string {
				return "hello-world"
			})

			req, err := http.NewRequest(method, path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			assert.Exactly(t, "hello-world", rr.Body.String())
		})
	}

}
