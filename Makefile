init:
	sudo apt-get install redis

default:
	@ go build .
	@ ./ReadAdvisor

install:
	@ go install

clean:
	@ rm ReadAdvisor

kill:
	@ pkill ReadAdvisor