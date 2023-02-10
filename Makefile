#!/bin/sh

all:
	@docker-compose up -d --build
	@echo "Calculating..."
	@docker commit parse-log debug/ubuntu > /dev/null \
		&& docker run -it --rm --entrypoint sh debug/ubuntu -c "cat average-time-request.csv" > average-time-request.csv \
		&& docker run -it --rm --entrypoint sh debug/ubuntu -c "cat service-request.csv" > service-request.csv \
		&& docker run -it --rm --entrypoint sh debug/ubuntu -c "cat consumer-request.csv" > consumer-request.csv ;
	@echo "done!"

build:
	docker-compose build

up:
	docker-compose up

clean:
	@rm *.csv api/*.txt
	@echo "done!"

stop:
	@if (docker ps -f name=parse-log | grep -q "parse-log"); then \
		docker stop parse-log; \
		docker stop parse-log-mongodb
	else \
		echo "$$(tput bold)There is no running container.$$(tput sgr0)"; \
	fi \

.PHONY: help

help:
	@echo "$$(tput bold)Available commands:$$(tput sgr0)"; echo "\n\t$$(tput bold)all:$$(tput sgr0) \tbuild the image of the DB and up the container.\n\t$$(tput bold)build:$$(tput sgr0) \tbuild the image of the DB.\n\t$$(tput bold)up:$$(tput sgr0) \tup the container.\n\t$$(tput bold)stop:$$(tput sgr0) \tstop the current container.\n\t$$(tput bold)clean:$$(tput sgr0) \tclean the project."
	@echo "$$(tput bold)\nUsage example:$$(tput sgr0) make all"

.DEFAULT_GOAL := help
