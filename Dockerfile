# First stage 

FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o web cmd/web/main.go

# Second stage

FROM node:20-alpine AS node_builder

WORKDIR /app

COPY package.json yarn.lock ./

RUN npm install -g npm@latest

RUN npm install -g yarn --force

RUN yarn

COPY . .

RUN yarn run build-css

# Third stage

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=node_builder /app/public ./public

COPY --from=builder /app/web .

COPY --from=builder /app/views ./views

RUN touch .env

EXPOSE 3000

CMD ["./web"]