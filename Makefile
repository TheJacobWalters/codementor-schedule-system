
test:
	go test 
 
run:
	go run main.go

redis:
	docker run -p 6379:6379 redis

test-single-cmd:
	wget "http://localhost:5050/addTask?Command=ls&Argument=/tmp"

test-schedule-cmd:
	wget "http://localhost:5050/addTask?Command=ls&Argument=/tmp&Time=*/1 * * * *"
	