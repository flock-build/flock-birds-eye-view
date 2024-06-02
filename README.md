# flock-birds-eye-view

## Prerequisite
- Conda (with Python 3.12 highly recommended)
- Go (can be installed via Conda)
- Docker

## Setup
1. Clone the repository
2. Run [RabbitMQ on Docker](https://www.rabbitmq.com/docs/download)
3. Run `python main.py` in processor
4. Run `go run .` in proxy

## How to use proxy?

### Typical OpenAI API call:
```bash
curl -X POST https://api.openai.com/v1/chat/completions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer sk-..." \
-d '{
  "model": "gpt-3.5-turbo-instruct",
  "prompt": "Say this is a test",
  "max_tokens": 7,
  "temperature": 0
}'
```

### Flocked OpenAI API call:
```bash
curl -X POST http://localhost:8080/openai/v1/completions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer sk-..." \
-H "FLOCK-AUTH: flk-..." \
-d '{
  "model": "gpt-3.5-turbo-instruct",
  "prompt": "Say this is a test",
  "max_tokens": 7,
  "temperature": 0
}'
```