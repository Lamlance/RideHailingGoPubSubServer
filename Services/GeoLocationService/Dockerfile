FROM golang:1.19.13-alpine3.18
WORKDIR /app
COPY . .
ENV PORT=3083
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_HOST=rabbitmq
ENV REDIS_PORT=6785
ENV REDIS_HOST=redis
EXPOSE 3083
RUN go mod tidy
RUN go build -o geo.exe
CMD [ "./geo.exe" ]
