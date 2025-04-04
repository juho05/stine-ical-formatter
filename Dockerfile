# Login: docker login ghcr.io -u juho05
# Build and deploy: docker buildx build --platform linux/arm64,linux/amd64 --tag ghcr.io/juho05/stine-ical-formatter:latest --push .

# Build
FROM --platform=$BUILDPLATFORM golang:alpine AS build
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN apk add --no-cache make nodejs npm

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} make clean && make

# Run
FROM alpine AS stine-ical-formatter
ARG BUILDPLATFORM
WORKDIR /
COPY --from=build /app/bin/stine-ical-formatter /stine-ical-formatter

EXPOSE 8080

CMD [ "/stine-ical-formatter" ]
