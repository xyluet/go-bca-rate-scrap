IMAGE := bca-rate

.PHONY: docker
docker:
	docker build -t $(IMAGE):latest -f ./cmd/bca-rate/Dockerfile .
