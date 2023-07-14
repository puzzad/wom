FROM golang:1.20 AS build

WORKDIR /go/src/app

COPY go.sum .
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -trimpath -ldflags=-buildid= -o /main ./cmd/serve

FROM gcr.io/distroless/base:nonroot
COPY --from=build /main /wom
COPY --from=build --chown=nonroot /go/src/app/templates /templates
WORKDIR /

ENTRYPOINT ["/wom"]
EXPOSE 3000 
