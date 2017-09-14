package mvt

import "github.com/rikardNL/tegola"

//Provider is the mechinism by which the system talks to different data providers.
type Provider interface {
	// MVTLayer returns a layer object based
	MVTLayer(layerName string, tile tegola.Tile, tags map[string]interface{}) (*Layer, error)
	// LayerNames returns a list of layer name the Provider knows about.
	LayerNames() []string
}
