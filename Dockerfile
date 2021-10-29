# Start by building the application.
FROM golang:1.17 as build

WORKDIR /go/src/assumerole
COPY . .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 go build src/assumerole.go
RUN GOBIN=/go/bin go install -v ./...
RUN ls -l /go/bin

# Now copy it into our base image.
FROM redhat/ubi8
RUN dnf -y update
COPY --from=build /go/src/assumerole/assumerole /usr/bin
USER 3001:3001

CMD [ "/usr/bin/assumerole" ]
