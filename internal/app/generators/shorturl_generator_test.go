package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateShortLink(t *testing.T) {
	shortenedLinks := map[string]string{
		"https://bitfieldconsulting.com/golang/slower":                                                                  "ZhvygENi",
		"https://medium.com/scaled-agile-framework/exploring-key-elements-of-spotifys-agile-scaling-model-471d2a23d7ea": "Cb68RzKG",
		"https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-1-2dbb96ce079a":      "jo9sYrzF",
	}

	for initialLink, shortenedLink := range shortenedLinks {
		shortLink, err := GenerateShortLink(initialLink)
		require.NoError(t, err)
		assert.Equal(t, shortLink, shortenedLink)
	}
}
