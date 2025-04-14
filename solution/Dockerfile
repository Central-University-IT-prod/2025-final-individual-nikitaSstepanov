FROM golang:alpine AS builder

EXPOSE 80

WORKDIR /usr/local/src

COPY ["go.mod", "go.sum", "./"]

RUN ["go", "mod", "download"]

COPY ./ ./

RUN ["go", "build", "-o", "./bin/advert-platform", "./cmd/advert-platform"]

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/advert-platform ./

COPY ./config ./config

COPY ./migrations ./migrations

CMD [ "./advert-platform" ]
