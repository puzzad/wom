FROM golang:1.20 AS build

WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go install ./cmd/serve

FROM gcr.io/distroless/base:nonroot
COPY --from=build /go/bin/serve /wom
COPY --from=build --chown=nonroot /go/src/app/templates /templates
WORKDIR /

ENTRYPOINT ["/wom"]
EXPOSE 3000 
