#!/usr/bin/env bash

set -e

BEFORE_DIFF=$(git diff | sha1sum )
BEFORE_STATUS=$(git status --porcelain | sha1sum)

make generate

AFTER_DIFF=$(git diff | sha1sum )
AFTER_STATUS=$(git status --porcelain | sha1sum)

if [[ $BEFORE_DIFF != $AFTER_DIFF || $BEFORE_STATUS != $AFTER_STATUS ]]; then
  echo "Unstable generate. Make sure to generate and check in changed files."
  exit 1
fi