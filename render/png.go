package render

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	geo "github.com/therealfakemoot/genesis/geo"
	Q "github.com/therealfakemoot/go-quantize"
)

func matchColor(point float64, d Q.Domain) (c color.Color) {
	colorSpace := Q.Domain{
		Min: 0,
		Max: 255,
	}
	// normalized := uint8(colorSpace.QuantizePoint(point))
	normalized := uint8(Q.Quantize(point, d, colorSpace))
	log.WithFields(log.Fields{
		"point": point,
		"color": normalized,
	}).Debug("matching color")

	return color.NRGBA{normalized, normalized, normalized, 255}
}

func GeneratePNG(m geo.Map) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, m.Width, m.Height))

	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			point := m.Points[x][y]
			c := matchColor(point, m.Domain)
			img.Set(x, y, c)
		}
	}

	return img
}

func ServePNG(w http.ResponseWriter, m geo.Map) {
	var err error
	buffer := new(bytes.Buffer)

	i := GeneratePNG(m)

	w.Header().Set("Content-type", "image/png")

	w.Header().Set("Content-Disposition", `inline;filename="map"`)
	err = png.Encode(buffer, i)
	if err != nil {
		log.WithError(err).Error("image encoding failure")

		e := struct {
			Error string
		}{Error: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(e)
		return
	}

	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	_, err = w.Write(buffer.Bytes())
	if err != nil {
		log.WithError(err).Error("response write failure")
	}

}
