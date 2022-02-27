#!/bin/sh

all: 
	@docker-compose up -d --build
	cd api; go run main.go $(input)
	
build:
	docker-compose build

up:
	docker-compose up

clean:
	@rm api/*.csv api/*.txt
	@echo "done."

stop:
	@if (docker ps -f name=parse_log_1 | grep -q "parse_log_1"); then \
		docker stop parse_log_1; \
	else \
		echo "$$(tput bold)There is no running container.$$(tput sgr0)"; \
	fi \

.PHONY: help

help:
	@echo "$$(tput bold)Available commands:$$(tput sgr0)"; echo "\n\t$$(tput bold)all:$$(tput sgr0) \tbuild the image of the DB and up the container.\n\t\tpass input=<fileName> to run the application\n\t$$(tput bold)build:$$(tput sgr0) \tbuild the image of the DB.\n\t$$(tput bold)up:$$(tput sgr0) \tup the container.\n\t$$(tput bold)stop:$$(tput sgr0) \tstop the current container.\n\t$$(tput bold)clean:$$(tput sgr0) \tclean the project."
	@echo "$$(tput bold)\nUsage example:$$(tput sgr0) make all input=log.txt"

.DEFAULT_GOAL := help
