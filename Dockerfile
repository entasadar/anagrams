FROM golang:alpine as build
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o anagrams cmd/main.go

FROM alpine:latest
COPY --from=build /app/anagrams /bin/anagrams
CMD /bin/anagrams
EXPOSE 8080
