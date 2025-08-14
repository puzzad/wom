FROM golang:1.25 AS build

WORKDIR /go/src/app

COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -trimpath -ldflags=-buildid= -o /main ./cmd/serve
RUN mkdir /rootfs; \
    mkdir /rootfs/migrations; \
    mkdir /rootfs/data; \
    cp /main /rootfs/wom
RUN mkdir /migrations; mkdir /data

FROM gcr.io/distroless/base:nonroot
COPY --from=build --chown=nonroot /rootfs /
ENTRYPOINT ["/wom"]
EXPOSE 8090
