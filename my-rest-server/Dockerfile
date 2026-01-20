##
## Build
##

FROM golang:latest AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /my-rest-server
EXPOSE 8080

CMD [ "/my-rest-server" ]


##
## Deploy
## todo: will figure out it later onwards, since dealing

# FROM gcr.io/distroless/base-debian12

# WORKDIR /

# COPY --from=build /my-rest-server /my-rest-server
# COPY --from=build /app/conf /conf
# COPY --from=build /app/assets /assets

# EXPOSE 8080

# USER nonroot:nonroot

# ENTRYPOINT ["/my-rest-server"]
