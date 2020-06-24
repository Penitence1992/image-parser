ARG goVersion=1.14.3-alpine


FROM golang:${goVersion} as builder
ARG gitCommit=""
ARG buildStamp=""

RUN mkdir /app

ADD . /app

RUN cd /app && \
    GO111MODULE=on go build -ldflags "-s -w -X 'main.gitCommit=${gitCommit}' -X 'main.buildStamp=${buildStamp}'" -o iparse pkg/cmd/cmd.go

FROM alpine

LABEL author=renjie email=penitence.rj@gmail.com

COPY --from=builder /app/iparse /usr/local/bin/

CMD ["ddns"]