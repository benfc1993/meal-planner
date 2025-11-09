#!/usr/bin/env bash
sudo cp -f ./meal-planner.service /etc/systemd/system/meal-planner.service
make build
sudo cp -f ./dist/ /usr/local/bin/meal-planner/
sudo systemctl enable meal-planner
sudo systemctl restart meal-planner
