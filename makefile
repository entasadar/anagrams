build:
	go mod download
	go build -o ./bin/anagrams cmd/main.go

run: build
	docker-compose up -d anagram-redis
	./bin/anagrams --redis="127.0.0.1:7079"

docker-build:
	docker build -t img-anagrams .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

format:
	go fmt $$(pwd)/...

lint:
	golangci-lint run

test-unit:
	go test $$(pwd)/internal/...

test-integration:
	@printf "\e[1m\e[1;34m#### Preparing environment for integration testing ####\e[0m \e[0m\n"
	@printf "Build docker image, please wait... \n"
	@`docker build -t img-anagrams . 2>/dev/null` 2>/dev/null ||:
	docker network create --driver=bridge --subnet=192.168.77.0/24 test_net
	docker run --rm --detach --network test_net --ip 192.168.77.78 --name test-redis redis:4.0-alpine
	docker run -p 7070:7070 --env REDIS=192.168.77.78:6379 --env PORT=7070 --rm --detach \
	       --network test_net --ip 192.168.77.80 --name test-app img-anagrams
	@printf "\n\e[1m\e[32m#### Starting integration testing ####\e[0m \e[0m \n"
	@-go test -v $$(pwd)/test/...
	@printf "\e[1m\e[32m#### Testing completed! ####\e[0m \e[0m\n\n"
	@printf "\e[1m\e[1;34m#### Cleaning the test environment ####\e[0m \e[0m\n"
	docker stop test-redis test-app
	docker network rm test_net
