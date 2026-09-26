FROM golang:1.27-alpine AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/http-server .

FROM scratch
COPY --from=build /out/http-server /http-server
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/http-server"]
