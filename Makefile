default:
	@ go build .
	@ ./ReadAdvisor

install:
	@ go install

clean:
	@ rm ReadAdvisor

kill:
	@ pkill ReadAdvisor