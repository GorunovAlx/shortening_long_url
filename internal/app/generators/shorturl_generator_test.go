package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortLink(t *testing.T) {
	initialLink_1 := "https://bitfieldconsulting.com/golang/slower"
	shortLink_1 := GenerateShortLink(initialLink_1)

	initialLink_2 := "https://medium.com/scaled-agile-framework/exploring-key-elements-of-spotifys-agile-scaling-model-471d2a23d7ea"
	shortLink_2 := GenerateShortLink(initialLink_2)

	initialLink_3 := "https://medium.com/capital-one-tech/doing-well-by-doing-bad-writing-bad-code-with-go-part-1-2dbb96ce079a"
	shortLink_3 := GenerateShortLink(initialLink_3)

	assert.Equal(t, shortLink_1, "ZhvygENi")
	assert.Equal(t, shortLink_2, "Cb68RzKG")
	assert.Equal(t, shortLink_3, "jo9sYrzF")
}
