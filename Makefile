.DEFAULT_GOAL := all

all:
	@echo "\nUsage:  make <command>"
	@echo "\nCommands:"
	@echo "   all          list all commands"
	@echo ""
	@echo "   dbgui        open DB GUI (password: xemmenu) (required dev_prep - MacOS)"
	@echo "   mqgui        open RabbitMQ Management GUI (user: xemmenu, password: xemmenu) (required dev_prep - MacOS)"
	@echo ""
	@echo "   dev_prep     prepare and run infra services in local environment"
	@echo "   dev_stop     stop infra services in local environment"
	@echo "   dev_clean    stop and remove infra services in local environment"
	@echo ""
	@echo "   run_order    run order service"
	@echo "   run_kitchen  run kitchen service"
	@echo "   run_payment  run payment service"
	@echo "\n"

dbgui:
	@open http://localhost:8000

mqgui:
	@open http://localhost:15672

dev_prep:
	@docker compose -f ./deployments/dev/docker-compose.dev.yaml up -d

dev_stop:
	@docker compose -f ./deployments/dev/docker-compose.dev.yaml down

dev_clean:
	@docker compose -f ./deployments/dev/docker-compose.dev.yaml down
	@docker volume prune -f

run_order:
	@cd services/order && make run

run_kitchen:
	@cd services/kitchen && make run

run_payment:
	@cd services/payment && make run

.PHONY: all dbgui mqgui dev_prep dev_stop dev_clean run_order run_kitchen run_payment
