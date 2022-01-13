FROM golang

RUN go get -u github.com/golang/dep/cmd/dep 

WORKDIR /go/src/github.com/arielizuardi/golang-backend-blog

COPY . .

RUN go build -o golang-backend-blog

EXPOSE 8080
CMD ./golang-backend-blog