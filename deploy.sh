#!/bin/bash -e 

# Generate vendoring deps
godep save

# Commit changes
git add vendor Godeps
git ci -m 'updated: vendoring deps'

# Push changes to heroku
git push heroku master
