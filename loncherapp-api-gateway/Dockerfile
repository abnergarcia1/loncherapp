# Building stage
FROM golang:1.14-alpine3.11 as builder

ENV GOPRIVATE=bitbucket.org/edgelabsolutions/*


RUN mkdir - /go/src/api
WORKDIR /go/src/api

ENV GO111MODULE=on

COPY go.sum go.mod ./



COPY app app

RUN apk add --no-cache git curl build-base bash openssh-client shadow


RUN printf "machine bitbucket.org\n\
            login ${USERNAME}\n\
            password ${PASSWORD}\n\
            \n\
            machine api.bitbucket.org\n\
            login ${USERNAME}\n\
            password ${PASSWORD}\n"\
            >> /root/.netrc
RUN chmod 600 /root/.netrc


RUN go mod download

RUN go build -o /api app/main.go


# Final image
FROM alpine:3.11


COPY --from=builder /api .

CMD ["./api"]