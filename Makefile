build: 
	docker build -t ins_ecosystem_testtask .

run: build
	docker run --rm -it ins_ecosystem_testtask	