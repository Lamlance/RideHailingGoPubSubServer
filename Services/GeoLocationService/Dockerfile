FROM golang:latest
WORKDIR /app
COPY . .
ENV PORT=3083
ENV RABBITMQ=5672
EXPOSE 3083
RUN go mod tidy
RUN go build -o geo.exe
CMD [ "./geo.exe" ]
