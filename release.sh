#!/bin/bash
#
# Tag a release of `builder` or `server`, build a Linux/amd64 Docker image,
# and push it to the Forgejo container registry at code.andydayton.com as
# both :<tag> and :latest. The git tag is namespaced per app (e.g.
# server-v0.1.0) so the two apps can release independently from the same
# repo. Requires `docker login code.andydayton.com` to have been run.
#
# Pass --dry-run to print every mutating command without executing it.

# -e: abort on any failed command
# -u: abort on use of an unset variable (catches typos in var names)
# -o pipefail: a pipeline fails if any stage fails, not just the last one
set -euo pipefail

USAGE="Usage: $0 <builder|server> <vX.Y.Z> [--dry-run]"

# Pull --dry-run out of $@ so the positional-argument check below stays simple
# regardless of where the flag was placed on the command line.
DRY_RUN=0
ARGS=()
for arg in "$@"; do
  if [[ "$arg" == "--dry-run" ]]; then
    DRY_RUN=1
  else
    ARGS+=("$arg")
  fi
done
# `${ARGS[@]+"${ARGS[@]}"}` expands to nothing when ARGS is empty. The bare
# `"${ARGS[@]}"` form trips `set -u` on bash 3.2 (the default on macOS).
set -- "${ARGS[@]+"${ARGS[@]}"}"

if [ "$#" -ne 2 ]; then
  echo "$USAGE"
  exit 1
fi

if [[ "$1" != "builder" && "$1" != "server" ]]; then
  echo "$USAGE"
  exit 1
fi

APP=$1
TAG=$2
GITTAG="$APP-$TAG"

# Reject typos like `v.0.1.0` or `0.1.0` before we do anything destructive.
if [[ ! "$TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: tag must match vX.Y.Z (got: $TAG)"
  exit 1
fi

# Wrapper used by every mutating step below. In dry-run mode it prints the
# command (with %q quoting so the output is copy-pasteable) instead of running
# it. Pre-flight checks deliberately do NOT go through `run` — we always want
# them to execute for real, even in dry-run.
run() {
  if [[ "$DRY_RUN" == "1" ]]; then
    printf '[dry-run]'
    printf ' %q' "$@"
    printf '\n'
  else
    "$@"
  fi
}

# ---- Pre-flight checks --------------------------------------------------
# All of these run before the confirmation prompt so we fail fast and never
# push a git tag that we can't follow up with a successful image push.

# Docker Desktop / dockerd must be reachable; otherwise the build/push at the
# end would fail after we've already pushed the git tag.
if ! docker info >/dev/null 2>&1; then
  echo "Error: Docker daemon is not running"
  exit 1
fi

# Verify we can reach the Forgejo registry as an authenticated user, so an
# expired/missing token fails here instead of after the git tag is already
# pushed. We probe a general-use sentinel image that lives on the registry
# purely for this purpose. `docker manifest inspect` is used (over `docker
# pull`) because it always hits the registry — no local-cache fast path —
# and doesn't write any bytes to disk.
#
# If the sentinel image ever needs to be recreated:
#
#   cat > /tmp/Dockerfile.sentinel <<'EOF'
#   FROM scratch
#   LABEL purpose="general-use auth preflight sentinel for code.andydayton.com"
#   EOF
#   docker buildx build --platform linux/amd64 \
#     -t code.andydayton.com/andy/sentinel-image:latest \
#     -f /tmp/Dockerfile.sentinel /tmp
#   docker push code.andydayton.com/andy/sentinel-image:latest
if ! docker manifest inspect code.andydayton.com/andy/sentinel-image:latest >/dev/null 2>&1; then
  echo "Error: cannot reach Forgejo registry at code.andydayton.com as authenticated user"
  echo "Run: docker login code.andydayton.com"
  exit 1
fi

# Refuse to release with uncommitted changes — the tag would point at HEAD
# but the working tree wouldn't match, which is confusing at best.
if [[ -n "$(git status --porcelain)" ]]; then
  echo "Error: working tree has uncommitted changes"
  git status --short
  exit 1
fi

# Catch the "tag already exists" case up front; otherwise `git tag` would
# fail later with a less obvious error.
if git rev-parse -q --verify "refs/tags/$GITTAG" >/dev/null; then
  echo "Error: tag $GITTAG already exists"
  exit 1
fi

# Remember where the user started so we can return them there on exit. Falls
# back to the bare commit SHA if HEAD is detached.
ORIGINAL_REF=$(git symbolic-ref --short HEAD 2>/dev/null || git rev-parse HEAD)

# Trap on EXIT so we restore the original branch on success, failure, or a
# user-aborted prompt — not just at the end of the happy path. We only run
# `git checkout` if HEAD has actually moved, to avoid noise on early exits.
cleanup() {
  if [[ "$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo)" != "$ORIGINAL_REF" ]]; then
    echo ""
    echo "Restoring original ref: $ORIGINAL_REF"
    git checkout "$ORIGINAL_REF" 2>/dev/null || true
  fi
}
trap cleanup EXIT

# ---- Confirmation summary ----------------------------------------------

echo ""
echo "Repo status:"
echo ""
git status --short --branch

echo ""
echo "Existing $APP tags:"
echo ""
# `|| true` because `git tag --list` returns 0 even with no matches, but we
# keep the guard in case future flags or pipes change that behavior — and to
# stay safe under `set -e`.
git tag --list "$APP-*" || true

echo ""
echo "Commit to be tagged:"
git log -1 --oneline

echo ""
echo "App:           $APP"
echo "Tag:           $TAG"
echo "Git tag:       $GITTAG"
echo "Original ref:  $ORIGINAL_REF"
if [[ "$DRY_RUN" == "1" ]]; then
  echo "Mode:          DRY RUN"
fi

echo ""
echo "Please review output above. Continue with tag/release workflow? (y/N)"
read -r CONFIRM

if [[ ! "$CONFIRM" =~ ^[Yy]$ ]]; then
  exit 0
fi

# ---- Release steps ------------------------------------------------------

# Tag locally, push only this tag (not `--tags`, which would push every
# unrelated local tag), then check the tag out so the build is reproducible
# from the tag ref rather than whatever HEAD happened to be.
run git tag "$GITTAG"
run git push origin "$GITTAG"
run git checkout "$GITTAG"

# Build a fresh Linux binary; the Dockerfile expects it pre-built in bin/.
run make clean-bin "bin/$APP-linux-amd64"

# Build the container image. Pinned to linux/amd64 because that's what the
# production host runs — buildx handles cross-compile from arm64 macs.
run docker buildx build --platform linux/amd64 \
  -f "app/$APP/Dockerfile" \
  -t "code.andydayton.com/andy/yesterdaysnews-$APP:$TAG" .

# Push the tagged image first, then re-tag and push as :latest. Doing it in
# this order means :latest never points at an image that isn't also pushed
# under its versioned tag.
run docker push "code.andydayton.com/andy/yesterdaysnews-$APP:$TAG"

run docker tag "code.andydayton.com/andy/yesterdaysnews-$APP:$TAG" "code.andydayton.com/andy/yesterdaysnews-$APP:latest"
run docker push "code.andydayton.com/andy/yesterdaysnews-$APP:latest"
