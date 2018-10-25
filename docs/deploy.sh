#!/usr/bin/env bash

cd src
bundle exec middleman build --clean
cd ..
cp -r src/build/* .
git add ./*
git commit -m "docs: build"
