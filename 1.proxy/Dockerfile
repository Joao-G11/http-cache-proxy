FROM golang:1.17-alpine


WORKDIR /proxy

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY proxy.go ./
COPY proxy.yaml ./

RUN go build -o proxy

CMD [ "./proxy" ]