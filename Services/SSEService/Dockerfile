FROM golang:latest
WORKDIR /app
COPY . .
ENV PORT=3082
ENV RABBITMQ_PORT=5672
ENV RABBITMQ_HOST=rabbitmq
EXPOSE 3082
RUN go mod tidy
RUN go build -o sse.exe
CMD [ "./sse.exe" ]
