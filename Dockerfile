# Start by building the application.
FROM golang:1.12 as build

WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build
#RUN go get -d -v ./...
#RUN go install -v ./...

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/go-deploy /
ENV SRC_DIR=/src
RUN mkdir -p /tmp /src /html_dir
ENTRYPOINT [ "/go-deploy"]
CMD [ "volume"]