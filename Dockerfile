# Start by building the application.
FROM golang:1.17 as build

WORKDIR /go/src/assumerole
COPY . .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 go build
RUN GOBIN=/go/bin go install -v ./...
RUN ls -l /go/bin

# Now copy it into our base image.
FROM ubuntu
RUN apt update; apt-get install ca-certificates -y; update-ca-certificates
COPY --from=build /go/bin/assumerole /
USER 3001:3001

CMD [ "ls -laR", "/var/run/secrets/eks.amazonaws.com/serviceaccount/" ]
CMD [ "/assumerole" ]
