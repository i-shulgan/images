# Browser Images
[![ci](https://github.com/i-shulgan/images/actions/workflows/ci.yml/badge.svg)](https://github.com/i-shulgan/images/actions/workflows/ci.yml)
[![chrome](https://github.com/i-shulgan/images/actions/workflows/chrome.yml/badge.svg)](https://github.com/i-shulgan/images/actions/workflows/chrome.yml)

Docker build files and the `images` tool that builds browser images for [Selenoid](https://github.com/i-shulgan/selenoid).

This is a maintained fork of [aerokube/images](https://github.com/aerokube/images), which is archived since 7.3.4; its last Chrome images on Docker Hub stopped at 128. The fork builds an image with the current stable Google Chrome every day.

## Chrome image

| Tag | Content |
|-----|---------|
| `ghcr.io/i-shulgan/chrome:<major>.<minor>`, e.g. `141.0` | latest build of that Chrome version |
| `ghcr.io/i-shulgan/chrome:<full version>`, e.g. `141.0.7390.122` | a specific Chrome build |
| `ghcr.io/i-shulgan/chrome:latest` | latest stable Chrome |

The image contains Chrome, the matching chromedriver (port 4444), Xvfb, VNC (`ENABLE_VNC=true`), PulseAudio, the DevTools proxy (port 7070), a file server for downloads (8080) and clipboard access (9090), just like `selenoid/vnc_chrome` images. It works with Selenoid's `enableVNC`, `enableVideo`, `screenResolution`, `/devtools/`, `/download/` and `/clipboard/`.

`browsers.json` for Selenoid:

```json
{
  "chrome": {
    "default": "141.0",
    "versions": {
      "141.0": {
        "image": "ghcr.io/i-shulgan/chrome:141.0",
        "port": "4444",
        "path": "/",
        "shmSize": 2147483648
      }
    }
  }
}
```

Pull the image before the first session (`docker pull ghcr.io/i-shulgan/chrome:141.0`): Selenoid does not pull images itself.

### How it is built

[`chrome.yml`](.github/workflows/chrome.yml) runs every day, on changes to Chrome or base image files and on demand (**Actions → chrome → Run workflow**, optionally with a `google-chrome-stable` package version and `force`):

1. reads the current `google-chrome-stable` version from Google's apt repository;
2. skips the run if `ghcr.io/i-shulgan/chrome:<full version>` is already published (scheduled and manual runs without `force`);
3. builds the base image from `selenium/base` (Ubuntu 24.04, X server, fonts, helpers) and the Chrome image on top of it;
4. runs [`ci/chrome-smoke.sh`](ci/chrome-smoke.sh): a session against chromedriver in the image (navigation, title, screenshot, VNC process), then a session through Selenoid (`enableVNC`, `/status`, `/devtools/`, delete);
5. pushes the three tags above. Pull requests only build and test.

## Building images by hand

Requirements: Docker and Go 1.27.2 or newer. Run from the repository root (or pass `--static-dir`):

```bash
go build -o images .
docker build -t browsers/base:local selenium/base
./images chrome --base-image browsers/base:local \
    -b 141.0.7390.122-1 -t my/chrome:141.0
```

`-b` takes a `google-chrome-stable` apt package version, a `.deb` file or its URL. Without `--base-image` the Dockerfiles use `browsers/base:7.4.2` from Docker Hub. `./images --help` lists other browsers and flags; only Chrome is built and tested in CI.

## Differences from upstream

* Daily Chrome image in `ghcr.io/i-shulgan/chrome`, tested with chromedriver and Selenoid.
* chromedriver is now the exact Chrome for Testing build for the Chrome version (or the newest of the same build / major); upstream took the oldest version containing the Chrome version as a substring.
* A corrupted driver archive is reported instead of panicking; a non-numeric Firefox version is an error instead of silently selecting the legacy image.
* `--base-image`, base image on Ubuntu 24.04, Chrome apt key in a keyring file instead of `apt-key`.
* Go 1.27.2 for the tool and the helpers in images (`fileserver`, `xseld`, `devtools`), updated dependencies, `pkger` replaced by reading `static/` from the checkout.
* CI: `go vet`, `staticcheck`, `go test -race`, `govulncheck` for all Go modules, shellcheck; semantic versions of the tool with release-please.

## Releases and versioning

Versions of the `images` tool follow [Semantic Versioning](https://semver.org), continue from upstream `7.3.4` and are derived from [Conventional Commits](https://www.conventionalcommits.org) by [release-please](https://github.com/googleapis/release-please): `fix:` bumps the patch version, `feat:` the minor one, `!` or `BREAKING CHANGE:` the major one. After merges to `master` release-please opens a `chore(master): release X.Y.Z` pull request; merging it creates the `vX.Y.Z` tag and GitHub release. Pull request titles must follow Conventional Commits (the `pr-title` check), since squash merges use them as commit messages.

Browser images are tagged by browser version, not by the tool version.
