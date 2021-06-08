.DEFAULT_GOAL := all

all:
	@echo "\nUsage:  make <command>"
	@echo "\nCommands:"
	@echo "   all          list all commands"
	@echo ""
	@echo "   dbgui        open DB GUI (required dev prep - MacOS)"
	@echo ""
	@echo "   dev_prep     prepare and run infra services in local environment"
	@echo "   dev_stop     stop infra services in local environment"
	@echo "   dev_clean    stop and remove infra services in local environment"
	@echo "\n"

dbgui:
	@open http://localhost:8000

dev_prep:
	@docker compose -f ./deployments/dev/docker-compose.dev.yaml up -d

dev_stop:
	@docker compose -f ./deployments/dev/docker-compose.dev.yaml down

dev_clean:
	@docker compose -f ./deployments/dev/docker-compose.dev.yaml down
	@docker volume prune -f

.PHONY: all dbgui dev_prep dev_stop dev_clean
