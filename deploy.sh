#!/bin/bash

./mhxr
cd ./pages

git config user.email "maccotsan@gmail.com"
git config user.name "maccotsan"

git commit -a -m "Update today schedule."
git push --force --quiet "https://${GH_TOKEN}@github.com/maccotsan/mhxr.git" gh-pages:gh-pages > /dev/null 2>&1
