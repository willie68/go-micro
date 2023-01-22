package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {
	ast := assert.New(t)
	versionStr := "1.2.3-special"

	version, err := ParseVersion(versionStr)
	ast.Nil(err, "Error in parsing")
	ast.Equal(1, version.Major, "Major version not equal")
	ast.Equal(2, version.Minor, "Minor version not equal")
	ast.Equal(3, version.Patch, "Patch version not equal")
	ast.Equal("special", version.Special, "Special version not equal")
}

func TestCompareEqual(t *testing.T) {
	ast := assert.New(t)
	versionStr := "1.2.3-special"

	version, err := ParseVersion(versionStr)
	ast.Nil(err, "Error in parsing")

	ast.False(version.IsEqual(Version{
		Major: 2,
	}))

	version2 := Version{
		Major:   1,
		Minor:   2,
		Patch:   3,
		Special: "special",
	}

	ast.False(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.True(version.IsEqual(version2))
	ast.False(version.IsEqual(Version{
		Major:   1,
		Minor:   2,
		Patch:   3,
		Special: "Willie",
	}))

	versionStr = "1.2.3"

	version, err = ParseVersion(versionStr)
	ast.Nil(err, "Error in parsing")

	version2 = Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}

	ast.False(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.True(version.IsEqual(version2))

	versionStr = "1.2"

	version, err = ParseVersion(versionStr)
	ast.Nil(err, "Error in parsing")

	version2 = Version{
		Major: 1,
		Minor: 2,
	}

	ast.False(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.True(version.IsEqual(version2))

	versionStr = "1"

	version, err = ParseVersion(versionStr)
	ast.Nil(err, "Error in parsing")

	version2 = Version{
		Major: 1,
	}

	ast.False(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.True(version.IsEqual(version2))
}

func TestCompareGreater(t *testing.T) {
	ast := assert.New(t)
	versionStr := "1.2.3-special"

	version, err := ParseVersion(versionStr)
	ast.Nil(err, "Error in parsing")

	ast.True(version.IsGreaterThan(Version{
		Major: 0,
	}))
	ast.False(version.IsGreaterThan(Version{
		Major: 2,
	}))
	ast.False(version.IsSmallerThan(Version{
		Major: 0,
	}))
	ast.True(version.IsSmallerThan(Version{
		Major: 2,
	}))

	version2 := Version{
		Major: 1,
		Minor: 2,
		Patch: 2,
	}

	ast.True(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.False(version2.IsGreaterThan(version))
	ast.True(version2.IsSmallerThan(version))
	ast.False(version.IsEqual(version2))

	version2 = Version{
		Major: 1,
		Minor: 2,
	}

	ast.True(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.False(version2.IsGreaterThan(version))
	ast.True(version2.IsSmallerThan(version))
	ast.False(version.IsEqual(version2))

	version2 = Version{
		Major: 1,
	}

	ast.True(version.IsGreaterThan(version2))
	ast.False(version.IsSmallerThan(version2))
	ast.False(version2.IsGreaterThan(version))
	ast.True(version2.IsSmallerThan(version))
	ast.False(version.IsEqual(version2))
}

func TestVersionString(t *testing.T) {
	ast := assert.New(t)

	version := Version{
		Major: 1,
	}
	ast.Equal("1.0.0", version.String())

	version.Minor = 9
	ast.Equal("1.9.0", version.String())

	version.Patch = 35
	ast.Equal("1.9.35", version.String())

	version.Special = "Willie"
	ast.Equal("1.9.35-Willie", version.String())
}
