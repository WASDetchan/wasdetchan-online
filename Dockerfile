FROM golang:1.26 AS build

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY pages/ ./pages/
COPY *.templ ./

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v -o /usr/local/bin/app




FROM node:25-slim AS styles

WORKDIR /usr/src/app

COPY package.json ./
RUN npm install

COPY css ./css
COPY pages/*.templ ./pages/
COPY *.templ ./
RUN npx @tailwindcss/cli -i css/app.css -o /static/app.css




FROM ubuntu

COPY --from=build /usr/local/bin/app /usr/local/bin/app 
COPY --from=styles /static/app.css /static/app.css

CMD ["/usr/local/bin/app"]
