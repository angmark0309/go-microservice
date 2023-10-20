FROM alpine:latest
RUN mkdir /app

COPY logger-service/logServiceApp /app

# Run the server executable
CMD [ "/app/logServiceApp" ]