//Package server implements the http frontend
package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/terranodo/tegola"
	"github.com/terranodo/tegola/basic"
	"github.com/terranodo/tegola/mvt"
	"github.com/terranodo/tegola/provider/postgis"
)

const (
	//MaxTileSize is	500k
	MaxTileSize = 500000
	//MaxZoom is the suggested max by Slippy Map Tilenames spec
	MaxZoom = 18
)

func init() {
	config := postgis.Config{
		Host:     "localhost",
		Port:     5432,
		Database: "gdey",
		User:     "gdey",
		Layers: map[string]string{
			"landuse": "gis.zoning_base_3857",
		},
	}
	var err error
	postgisProvider, err = postgis.NewProvider(config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create a new provider. %v", err))
	}

}

//	creates a debug layer with z/x/y encoded as a point
func debugLayer(tile tegola.Tile) *mvt.Layer {
	//	create a line
	line1 := &basic.Line{
		basic.Point{0, 0},
		basic.Point{4096, 0},
		basic.Point{4096, 4096},
		basic.Point{0, 4096},
		basic.Point{0, 0},
	}

	//	tile outlines
	outline := mvt.Feature{
		Tags: map[string]interface{}{
			"type": "debug_outline",
		},
		Geometry: line1,
	}

	//	middle of tile
	point1 := &basic.Point{2048, 2048}

	//	new feature
	zxy := mvt.Feature{
		Tags: map[string]interface{}{
			"type":    "debug_text",
			"name_en": fmt.Sprintf("Z:%v, X:%v, Y:%v", tile.Z, tile.X, tile.Y),
		},
		Geometry: point1,
	}

	layer := mvt.Layer{
		Name: "debug",
	}

	layer.AddFeatures(zxy, outline)

	return &layer
}

var postgisProvider *postgis.Provider

//	URI scheme: /maps/:map_id/:z/:x/:y
//		map_id - id in the config file with an accompanying data source
//		z, x, y - tile coordinates as described in the Slippy Map Tilenames specification
//			z - zoom level
//			x - row
//			y - column
func handleZXY(w http.ResponseWriter, r *http.Request) {
	//	check http verb
	switch r.Method {
	//	preflight check for CORS request
	case "OPTIONS":
		//	TODO: how configurable do we want the CORS policy to be?
		//	set CORS header
		w.Header().Add("Access-Control-Allow-Origin", "*")

		//	options call does not have a body
		w.Write(nil)
		return
	//	tile request
	case "GET":
		//	pop off URI prefix
		uri := r.URL.Path[len("/maps/"):]

		//	break apart our URI
		uriParts := strings.Split(uri, "/")

		//	check that we have the correct number of arguments in our URI
		if len(uriParts) != 4 {
			http.Error(w, "uri requires four params: /:map_id/:z/:x/:y", http.StatusBadRequest)
			return
		}

		//	TODO: look up layer data source provided by config

		//	trim the "y" param in the url in case it has an extension
		yparts := strings.Split(uriParts[3], ".")
		uriParts[3] = yparts[0]

		//	parse our URL vals to ints
		z, err := strconv.Atoi(uriParts[1])
		if err != nil {
			http.Error(w, "invalid z value: "+uriParts[1], http.StatusBadRequest)
			return
		}

		x, err := strconv.Atoi(uriParts[2])
		if err != nil {
			http.Error(w, "invalid x value: "+uriParts[2], http.StatusBadRequest)
			return
		}

		y, err := strconv.Atoi(uriParts[3])
		if err != nil {
			http.Error(w, "invalid y value: "+uriParts[3], http.StatusBadRequest)
			return
		}

		//	new tile
		tile := tegola.Tile{
			Z: z,
			X: x,
			Y: y,
		}

		//	generate a tile
		var mvtTile mvt.Tile
		var pbyte []byte

		//	check that our request is below max zoom
		if tile.Z < MaxZoom {
			//	fetch requested layer from our data provider
			mvtLayer, err := postgisProvider.MVTLayer("landuse", tile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error Getting MVTLayer: %v", err.Error()), http.StatusBadRequest)
				return
			}

			//	add layers
			mvtTile.AddLayers(mvtLayer)

			//	add debug layer
			debugLayer := debugLayer(tile)
			mvtTile.AddLayers(debugLayer)
		}

		vtile, err := mvtTile.VTile(tile.Extent())
		if err != nil {
			http.Error(w, fmt.Sprintf("Error Getting VTile: %v", err.Error()), http.StatusBadRequest)
			return
		}

		//	marshal our tile into a protocol buffer
		pbyte, err = proto.Marshal(vtile)
		if err != nil {
			http.Error(w, "error marshalling tile", http.StatusInternalServerError)
			return
		}

		//	check for tile size warnings
		if len(pbyte) > MaxTileSize {
			log.Printf("tile is rather large - %v", len(pbyte))
		}

		//	TODO: how configurable do we want the CORS policy to be?
		//	set CORS header
		w.Header().Add("Access-Control-Allow-Origin", "*")

		//	mimetype for protocol buffers
		w.Header().Add("Content-Type", "application/x-protobuf")

		w.Write(pbyte)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}