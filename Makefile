build:
	docker build -t forum . 

run: 
	docker run -p 8989:8989 forum

