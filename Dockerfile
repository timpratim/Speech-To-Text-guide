# Use the official Golang base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Install whisper
RUN git clone https://github.com/ggerganov/whisper.cpp.git &&\
    cd whisper.cpp && make &&\
    make libwhisper.so libwhisper.a &&\
    cp whisper.h /usr/local/include &&\
    cp ggml.h /usr/local/include &&\
    cp libwhisper.a /usr/local/lib &&\
    cp libwhisper.so /usr/local/lib &&\
    cd ..

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all necessary dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go app
RUN go build -o ytt

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# execute cli help to check if everything is ok
RUN ./ytt -h


# Run the compiled binary with a default command
ENTRYPOINT ["/app/ytt"]