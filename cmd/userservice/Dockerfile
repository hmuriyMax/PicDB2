FROM golang:1.18

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY /home/maxim/GoProjects/PicDB2/go.mod /home/maxim/GoProjects/PicDB2/go.sum ./
RUN go mod init userservice && go mod tidy
RUN go get github.com/hmuriyMax/PicDB2/pkg/user

COPY . .
RUN go build -v -o /usr/local/bin/app ./...

EXPOSE 80

CMD ["app"]