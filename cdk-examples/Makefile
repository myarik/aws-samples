.PHONY: dev lint deploy destroy deps service-tests infra-tests coverage-tests build format
PYTHON := ".venv/bin/python3"
.ONESHELL:  # run all commands in a single shell, ensuring it runs within a local virtual env

dev:
	pip install --upgrade pip poetry
	pre-commit install
# ensures poetry creates a local virtualenv (.venv)
	poetry config --local virtualenvs.in-project true
	poetry install --no-root
	npm ci

format:
	poetry run ruff check . --fix

format-fix:
	poetry run ruff format .

lint: format
	@echo "Running mypy"
	$(MAKE) mypy-lint

mypy-lint:
	poetry run mypy --pretty service cdk tests

deps:
	poetry export --only=dev --format=requirements.txt > dev_requirements.txt
	poetry export --without=dev --format=requirements.txt > lambda_requirements.txt

build: deps
	mkdir -p .build/lambdas ; cp -r service .build/lambdas
	mkdir -p .build/common_layer ; poetry export --without=dev --format=requirements.txt > .build/common_layer/requirements.txt
	mkdir -p .build/input_validation_layer ; poetry export --with=input_validation --format=requirements.txt > .build/input_validation_layer/requirements.txt
	mkdir -p .build/multi_layers ; poetry export --only=multi_layers --format=requirements.txt > .build/multi_layers/requirements.txt

infra-tests: build
	poetry run pytest -l -s --pdb tests/infrastructure

service-tests: build
	poetry run pytest -l -s --pdb tests/service

coverage-tests:
	poetry run pytest tests/service tests/infrastructure  --cov-config=.coveragerc --cov=service --cov-report xml

deploy: build
	npx cdk deploy --app="${PYTHON} ${PWD}/app.py" --require-approval=never

destroy:
	npx cdk destroy --app="${PYTHON} ${PWD}/app.py" --force
