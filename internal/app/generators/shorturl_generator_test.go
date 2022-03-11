package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortLink(t *testing.T) {
	initialLink1 := "https://bitfieldconsulting.com/golang/slower"
	shortLink1 := GenerateShortLink(initialLink1)

	initialLink2 := "https://medium.com/scaled-agile-framework/exploring-key-elements-of-spotifys-agile-scaling-model-471d2a23d7ea"
	shortLink2 := GenerateShortLink(initialLink2)

	initialLink3 := "https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-1-2dbb96ce079a"
	shortLink3 := GenerateShortLink(initialLink3)

	assert.Equal(t, shortLink1, "ZhvygENi")
	assert.Equal(t, shortLink2, "Cb68RzKG")
	assert.Equal(t, shortLink3, "jo9sYrzF")
}
