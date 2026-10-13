package build

import (
	assert "github.com/stretchr/testify/require"
	"testing"
)

func TestVersionFromPackageName(t *testing.T) {
	assert.Equal(t, "45.0.2", extractVersion("firefox_45.0.2+build1-0ubuntu0.14.04.1+aerokube0_amd64.deb"))
	assert.Equal(t, "71.0", extractVersion("firefox_71.0_b11+build1-0ubuntu0.18.04.1_amd64.deb"))
	assert.Equal(t, "72.0", extractVersion("firefox-trunk_72.0_a1_hg20191114r501767-0ubuntu0.18.04.1_umd1_amd64.deb"))
	assert.Equal(t, "48.0.2564.109", extractVersion("google-chrome-stable_48.0.2564.109-1+aerokube0_amd64.deb"))
	assert.Equal(t, "79.0.3945.36", extractVersion("google-chrome-beta_79.0.3945.36-1_amd64.deb"))
	assert.Equal(t, "80.0.3964.0", extractVersion("google-chrome-unstable_80.0.3964.0-1_amd64.deb"))
	assert.Equal(t, "38.0.2220.31", extractVersion("opera-stable_38.0.2220.31_amd64.deb"))
	assert.Equal(t, "65.0.3467.32", extractVersion("opera-beta_65.0.3467.32_amd64.deb"))
	assert.Equal(t, "66.0.3502.0", extractVersion("opera-developer_66.0.3502.0_amd64.deb"))
	assert.Equal(t, "19.9.3.358", extractVersion("yandex-browser-beta_19.9.3.358-1_amd64.deb"))
}

func TestVersionToTag(t *testing.T) {
	assert.Equal(t, "45.0.2", extractVersion("45.0.2+build1-0ubuntu0.14.04.1+aerokube0"))
	assert.Equal(t, "48.0.2564.109", extractVersion("48.0.2564.109-1+aerokube0"))
	assert.Equal(t, "78.0.3904.108", extractVersion("78.0.3904.108-1"))
	assert.Equal(t, "38.0.2220.31", extractVersion("38.0.2220.31"))
	assert.Equal(t, "71.0", extractVersion("71.0~b11+build1-0ubuntu0.18.04.1"))
	assert.Equal(t, "74.0", extractVersion("74.0~a1~hg20200205r512485-0ubuntu0.18.04.1~umd1"))
}

func TestParseVersion(t *testing.T) {
	assert.Equal(t, "78", majorVersion("78.0.3904.108"))
	assert.Equal(t, "78", majorVersion("78"))
	assert.Equal(t, "78.0", majorMinorVersion("78.0.3904.108"))
	assert.Equal(t, "78.0", majorMinorVersion("78.0"))
	assert.Equal(t, "78.0.3904", buildVersion("78.0.3904.108"))
	assert.Equal(t, "78.0.3904", buildVersion("78.0.3904"))
}
