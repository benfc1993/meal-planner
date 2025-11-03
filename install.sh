#!/usr/bin/env bash
sudo cp -f ./meal-planner.service /etc/systemd/system/meal-planner.service
make build
mkdir ./_temp
sudo mv /usr/local/bin/meal-planner/my.db ./_temp/
sudo rm -rf /usr/local/bin/meal-planner/
sudo mv ./dist/ /usr/local/bin/meal-planner/
sudo mv -f ./_temp/my.db /usr/local/bin/meal-planner/my.db
rm -rf ./_temp
sudo systemctl enable meal-planner
sudo systemctl restart meal-planner
