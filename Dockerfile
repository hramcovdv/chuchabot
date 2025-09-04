FROM golang:1.25-alpine AS build
WORKDIR /src/
COPY . .
RUN go mod download \
&& go mod verify \
&& go build -v -o chuchabot main.go

FROM alpine:3.22
COPY --from=build /src/chuchabot /bin/chuchabot
ENTRYPOINT ["/bin/chuchabot"]
