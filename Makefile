up:
	docker compose up

remove:
	docker compose down -v
	sudo rm -rf ./tarantool/data/*
	sudo rm -rf /tmp/tarantool/run/*

test:
	go test -v -count=1 -run $(TEST) ./$(PKG)