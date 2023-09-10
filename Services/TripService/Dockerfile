FROM golang:latest
WORKDIR /app
COPY . .
ENV PORT=3081
ENV RABBITMQ=5672
EXPOSE 3081
RUN go mod tidy
RUN go build -o trip.exe
CMD [ "./trip.exe" ]
