build: stop clean
	go build -o bin/swaggerbond
	
clean:
	rm -rf bin

start: build
	bin/swaggerbond -p 8080 -d bin/swagger-files -i 5 &

demo: build
	bin/swaggerbond -p 8080 -d bin/swagger-files -i 5 -s &

stop:
	-pkill swaggerbond