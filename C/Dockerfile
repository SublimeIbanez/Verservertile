FROM alpine:latest

WORKDIR /app

RUN apk update && apk add alpine-sdk

COPY . /c_server/

CMD ["sh"]
# CMD ["sh", "-c", "gcc -o myapp main.c && ./myapp"]