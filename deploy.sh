#!bash

git worktree add pages gh-pages
go build
./mhxr
cd ./pages
git commit -a -m "Update today schedule."
git push --force --quiet "https://${GH_TOKEN}@github.com/maccotsan/mhxr.git" gh-pages:gh-pages > /dev/null 2>&1
