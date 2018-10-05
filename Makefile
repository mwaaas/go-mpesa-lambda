app_version=$(shell git rev-parse HEAD 2> /dev/null | sed "s/\(.*\)/\1/")
env=development
profile=mwas
debug=false

compile:
	docker-compose run app go build -o devops/dist/main

build_ansible:
	docker-compose -f docker-compose-ansible.yml build ansible

deploy: compile build_ansible
	docker-compose -f docker-compose-ansible.yml run \
	ansible ansible-playbook devops/deploy.yml \
	--connection local -e "app_version=$(app_version)" \
	-e "env=$(env)" -e "profile=$(profile)" -e debug=$(debug);