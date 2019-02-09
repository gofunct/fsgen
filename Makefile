install:
	go install

init: install
	cd modules; go generate