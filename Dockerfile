#- Build Stage
FROM golang:1.18.0-alpine3.15 AS builder

WORKDIR /go/src/github.com/3n0ugh/movpic
COPY . .

RUN echo ls
RUN go mod download
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -trimpath -a -v -o /go/bin/movpic

#- Run Stage
# FROM gcr.io/distroless/static-debian11
FROM golang:1.18.0-alpine3.15

WORKDIR /app

COPY --from=builder /go/bin/movpic /
COPY --from=builder /go/src/github.com/3n0ugh/movpic/framework/swagger/dist /app/framework/swagger/dist

EXPOSE 8080

CMD [ "/movpic" ]
