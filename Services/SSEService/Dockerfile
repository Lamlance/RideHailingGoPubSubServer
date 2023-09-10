FROM golang:latest
WORKDIR /app
COPY . .
ENV PORT=3082
ENV RABBITMQ=5672
EXPOSE 3082
RUN go mod tidy
RUN go build -o sse.exe
CMD [ "./sse.exe" ]
