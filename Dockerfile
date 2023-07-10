FROM golang:1.20.5-alpine3.18

ENV GOPROXY=https://proxy.golang.org,direct
ENV GONOPROXY=example.com/private

WORKDIR /app
COPY . .
RUN go build -o main main.go

EXPOSE 8080
CMD [ "/app/main" ]