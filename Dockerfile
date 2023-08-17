# First stage 

FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache make

RUN make build

# Second stage

FROM node:20-alpine AS node_builder

WORKDIR /app

COPY package.json yarn.lock ./

RUN npm install -g npm@latest

RUN npm install -g yarn --force

RUN yarn

RUN apk add --no-cache make

COPY . .

RUN make build-css

# Third stage

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=node_builder /app/public ./public

COPY --from=node_builder /app/node_modules ./node_modules

COPY --from=builder /app/web .

COPY --from=builder /app/views ./views

EXPOSE 3000

CMD ["./web"]