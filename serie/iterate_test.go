package serie_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xinzf/datatable/serie"
)

func TestIterate(t *testing.T) {
	xs := []float64{
		32.32, 56.98, 21.52, 44.32,
		55.63, 13.75, 43.47, 43.34,
		12.34,
	}

	s := serie.Float64(xs)

	index := 0
	for it := s.Iterator(); it.Next(); {
		assert.Equal(t, xs[index], it.Current())
		index++
	}
}
