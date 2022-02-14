FROM golang:1.17-buster as go
# Set the working directory
WORKDIR /server
# Download dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download
# Import code from the context
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/app ./cmd/app/main.go

FROM alpine as app
# Set the working directory
WORKDIR /server
# Copy everything required from the go stage into the final stage
COPY --from=go /server/bin/app ./app
# Set the host and port
ENV S_PORT="8081"
EXPOSE 8081
EXPOSE 9000/udp
## Further actions will be performed as a non-privileged user
USER nobody:nobody
ENTRYPOINT ["./app"]