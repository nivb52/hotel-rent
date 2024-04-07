FROM golang:1.20.3-alpine

WORKDIR /app

COPY .env ./ 
COPY .env.local ./
ENV IS_DOCKERIZE=true
ENV HOSTNAME=":5000"
ENV DB_CONNECTION_STRING=mongodb://root:myrootpassword@mongo:27018


COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 5000

CMD [ "./main" ]
