#!/bin/bash
docker-compose up loadbal worker prometheus grafana -d
sleep 1m
 
docker-compose up -d
sleep 10m

docker-compose scale worker=3
sleep 10m

docker-compose scale worker=5 loadbal=2
docker-compose restart app
sleep 10m

docker-compose scale worker=2 loadbal=1
sleep 10m

docker-compose stop app
