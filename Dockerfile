# PS: Здесь я использовал multi stage сборку что уменишить конечный размер образа
# 1 stage - я поделил его на слои, так что бы go mod download
#кэшировалось и не занимала время при сборке если не менялось
# 2 stage - я использовал alpine что бы уменшить размер, из первого стеджа нам нужнен толкьо бинарь и енвы

FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -v -o test_task_kami

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/test_task_kami .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./test_task_kami"]



