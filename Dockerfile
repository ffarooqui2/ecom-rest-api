# use officieal Goland image
FROM golang:1.16.3

# set working directory
WORKDIR /app

# copy the source code
COPY . .

# download and install the dependencies
RUN go get -d -v ./...

# build the application
RUN go build -o api .

# expose the port
EXPOSE 8080

# run the executable
CMD ["./api"] 