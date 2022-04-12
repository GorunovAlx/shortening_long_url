package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const UserID = 2351092691

func TestGenerateShortLink(t *testing.T) {
	shortenedLinks := map[string]string{
		"https://bitfieldconsulting.com/golang/slower":                                                                  "duiLQBQW",
		"https://medium.com/scaled-agile-framework/exploring-key-elements-of-spotifys-agile-scaling-model-471d2a23d7ea": "com6ngTL",
		"https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-1-2dbb96ce079a":      "4tEZyWRC",
	}

	for initialLink, shortenedLink := range shortenedLinks {
		shortLink, err := GenerateShortLink(initialLink, UserID)
		require.NoError(t, err)
		assert.Equal(t, shortLink, shortenedLink)
	}
}
