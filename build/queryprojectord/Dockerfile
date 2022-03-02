# Build go binary 
FROM golang AS builder

WORKDIR /go/src/github.com/pperaltaisern/financing

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/queryprojectord/...

# Build image
FROM alpine:latest

WORKDIR /go/src/github.com/pperaltaisern/financing

RUN apk update && apk add curl

COPY --from=builder /go/src/github.com/pperaltaisern/financing/queryprojectord .

EXPOSE 8080

CMD ["./queryprojectord"]