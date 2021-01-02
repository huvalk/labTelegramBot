FROM tetafro/golang-gcc:latest

COPY . /go/src/labTelegramBot

WORKDIR /go/src/labTelegramBot

RUN apk add make && make build-app

CMD /go/src/labTelegramBot/app