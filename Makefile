APP_NAME=quay.io/cporter/assumerole
TAG=1.8

build: ## Build the container
	docker build -t $(APP_NAME):$(TAG) .
	docker push $(APP_NAME):$(TAG)

build-nc: ## Build the container without caching
	docker build --no-cache -t $(APP_NAME):$(TAG) .
	docker push $(APP_NAME):$(TAG)

run: ## Run container on port 5000
	docker run -p 80:8080/tcp $(APP_NAME):$(TAG) --name $(APP_NAME)

deploy: ## deploy using kubernetes
	kubectl create namespace s3write
	kubectl create -f s3write.yaml
