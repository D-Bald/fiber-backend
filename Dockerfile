FROM golang:latest

RUN mkdir /app
WORKDIR /app
ENV DB_USER=${DB_USER}
ENV DB_USER_PASSWORD=${DB_USER_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV MONGO_PORT=${MONGO_PORT}
ENV FIBER_PORT=${FIBER_PORT}
ENV SECRET=${SECRET}
ENV ADMIN_PASSWORD=${ADMIN_PASSWORD}

COPY . .
RUN go build -o main .
EXPOSE ${FIBER_PORT}

CMD ["/app/main"]