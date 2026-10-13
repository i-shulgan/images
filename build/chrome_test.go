package build

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestMatchChromeDriver(t *testing.T) {
	versions := map[string]string{
		"120.0.6099.1":   "u1",
		"120.0.6099.109": "u109",
		"120.0.6099.71":  "u71",
		"120.0.6100.5":   "other-build",
		"121.0.6167.85":  "next",
	}
	for chrome, want := range map[string]string{
		"120.0.6099.71":  "120.0.6099.71",  // exact
		"120.0.6099.1":   "120.0.6099.1",   // exact, not 120.0.6099.109
		"120.0.6099.200": "120.0.6099.109", // newest of the same build
		"120.0.6200.1":   "120.0.6100.5",   // newest of the same major
		"121.0.6167.100": "121.0.6167.85",
	} {
		v, ok := matchChromeDriver(chrome, versions)
		assert.True(t, ok, chrome)
		assert.Equal(t, want, v, chrome)
	}
	_, ok := matchChromeDriver("12.0.1.1", versions)
	assert.False(t, ok, "12 must not match 120 or 121")
	_, ok = matchChromeDriver("119.0.1.1", versions)
	assert.False(t, ok)
}

func TestParseChromeDriverVersionExplicit(t *testing.T) {
	c := &Chrome{Requirements{DriverVersion: "1.2.3"}}
	v, err := c.parseChromeDriverVersion("120.0.6099.71", nil)
	assert.NoError(t, err)
	assert.Equal(t, "1.2.3", v)
}
