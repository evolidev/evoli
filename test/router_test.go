package test

import (
	evoli "github.com/evolidev/evoli/pkg"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasic(t *testing.T) {
	t.Run("Basic Get route should have the returned string of the callback in the body", func(t *testing.T) {
		router := evoli.NewRouter()
		router.Get("/test", func() string {
			return "hello-world"
		})

		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Exactly(t, "hello-world", rr.Body.String())
	})
}
