package config

import (
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major   int
	Minor   int
	Patch   int
	Special string
}

func (v *Version) String() string {
	versionStr := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Special != "" {
		versionStr = fmt.Sprintf("%s-%s", versionStr, v.Special)
	}
	return versionStr
}

func ParseVersion(versionStr string) (Version, error) {
	version := Version{
		Major:   0,
		Minor:   0,
		Patch:   0,
		Special: "",
	}
	numbers := strings.Split(versionStr, ".")
	if len(numbers) > 0 {
		number, err := strconv.ParseInt(numbers[0], 10, 32)
		if err != nil {
			return version, err
		}
		version.Major = int(number)
	}
	if len(numbers) > 1 {
		number, err := strconv.ParseInt(numbers[1], 10, 32)
		if err != nil {
			return version, err
		}
		version.Minor = int(number)
	}
	if len(numbers) > 2 {
		numberStr := numbers[2]
		if strings.Contains(numberStr, "-") {
			parts := strings.Split(numberStr, "-")
			if len(parts) > 1 {
				numberStr = parts[0]
				version.Special = parts[1]
			}
		}
		number, err := strconv.ParseInt(numberStr, 10, 32)
		if err != nil {
			return version, err
		}
		version.Patch = int(number)
	}
	return version, nil
}

func (v *Version) IsGreaterThan(o Version) bool {
	if v.Major > o.Major {
		return true
	}
	if v.Major < o.Major {
		return false
	}

	// Major Version is equal
	if v.Minor > o.Minor {
		return true
	}
	if v.Minor < o.Minor {
		return false
	}

	// Minor Version is equal
	if v.Patch > o.Patch {
		return true
	}
	if v.Patch < o.Patch {
		return false
	}

	return false
}

func (v *Version) IsSmallerThan(o Version) bool {
	if v.Major < o.Major {
		return true
	}
	if v.Major > o.Major {
		return false
	}

	// Major Version is equal
	if v.Minor < o.Minor {
		return true
	}
	if v.Minor > o.Minor {
		return false
	}

	// Minor Version is equal
	if v.Patch < o.Patch {
		return true
	}
	if v.Patch > o.Patch {
		return false
	}

	return false
}

func (v *Version) IsEqual(o Version) bool {
	if v.Major != o.Major {
		return false
	}

	// Major Version is equal
	if v.Minor != o.Minor {
		return false
	}

	// Minor Version is equal
	if v.Patch != o.Patch {
		return false
	}

	if v.Special != o.Special {
		return false
	}
	return true
}
