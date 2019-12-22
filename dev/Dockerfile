FROM golang:1.13.5-buster

RUN curl -fLo /usr/local/bin/air \
      https://raw.githubusercontent.com/cosmtrek/air/master/bin/linux/air \
    && chmod +x /usr/local/bin/air

ENV GO111MODULE "on"
