FROM golang:1.17-stretch as build-stage

WORKDIR /app

COPY . /app

ENV CGO_ENABLED=0
RUN go build -o /app/build/bff-website ./main.go

FROM alpine

COPY --from=build-stage /app/build/bff-website /bin/bff-website

EXPOSE 8000

ENTRYPOINT ["/bin/bff-website"]
