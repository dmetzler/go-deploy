# Start by building the application.
FROM golang:1.12 as build
ENV GO111MODULE=on
WORKDIR /go/src/go-deploy
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...
RUN mkdir -p /dest/tmp /dest/src /dest/html_dir


# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /dest /
COPY --from=build /go/bin/go-deploy /
ENV SRC_DIR=/src

ENTRYPOINT [ "/go-deploy"]
CMD [ "volume"]