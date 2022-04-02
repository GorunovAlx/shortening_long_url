package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const UserID = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func TestGenerateShortLink(t *testing.T) {
	shortenedLinks := map[string]string{
		"https://bitfieldconsulting.com/golang/slower":                                                                  "UovDy88t",
		"https://medium.com/scaled-agile-framework/exploring-key-elements-of-spotifys-agile-scaling-model-471d2a23d7ea": "NFf4EeL2",
		"https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-1-2dbb96ce079a":      "hKKSoByB",
	}

	id, err := GetUserID(UserID)
	require.NoError(t, err)

	for initialLink, shortenedLink := range shortenedLinks {
		shortLink, err := GenerateShortLink(initialLink, id)
		require.NoError(t, err)
		assert.Equal(t, shortLink, shortenedLink)
	}
}
