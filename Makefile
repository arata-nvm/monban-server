ENV_LOCAL_FILE := env.local
ENV_LOCAL := $(shell cat $(ENV_LOCAL_FILE))

.PHONY: run
run:
	$(ENV_LOCAL) go run .