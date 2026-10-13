package build

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestCopySourceFiles(t *testing.T) {
	old := StaticFiles
	defer func() { StaticFiles = old }()
	StaticFiles = fstest.MapFS{
		"chrome/Dockerfile":       {Data: []byte("FROM base"), Mode: 0644},
		"chrome/entrypoint.sh":    {Data: []byte("#!/bin/sh"), Mode: 0755},
		"chrome/devtools/main.go": {Data: []byte("package main"), Mode: 0644},
		"firefox/apt/Dockerfile":  {Data: []byte("FROM other"), Mode: 0644},
	}
	dest := t.TempDir()
	dir, err := copySourceFiles("chrome", dest)
	if err != nil {
		t.Fatal(err)
	}
	if dir != filepath.Join(dest, "chrome") {
		t.Fatalf("unexpected dir %s", dir)
	}
	data, err := os.ReadFile(filepath.Join(dir, "devtools", "main.go"))
	if err != nil || string(data) != "package main" {
		t.Fatalf("nested file not copied: %q %v", data, err)
	}
	info, err := os.Stat(filepath.Join(dir, "entrypoint.sh"))
	if err != nil || info.Mode().Perm()&0100 == 0 {
		t.Fatalf("executable bit lost: %v %v", info, err)
	}
	if _, err := os.Stat(filepath.Join(dest, "firefox")); !os.IsNotExist(err) {
		t.Fatal("only the requested directory must be copied")
	}
	if _, err := copySourceFiles("missing", dest); err == nil {
		t.Fatal("missing directory must be an error")
	}
}

func TestStaticFilesOfRepository(t *testing.T) {
	old := StaticFiles
	defer func() { StaticFiles = old }()
	StaticFiles = os.DirFS("../static")
	dir, err := copySourceFiles("chrome", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"Dockerfile", "entrypoint.sh", "apt/Dockerfile", "devtools/devtools.go"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Errorf("%s: %v", f, err)
		}
	}
}
