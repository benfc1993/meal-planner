#!/usr/bin/env bash
sudo cp -f ./meal-planner.service /etc/systemd/system/meal-planner.service
make build
mkdir _tmp
sudo cp -fr ./dist/meal-planner /usr/local/bin/meal-planner/meal-planner
sudo systemctl enable meal-planner
sudo systemctl restart meal-planner
