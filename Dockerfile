FROM golang:1.16.0-alpine3.13
LABEL author="Ruslan Kasimov"
ENV NAME "history"
WORKDIR /opt/${NAME}
COPY ./go.mod .
RUN go mod download
RUN go get -u github.com/swaggo/swag/cmd/swag
COPY . .

ARG DOMAIN
ENV DOMAIN=${DOMAIN}

RUN swag init -g ./main.go -o "docs" --md . && \
    sed -i -e "s#{{ IM_HOST_CHANGE_ME }}#${DOMAIN}#g" "/opt/${NAME}/docs/swagger.yaml" && \
    sed -i -e "s#{{ IM_HOST_CHANGE_ME }}#${DOMAIN}#g" "/opt/${NAME}/docs/swagger.json" && \
    sed -i -e "s#{{ IM_HOST_CHANGE_ME }}#${DOMAIN}#g" "/opt/${NAME}/docs/docs.go"

RUN CGO_ENABLED=0 GOARCH=amd64 go build -o "/opt/${NAME}/bin/${NAME}" "/opt/${NAME}/main.go"

FROM alpine:3.13
LABEL author="Ruslan Kasimov"
ARG ENV
ENV NAME "history"
WORKDIR /opt/${NAME}
COPY --from=0 /opt/${NAME}/bin/. ./.
CMD ./${NAME}