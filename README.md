## Trivial Telegram BOT with STT for Voice Messages

[![Go](https://github.com/MarcusAdriano/tgbot-stt/actions/workflows/go.yml/badge.svg)](https://github.com/MarcusAdriano/tgbot-stt/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/MarcusAdriano/tgbot-stt/branch/main/graph/badge.svg?token=8EiBcDQOPO)](https://codecov.io/github/MarcusAdriano/tgbot-stt)

It's a very simple bot thats converts voice messages to text using openai's whisper API.

## Goals

- [x] Improve my golang skills
- [x] Unit testing with mocks

## Getting started

### Prerequisites

- [Go](https://golang.org/doc/install)
- [OpenAI API Key](https://platform.openai.com/account/api-keys)
- [Telegram Bot Token](https://core.telegram.org/bots#6-botfather)
- [Docker](https://docs.docker.com/get-docker/)

### Running

#### Using Docker

Fill .env file with your OpenAI API Key and Telegram Bot Token.

Build the docker image:

```bash
make docker-build
```

Run the container:

```bash
docker run --env-file .env -d --name stt-bot marcusadriano/tgbot-stt-chatgpt:latest
```


## Reference
- [OpenAI's Whisper API](https://platform.openai.com/docs/api-reference/audio)
- [Telegram Bot API](https://core.telegram.org/bots/api)
