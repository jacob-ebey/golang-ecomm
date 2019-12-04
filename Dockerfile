# Build the frontend
FROM mhart/alpine-node:12 as nodebuilder
WORKDIR /app

COPY frontend .

RUN npm install
RUN npm run build


# Build the backend
FROM golang:latest as gobuilder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o program ./main/main.go


# The acutal application
FROM alpine:latest

WORKDIR /root

RUN apk --no-cache add ca-certificates

COPY --from=nodebuilder /app/build ./frontend

COPY --from=gobuilder /app/program .
RUN chmod +x ./program

COPY --from=gobuilder /app/wait-for-it.sh .
RUN chmod +x ./wait-for-it.sh
RUN apk add --no-cache bash

EXPOSE 8080

CMD ./wait-for-it.sh ${POSTGRESS_ADDRESS} -- ./program