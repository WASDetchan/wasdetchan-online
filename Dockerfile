FROM golang:1.26 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /usr/local/bin/app

FROM ubuntu

COPY --from=build /usr/local/bin/app /usr/local/bin/app 

CMD ["/usr/local/bin/app"]
