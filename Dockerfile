FROM golang:1.22.5 AS build
WORKDIR /app
COPY ./ ./

RUN go build -o /bin/spotipub .
CMD ["/bin/spotipub"]