package imgix

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMe_main(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))
	actual := ub.CreateURL("path/to/image.jpg")
	expected := "https://demo.imgix.net/path/to/image.jpg"
	assert.Equal(t, expected, actual)
}

func TestReadMe_usageWithParams(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))
	actual := ub.CreateURL("path/to/image.jpg", Param("w", "320"), Param("auto", "format", "compress"))
	expected := "https://demo.imgix.net/path/to/image.jpg?auto=format%2Ccompress&w=320"
	assert.Equal(t, expected, actual)
}

func TestReadMe_SecuredURLUsage(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithToken("MYT0KEN"), WithLibParam(false))
	expected := "https://demo.imgix.net/path/to/image.jpg?s=c8bd1807209f7f1d96dd7123f92febb4"
	actual := ub.CreateURL("path/to/image.jpg")
	assert.Equal(t, expected, actual)
}

func TestReadMe_usageSrcsetGeneration(t *testing.T) {
	ub := NewURLBuilder("demos.imgix.net", WithToken("foo123"))
	srcset := ub.CreateSrcset("image.png", []IxParam{})
	splitSrcset := strings.Split(srcset, "\n")
	assert.Equal(t, len(splitSrcset), 31)
}

func TestReadMe_SignedSrcSetCreation(t *testing.T) {
	// Instead of using dotenv, just set the environment variable directly.
	const key = "IX_TOKEN"
	const value = "MYT0KEN"
	os.Setenv(key, value)

	ixToken := os.Getenv(key)
	assert.Equal(t, value, ixToken)

	ub := NewURLBuilder("demos.imgix.net",
		WithToken(ixToken),
		WithLibParam(false))
	srcset := ub.CreateSrcset("image.png", []IxParam{})

	expectedLength := 31
	splitSrcSet := strings.Split(srcset, ",\n")

	for _, u := range splitSrcSet {
		assert.Contains(t, u, "s=")
	}

	actualLength := len(splitSrcSet)
	assert.Equal(t, expectedLength, actualLength)
}

func TestReadMe_FixedWidthSrcSetDefault(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))
	params := []IxParam{Param("h", "800"), Param("ar", "4:3")}
	expected := "https://demo.imgix.net/image.png?ar=4%3A3&dpr=1&h=800&q=75 1x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=2&h=800&q=50 2x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=3&h=800&q=35 3x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=4&h=800&q=23 4x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=5&h=800&q=20 5x"
	actual := ub.CreateSrcset("image.png", params)
	assert.Equal(t, expected, actual)
}

func TestReadMe_FixedWidthSrcSetVariableQualityDisabled(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))
	params := []IxParam{Param("h", "800"), Param("ar", "4:3")}
	expected := "https://demo.imgix.net/image.png?ar=4%3A3&dpr=1&h=800 1x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=2&h=800 2x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=3&h=800 3x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=4&h=800 4x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=5&h=800 5x"
	actual := ub.CreateSrcset("image.png", params, WithVariableQuality(false))
	assert.Equal(t, expected, actual)
}

func TestReadMe_FixedWidthSrcSetNoOpts(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))
	params := []IxParam{Param("h", "800"), Param("ar", "4:3")}
	expected := "https://demo.imgix.net/image.png?ar=4%3A3&dpr=1&h=800&q=75 1x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=2&h=800&q=50 2x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=3&h=800&q=35 3x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=4&h=800&q=23 4x,\n" +
		"https://demo.imgix.net/image.png?ar=4%3A3&dpr=5&h=800&q=20 5x"
	actual := ub.CreateSrcset("image.png", params)
	assert.Equal(t, expected, actual)
}

func TestReadMe_FluidWidthSrcSetFromWidths(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))
	ixParams := []IxParam{Param("mask", "ellipse")}
	actual := ub.CreateSrcsetFromWidths("image.jpg", ixParams, []int{100, 200, 300, 400})
	expected := "https://demo.imgix.net/image.jpg?mask=ellipse&w=100 100w,\n" +
		"https://demo.imgix.net/image.jpg?mask=ellipse&w=200 200w,\n" +
		"https://demo.imgix.net/image.jpg?mask=ellipse&w=300 300w,\n" +
		"https://demo.imgix.net/image.jpg?mask=ellipse&w=400 400w"
	assert.Equal(t, expected, actual)
}

func TestReadMe_FluidWidthSrcSet(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))

	actual := ub.CreateSrcset(
		"image.png",
		[]IxParam{},
		WithMinWidth(100),
		WithMaxWidth(380),
		WithTolerance(0.08))

	expected := "https://demo.imgix.net/image.png?w=100 100w,\n" +
		"https://demo.imgix.net/image.png?w=116 116w,\n" +
		"https://demo.imgix.net/image.png?w=135 135w,\n" +
		"https://demo.imgix.net/image.png?w=156 156w,\n" +
		"https://demo.imgix.net/image.png?w=181 181w,\n" +
		"https://demo.imgix.net/image.png?w=210 210w,\n" +
		"https://demo.imgix.net/image.png?w=244 244w,\n" +
		"https://demo.imgix.net/image.png?w=283 283w,\n" +
		"https://demo.imgix.net/image.png?w=328 328w,\n" +
		"https://demo.imgix.net/image.png?w=380 380w"
	assert.Equal(t, expected, actual)
}

func TestReadMe_FluidWidthSrcsetTolerance20(t *testing.T) {
	ub := NewURLBuilder("demo.imgix.net", WithLibParam(false))

	srcsetOptions := []SrcsetOption{
		WithMinWidth(100),
		WithMaxWidth(384),
		WithTolerance(0.20),
	}

	actual := ub.CreateSrcset(
		"image.png",
		[]IxParam{},
		srcsetOptions...)

	expected := "https://demo.imgix.net/image.png?w=100 100w,\n" +
		"https://demo.imgix.net/image.png?w=140 140w,\n" +
		"https://demo.imgix.net/image.png?w=196 196w,\n" +
		"https://demo.imgix.net/image.png?w=274 274w,\n" +
		"https://demo.imgix.net/image.png?w=384 384w"
	assert.Equal(t, expected, actual)
}

func TestReadMe_TargetWidths(t *testing.T) {
	expected := []int{300, 378, 476, 600, 756, 953, 1200, 1513, 1906, 2401, 3000}
	actual := TargetWidths(300, 3000, 0.13)
	assert.Equal(t, expected, actual)

	sm := expected[:3]
	expectedSm := []int{300, 378, 476}
	assert.Equal(t, expectedSm, sm)

	md := expected[3:7]
	expectedMd := []int{600, 756, 953, 1200}
	assert.Equal(t, expectedMd, md)

	lg := expected[7:]
	expectedLg := []int{1513, 1906, 2401, 3000}
	assert.Equal(t, expectedLg, lg)

	ub := NewURLBuilder("demos.imgix.net")
	ub.SetUseLibParam(false)
	srcset := ub.CreateSrcsetFromWidths("image.png", []IxParam{}, sm)
	actualSrcset := "https://demos.imgix.net/image.png?w=300 300w,\n" +
		"https://demos.imgix.net/image.png?w=378 378w,\n" +
		"https://demos.imgix.net/image.png?w=476 476w"
	assert.Equal(t, actualSrcset, srcset)
}
