#!/usr/bin/env bash
# Smoke test of a Chrome image: talks to chromedriver in the image directly,
# then creates a session through Selenoid the way tests do.
#
#   ci/chrome-smoke.sh IMAGE BROWSER_VERSION [SELENOID_IMAGE]
#
# BROWSER_VERSION is the version key used in browsers.json, e.g. 141.0.
set -euo pipefail

IMAGE=$1
VERSION=$2
SELENOID_IMAGE=${3:-ghcr.io/i-shulgan/selenoid:latest}
WORK=$(mktemp -d)
DIRECT_PORT=4446
DEVTOOLS_PORT=7071
SELENOID_PORT=4445

cleanup() {
  docker rm -f chrome-smoke selenoid-smoke >/dev/null 2>&1 || true
  # Browser containers started by Selenoid.
  docker ps -aq --filter "ancestor=$IMAGE" | xargs -r docker rm -f >/dev/null 2>&1 || true
  rm -rf "$WORK"
}
trap cleanup EXIT

fail() {
  echo "::error::$*" >&2
  docker logs chrome-smoke 2>&1 | tail -50 || true
  docker logs selenoid-smoke 2>&1 | tail -50 || true
  exit 1
}

wait_for() {
  local url=$1
  for _ in $(seq 1 60); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi
    sleep 1
  done
  fail "$url did not answer in 60s"
}

# run_session BASE_URL CAPS: creates a session, opens a page and checks it.
run_session() {
  local base=$1 caps=$2 reply sid version title screenshot
  reply=$(curl -fsS -X POST -H 'Content-Type: application/json' -d "$caps" "$base/session") || fail "new session failed: $reply"
  sid=$(jq -r '.value.sessionId // .sessionId // empty' <<<"$reply")
  [ -n "$sid" ] || fail "no session ID in $reply"
  version=$(jq -r '.value.capabilities.browserVersion // empty' <<<"$reply")
  echo "session $sid, Chrome $version" >&2
  [[ "$version" == "${VERSION%%.*}."* ]] || fail "expected Chrome $VERSION, got $version"

  curl -fsS -X POST -H 'Content-Type: application/json' \
    -d '{"url":"data:text/html,<title>selenoid smoke</title><h1>ok</h1>"}' \
    "$base/session/$sid/url" >/dev/null || fail "navigation failed"
  title=$(curl -fsS "$base/session/$sid/title" | jq -r .value)
  [ "$title" = "selenoid smoke" ] || fail "unexpected title: $title"
  screenshot=$(curl -fsS "$base/session/$sid/screenshot" | jq -r .value | base64 -d | wc -c)
  [ "$screenshot" -gt 1000 ] || fail "screenshot is too small: $screenshot bytes"
  echo "$sid"
}

delete_session() {
  curl -fsS -X DELETE "$1/session/$2" >/dev/null || fail "delete session failed"
}

echo "== chromedriver in $IMAGE"
docker run -d --name chrome-smoke --shm-size=2g -p "$DIRECT_PORT:4444" -p "$DEVTOOLS_PORT:7070" -e ENABLE_VNC=true "$IMAGE" >/dev/null
wait_for "http://localhost:$DIRECT_PORT/status"
sid=$(run_session "http://localhost:$DIRECT_PORT" '{"capabilities":{"alwaysMatch":{"browserName":"chrome"}}}')
docker exec chrome-smoke pgrep -x x11vnc >/dev/null || fail "x11vnc is not running with ENABLE_VNC=true"
echo "Chrome command line and profile directories:"
docker exec chrome-smoke sh -c 'ps -eo args | grep -E "^/opt/google/chrome/chrome " | head -1; ls -a /tmp'
curl -fsS "http://localhost:$DEVTOOLS_PORT/json/protocol" | jq -e '.domains | length > 0' >/dev/null || fail "devtools helper cannot reach Chrome"
delete_session "http://localhost:$DIRECT_PORT" "$sid"
docker rm -f chrome-smoke >/dev/null

echo "== through Selenoid ($SELENOID_IMAGE)"
mkdir -p "$WORK/config" "$WORK/video"
cat >"$WORK/config/browsers.json" <<JSON
{"chrome":{"default":"$VERSION","versions":{"$VERSION":{"image":"$IMAGE","port":"4444","path":"/","shmSize":2147483648}}}}
JSON
docker run -d --name selenoid-smoke -p "$SELENOID_PORT:4444" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$WORK/config:/etc/selenoid:ro" -v "$WORK/video:/opt/selenoid/video" \
  "$SELENOID_IMAGE" -limit 2 -service-startup-timeout 2m -session-attempt-timeout 2m >/dev/null
wait_for "http://localhost:$SELENOID_PORT/status"
base="http://localhost:$SELENOID_PORT/wd/hub"
sid=$(run_session "$base" "{\"capabilities\":{\"alwaysMatch\":{\"browserName\":\"chrome\",\"browserVersion\":\"$VERSION\",\"selenoid:options\":{\"enableVNC\":true}}}}")
status=$(curl -fsS "http://localhost:$SELENOID_PORT/status")
[ "$(jq .used <<<"$status")" = 1 ] || fail "Selenoid does not count the session: $status"
[ "$(jq "[.browsers.chrome[\"$VERSION\"][].sessions[] | select(.id == \"$sid\") | .vnc] | first" <<<"$status")" = true ] || fail "VNC is not enabled: $status"
curl -fsS "http://localhost:$SELENOID_PORT/devtools/$sid/json/protocol" | jq -e '.domains | length > 0' >/dev/null || fail "devtools endpoint does not work"
delete_session "$base" "$sid"
[ "$(curl -fsS "http://localhost:$SELENOID_PORT/status" | jq .used)" = 0 ] || fail "session was not released"
echo "== smoke test passed"
