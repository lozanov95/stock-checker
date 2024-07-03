include .env

repoURI := $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com/lambda
build:
	docker build --platform linux/amd64 -t $(IMAGE_NAME):test .
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ACCOUNT_ID).dkr.ecr.$(AWS_REGION).amazonaws.com
	docker tag $(IMAGE_NAME):test $(repoURI):latest
	docker push $(repoURI):latest