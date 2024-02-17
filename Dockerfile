#Build stage
FROM golang:1.21-bullseye AS builder

RUN apt-get update

WORKDIR /wishlist

COPY . .
 
COPY .env .env

RUN go mod download

RUN go build -o ./out/dist ./cmd

#production stage

FROM busybox

RUN mkdir -p /wishlist/out/dist

COPY --from=builder /wishlist/out/dist /wishlist/out/dist

COPY --from=builder /wishlist/.env /wishlist/out

WORKDIR /wishlist/out/dist

EXPOSE 8085

CMD ["./dist"]