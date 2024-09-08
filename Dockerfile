FROM golang:1.22-rc-alpine
WORKDIR /app
COPY go.mod ./

RUN go mod download
#RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
COPY . ./

RUN go build -o /go_telegram_bot

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8080
# command to be used to execute when the image is used to start a container
CMD [ "/go_telegram_bot" ]

#
#
#FROM golang:1.22-rc-alpine as build
#WORKDIR /app
#
#RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates
#
#COPY . ./
#
#RUN go mod download
#
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/app ./cmd/app
#
#FROM scratch
#
#WORKDIR /app
#
#COPY --from=build /app/bin/app ./bin/
#COPY --from=build /app/.env ./.env
#COPY --from=build /app/database/ ./database/
#
#ENTRYPOINT ["./bin/app"]
#CMD [ "/go_telegram_bot" ]
# tells Docker that the container listens on specified network ports at runtime
#EXPOSE 8080
# command to be used to execute when the image is used to start a container
#CMD [ "/avito_test/cmd/app" ]