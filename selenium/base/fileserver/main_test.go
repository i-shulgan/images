package main

import (
	"encoding/json"
	assert "github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

var (
	dir string
	srv *httptest.Server
)

func init() {
	dir, _ = os.MkdirTemp("", "fileserver")
	srv = httptest.NewServer(mux(dir))
}

func withUrl(path string) string {
	return srv.URL + path
}

func TestDownloadAndRemoveFile(t *testing.T) {
	tempFile, _ := os.CreateTemp(dir, "fileserver")
	_ = os.WriteFile(tempFile.Name(), []byte("test-data"), 0644)
	tempFileName := filepath.Base(tempFile.Name())
	resp, err := http.Get(withUrl("/" + tempFileName))
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	_, err = os.Stat(tempFile.Name())
	assert.NoError(t, err)

	rsp, err := http.Get(withUrl("/?json"))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rsp.StatusCode)
	var files []string
	assert.NoError(t, json.NewDecoder(rsp.Body).Decode(&files))
	assert.Equal(t, []string{tempFileName}, files)

	req, _ := http.NewRequest(http.MethodDelete, withUrl("/"+tempFileName), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	_, err = os.Stat(tempFile.Name())
	assert.Error(t, err)
}

func TestRemoveMissingFile(t *testing.T) {
	req, _ := http.NewRequest(http.MethodDelete, withUrl("/missing-file"), nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
