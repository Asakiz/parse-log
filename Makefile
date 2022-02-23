#!/bin/sh

all:
	docker-compose up -d --build

build:
	docker-compose build

up:
	docker-compose up

stop:
	@if (docker ps -f name=parse_log_1 | grep -q "parse_log_1"); then \
		docker stop parse_log_1; \
	else \
		echo "$$(tput bold)There is no running container.$$(tput sgr0)"; \
	fi \

.PHONY: help

help:
	@echo "$$(tput bold)Available commands:$$(tput sgr0)"; echo "\n\t$$(tput bold)all:$$(tput sgr0) \tbuild the image of the DB and up the container.\n\t$$(tput bold)build:$$(tput sgr0) \tbuild the image of the DB.\n\t$$(tput bold)up:$$(tput sgr0) \tup the container.\n\t$$(tput bold)stop:$$(tput sgr0) \tstop the current container."

.DEFAULT_GOAL := help
