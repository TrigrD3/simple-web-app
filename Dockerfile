FROM golang:1.19-alpine

WORKDIR /app

COPY main.go ./
RUN go mod init simple-web-app && go mod tidy
RUN go build -o /simple-web-app .

EXPOSE 8080

CMD [ "/simple-web-app" ]
