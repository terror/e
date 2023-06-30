set dotenv-load

export EDITOR := 'nvim'

alias f := fmt
alias r := run

default:
  just --list

all: test lint forbid fmt-check

run *args:
	go run ./src

test:
	go test -v ./src

fmt:
	golines -m 80 -w ./src
	just retab

fmt-check:
	gofmt -l .

forbid:
	./bin/forbid

install *bin:
  go build -o {{bin}} ./src
  mv {{bin}} ~/.bin

lint:
  golangci-lint run ./src

retab:
	./bin/retab

dev-deps:
	brew install golangci-lint golines
