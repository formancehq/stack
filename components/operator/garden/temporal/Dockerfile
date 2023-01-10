FROM golang:1.18 AS builder

RUN mkdir -p /go/src/temporalite
RUN git clone https://github.com/temporalio/temporalite.git /go/src/temporalite
WORKDIR /go/src/temporalite
RUN git switch -c tmp v0.2.0
RUN go mod download
RUN go get -d -v ./...
RUN go build -o /temporalite ./cmd/temporalite
EXPOSE 7233 8233

ENTRYPOINT ["/temporalite", "start", "--ephemeral", "-n", "default", "--ip" , "0.0.0.0"]
