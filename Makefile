APP_NAME=quay.io/cporter/assumerole
EXEC_NAME=assumerole

BUILD_NUMBER_FILE=buildnum.txt

.PHONY: $(BUILD_NUMBER_FILE)

$(BUILD_NUMBER_FILE):
	@if ! test -f $(BUILD_NUMBER_FILE); then echo 1000 > $(BUILD_NUMBER_FILE); fi
	@echo $$(($$(cat $(BUILD_NUMBER_FILE)) + 1)) > $(BUILD_NUMBER_FILE)

image: TAG = 1.$(shell cat $(BUILD_NUMBER_FILE))

image: $(BUILD_NUMBER_FILE)
	$(info Starting image build with tag $(TAG))
	docker build -t $(APP_NAME):$(TAG) .
	docker push $(APP_NAME):$(TAG)
	## rewrite YAML
	sed -E -i .bak 's/:[0-9]\.[0-9]{4}/:$(TAG)/g' assumerole.yaml

run: ## Run container on port 5000
	docker run -p 80:8080/tcp $(APP_NAME):$(TAG) --name $(APP_NAME)

deploy: ## deploy using kubernetes
	kubectl create namespace s3write
	kubectl create -f s3write.yaml
