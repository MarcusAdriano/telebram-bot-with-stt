APP_NAME=tgbot-stt-whisper


build:
	go test ./...
	CGO_ENABLED=0 GOOS=linux go build -o bin/$(APP_NAME)

docker-build:
	docker build -t marcusadriano/tgbot-stt-chatgpt:latest .