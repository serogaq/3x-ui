#!/bin/sh

git remote add upstream https://github.com/MHSanaei/3x-ui.git
git fetch upstream
git checkout main
git reset --hard upstream/main
git push --force origin main
git remote remove upstream
git fetch --prune --all
git tag -d $(git tag -l)
git fetch --tags origin
