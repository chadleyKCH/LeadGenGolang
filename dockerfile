# Stage 1: Build the Go program
FROM golang:1.20-alpine AS build

WORKDIR /app

COPY . .

RUN go build -o main .

# Stage 2: Create the runtime environment
FROM alpine:3.15

# Install necessary dependencies
RUN apk add --no-cache \
    openjdk11-jre \
    wget \
    unzip \
    curl \
    gnupg \
    ttf-dejavu \
    libstdc++ \
    libx11 \
    libxcomposite \
    libxrender \
    libxcursor \
    libxi \
    libxtst \
    libxrandr \
    libxscrnsaver \
    libxext \
    libxfixes \
    ca-certificates \
    chromium \
    chromium-chromedriver

# Set the environment variable for Chrome
ENV CHROME_BIN=/usr/bin/chromium

# Install Selenium standalone server
RUN mkdir /opt/selenium && \
    curl -sS -o /opt/selenium/selenium-server-standalone.jar http://selenium-release.storage.googleapis.com/3.141/selenium-server-standalone-3.141.59.jar

# Copy the Go binary from the build stage
COPY --from=build /app/main /app/main

# Expose the Selenium standalone server port (default is 4444)
EXPOSE 4444
EXPOSE 4445

WORKDIR /app

# Start the Selenium standalone server and your Go program
CMD java -jar /opt/selenium/selenium-server-standalone.jar & sleep 10 && ./main
