#!/usr/bin/env bash

cd src
bundle exec middleman build --clean
cd ..
cp -f src/build/* .
git add ./*
git commit -m "docs: build"
