# syntax=docker.io/docker/dockerfile:1

# Build stage: static binary, no CGO, so the final image can be scratch.
FROM docker.io/library/golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /bin/site .

# Runtime stage: nothing but the binary, its assets and a non-root user.
# scratch is a reserved empty base (not a registry image), so no registry prefix.
FROM scratch
COPY --from=build /bin/site /site
COPY views/ /views
COPY static/ /static
COPY content/ /content
ENV PORT=8080
EXPOSE 8080
USER 65534:65534
ENTRYPOINT ["/site"]
