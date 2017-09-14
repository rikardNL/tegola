package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dimfeld/httptreemux"
	"github.com/rikardNL/tegola/mvt"
	"github.com/rikardNL/tegola/server"
)

func TestHandleMapLayerZXY(t *testing.T) {
	//	setup a new provider
	testcases := []struct {
		handler    http.Handler
		uri        string
		uriPattern string
		reqMethod  string
		expected   mvt.Tile
	}{
		{
			handler:    server.HandleMapLayerZXY{},
			uri:        "/maps/test-map/test-layer/1/2/3.pbf",
			uriPattern: "/maps/:map_name/:layer_name/:z/:x/:y",
			reqMethod:  "GET",
			expected:   mvt.Tile{},
		},
	}

	for i, test := range testcases {
		var err error

		//	setup a new router. this handles parsing our URL wildcards (i.e. :map_name, :z, :x, :y)
		router := httptreemux.New()
		//	setup a new router group
		group := router.NewGroup("/")
		group.UsingContext().Handler(test.reqMethod, test.uriPattern, server.HandleMapZXY{})

		r, err := http.NewRequest(test.reqMethod, test.uri, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Errorf("Failed test %v. handler returned wrong status code: got (%v) expected (%v)", i, w.Code, http.StatusOK)
		}
	}
}
