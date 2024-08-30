FROM golang:1.23 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .env *.go ./
COPY static/ ./static
COPY templates/ ./templates
RUN ls --recursive ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /server /app/.env /
EXPOSE 9090
USER nonroot:nonroot
ENTRYPOINT ["/server"]
