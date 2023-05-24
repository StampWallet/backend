FROM golang:alpine as build

WORKDIR /usr/src/app
COPY . .
RUN go mod download && go install github.com/golang/mock/mockgen@v1.6.0
RUN go generate ./... && go build -v -o /usr/local/bin/app ./cmd

FROM alpine:latest
COPY --from=build /usr/local/bin/app /
CMD /app automigrate && /app start
