FROM arm32v7/golang AS builder

RUN go get github.com/tarm/serial && go get github.com/mattn/go-sqlite3 && go get github.com/gorilla/mux
RUN mkdir /app
ADD ./src /app/ 
WORKDIR /app 
RUN go build -ldflags "-linkmode external -extldflags -static" -o main . 

FROM alpine:latest AS runner 
#RUN apk --no-cache add ca-certificates
WORKDIR /app
RUN mkdir db
COPY --from=builder /app/main .
#--build-arg <varname>=<value>
ARG listen_port=3000 
ARG serial_device=/dev/ttyUSB0 
ARG monitor_interval=37000 
ARG writer_interval=97000 
ARG db_path=/app/db/ups.sqlite
#--env <key>=<value>
ENV LISTEN_PORT=${listen_port}
ENV SERIAL_DEVICE=${serial_device}
ENV MONITOR_INTERVAL=${monitor_interval}
ENV WRITER_INTERVAL=${writer_interval}
ENV DB_PATH=${db_path}
EXPOSE ${LISTEN_PORT}
CMD ["/app/main"]
