build-push-docker:
	docker build --pull . -t kuwerin/bff-website:latest
	docker push kuwerin/bff-website:latest
