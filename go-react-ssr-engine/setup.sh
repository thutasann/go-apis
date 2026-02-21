#!/bin/bash

# reactgo project scaffold - run from where you want the project root

mkdir -p reactgo/{cmd/server,internal/{engine,router,bundler,cache,hydration,props,config},pkg/html,pages/posts,public,components,.build/{server,client}}

cd reactgo

# Go files
touch cmd/server/main.go
touch internal/engine/{engine.go,pool.go,renderer.go}
touch internal/router/{router.go,tree.go,middleware.go}
touch internal/bundler/{bundler.go,watcher.go,transform.go}
touch internal/cache/{cache.go,invalidator.go}
touch internal/hydration/{hydration.go,manifest.go}
touch internal/props/{loader.go,context.go}
touch internal/config/config.go
touch pkg/html/{document.go,head.go}

# Config
touch reactgo.config.json

# Placeholder React pages
touch pages/index.tsx
touch pages/about.tsx
touch "pages/posts/[id].tsx"

# Module init
touch go.mod go.sum

# Gitignore build artifacts
echo ".build/" > .gitignore

echo "done - structure ready at ./reactgo"