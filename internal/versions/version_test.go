package versions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSemverUpToDate(t *testing.T) {
	info := Compare("2.3.4", []string{"0.1", "1.0.2", "2.3.4"})
	expectedInfo := VersionInfo{
		Current:       "2.3.4",
		Available:     nil,
		UpToDate:      true,
		NonComparable: false,
	}
	assert.Equal(t, info, expectedInfo)
}

func TestSemverOutOfDate(t *testing.T) {
	info := Compare("2.3.4", []string{"0.1", "1.0.2", "2.3.4", "2.3.5", "3.0"})
	expectedInfo := VersionInfo{
		Current:       "2.3.4",
		Available:     []string{"2.3.5", "3.0"},
		UpToDate:      false,
		NonComparable: false,
	}
	assert.Equal(t, info.Current, expectedInfo.Current)
	assert.ElementsMatch(t, info.Available, expectedInfo.Available)
	assert.Equal(t, info.UpToDate, expectedInfo.UpToDate)
	assert.Equal(t, info.NonComparable, expectedInfo.NonComparable)
}

func TestNonComparable(t *testing.T) {
	info := Compare("bbb", []string{"aaa", "bbb"})
	expectedInfo := VersionInfo{
		Current:       "bbb",
		Available:     nil,
		UpToDate:      true,
		NonComparable: true,
	}
	assert.Equal(t, info, expectedInfo)
}

func TestRCPrereleaseOutOfDate(t *testing.T) {
	info := Compare("20200420RC01", []string{"20200419RC02", "20200420RC01", "20200420RC02"})
	expectedInfo := VersionInfo{
		Current:       "20200420RC01",
		Available:     []string{"20200420RC02"},
		UpToDate:      false,
		NonComparable: false,
	}
	assert.Equal(t, info, expectedInfo)
}

func TestRCHashUpToDate(t *testing.T) {
	info := Compare("20200420RC91e61f094", []string{"20200419RCfc9ea1ba4", "20200420RC2748f6a34", "20200420RCeeba0ab3c"})
	expectedInfo := VersionInfo{
		Current:       "20200420RC91e61f094",
		Available:     nil,
		UpToDate:      true,
		NonComparable: false,
	}
	assert.Equal(t, info, expectedInfo)
}

func TestRCHashOutOfDate(t *testing.T) {
	info := Compare("20200420RC91e61f094", []string{"20200419RCfc9ea1ba4", "20200420RC2748f6a34", "20200420RCeeba0ab3c", "20200421RCaf230b946"})
	expectedInfo := VersionInfo{
		Current:       "20200420RC91e61f094",
		Available:     []string{"20200421RCaf230b946"},
		UpToDate:      false,
		NonComparable: false,
	}
	assert.Equal(t, info, expectedInfo)
}

func TestRCDateOutOfDate(t *testing.T) {
	info := Compare("1.0RC20200420", []string{"0.9RC20200310", "1.0RC20200419", "1.0RC20200420", "1.0RC20200421", "1.1RC20200422"})
	expectedInfo := VersionInfo{
		Current:       "1.0RC20200420",
		Available:     []string{"1.0RC20200421", "1.1RC20200422"},
		UpToDate:      false,
		NonComparable: false,
	}
	assert.Equal(t, info.Current, expectedInfo.Current)
	assert.ElementsMatch(t, info.Available, expectedInfo.Available)
	assert.Equal(t, info.UpToDate, expectedInfo.UpToDate)
	assert.Equal(t, info.NonComparable, expectedInfo.NonComparable)
}
