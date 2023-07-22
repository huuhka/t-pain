# Start with the build stage where we build our application
FROM ubuntu:22.04 AS build

# Update system and install required packages
RUN apt-get update && \
    apt-get install -y build-essential libssl-dev ca-certificates libasound2 wget curl git ffmpeg tzdata

# Set the environment variable for the Speech SDK location
ENV SPEECHSDK_ROOT="/usr/local/speechsdk"

# Download and install Speech SDK
RUN mkdir -p "$SPEECHSDK_ROOT" && \
    wget -O SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary && \
    tar --strip 1 -xzf SpeechSDK-Linux.tar.gz -C "$SPEECHSDK_ROOT"

# Install Golang version 1.20.6
RUN curl -O https://dl.google.com/go/go1.20.6.linux-amd64.tar.gz && \
    tar -xvf go1.20.6.linux-amd64.tar.gz && \
    rm go1.20.6.linux-amd64.tar.gz && \
    mv go /usr/local

# Set Go environment variables
ENV GOROOT="/usr/local/go"
ENV GOPATH="$HOME/go"
ENV PATH="$GOPATH/bin:$GOROOT/bin:$PATH"

# Set the CGO environment variables
ENV CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
ENV CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
ENV LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/t-pain/main.go

# Now create a second stage to only contain what's necessary to run the application
FROM ubuntu:22.04

# Install required packages for running the application
RUN apt-get update && apt-get install -y ca-certificates libasound2 ffmpeg tzdata && rm -rf /var/lib/apt/lists/*

# Copy the Speech SDK lib from the build stage
COPY --from=build /usr/local/speechsdk/lib/x64 /usr/local/speechsdk/lib/x64

# Set the CGO environment variables
ENV SPEECHSDK_ROOT="/usr/local/speechsdk"
ENV LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"

# Set the current working directory inside the container
WORKDIR /app

# Copy the binary from the build stage to the final stage
COPY --from=build /app/main /app/main

# Run the binary program produced by `go install`
CMD ["./main"]