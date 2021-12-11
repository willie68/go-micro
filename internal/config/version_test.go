package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsing(t *testing.T) {
	versionStr := "1.2.3-special"

	version, err := ParseVersion(versionStr)
	assert.Nil(t, err, "Error in parsing")
	assert.Equal(t, 1, version.Major, "Major version not equal")
	assert.Equal(t, 2, version.Minor, "Minor version not equal")
	assert.Equal(t, 3, version.Patch, "Patch version not equal")
	assert.Equal(t, "special", version.Special, "Special version not equal")
}

func TestCompareEqual(t *testing.T) {
	versionStr := "1.2.3-special"

	version, err := ParseVersion(versionStr)
	assert.Nil(t, err, "Error in parsing")

	version2 := Version{
		Major:   1,
		Minor:   2,
		Patch:   3,
		Special: "special",
	}

	assert.False(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.True(t, version.IsEqual(version2))

	versionStr = "1.2.3"

	version, err = ParseVersion(versionStr)
	assert.Nil(t, err, "Error in parsing")

	version2 = Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}

	assert.False(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.True(t, version.IsEqual(version2))

	versionStr = "1.2"

	version, err = ParseVersion(versionStr)
	assert.Nil(t, err, "Error in parsing")

	version2 = Version{
		Major: 1,
		Minor: 2,
	}

	assert.False(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.True(t, version.IsEqual(version2))

	versionStr = "1"

	version, err = ParseVersion(versionStr)
	assert.Nil(t, err, "Error in parsing")

	version2 = Version{
		Major: 1,
	}

	assert.False(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.True(t, version.IsEqual(version2))
}

func TestCompareGreater(t *testing.T) {
	versionStr := "1.2.3-special"

	version, err := ParseVersion(versionStr)
	assert.Nil(t, err, "Error in parsing")

	version2 := Version{
		Major: 1,
		Minor: 2,
		Patch: 2,
	}

	assert.True(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.False(t, version2.IsGreaterThan(version))
	assert.True(t, version2.IsSmallerThan(version))
	assert.False(t, version.IsEqual(version2))

	version2 = Version{
		Major: 1,
		Minor: 2,
	}

	assert.True(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.False(t, version2.IsGreaterThan(version))
	assert.True(t, version2.IsSmallerThan(version))
	assert.False(t, version.IsEqual(version2))

	version2 = Version{
		Major: 1,
	}

	assert.True(t, version.IsGreaterThan(version2))
	assert.False(t, version.IsSmallerThan(version2))
	assert.False(t, version2.IsGreaterThan(version))
	assert.True(t, version2.IsSmallerThan(version))
	assert.False(t, version.IsEqual(version2))
}
