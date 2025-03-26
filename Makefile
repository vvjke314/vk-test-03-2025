up:
	docker compose up

remove:
	docker compose down -v
	rm -rf ./tarantool/data/*
	rm -rf /tmp/tarantool/run/*