# Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY . .

# Install dependencies for using chrome and chromedriver
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        wget \
        unzip \
        libxss1 \
        libappindicator1 \
        libnss3 \
        libasound2 \
        libgconf-2-4 \
        lsb-release \
        fonts-liberation \
        xdg-utils \
        openjdk-11-jre-headless

# Install Google Chrome
RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt install -y ./google-chrome-stable_current_amd64.deb && \
    rm google-chrome-stable_current_amd64.deb

RUN mv /usr/bin/google-chrome-stable /usr/local/bin/google-chrome

# Install chromedriver
RUN wget -N https://chromedriver.storage.googleapis.com/$(curl -sS chromedriver.storage.googleapis.com/LATEST_RELEASE)/chromedriver_linux64.zip && \
    unzip chromedriver_linux64.zip && \
    chmod +x chromedriver && \
    mv chromedriver /usr/local/bin/ && \
    rm chromedriver_linux64.zip

RUN wget https://selenium-release.storage.googleapis.com/3.141/selenium-server-standalone-3.141.59.jar && \
    echo "java -jar selenium-server-standalone-3.141.59.jar &" > /usr/local/bin/start-selenium && \
    chmod +x /usr/local/bin/start-selenium

# Build the Go binary
RUN go mod tidy
RUN go build -o main .

# Start a new stage to minimize the final image size
FROM gcr.io/distroless/base-debian10

# Copy files from the builder stage
COPY --from=builder /usr/local/bin/google-chrome /usr/local/bin/google-chrome
COPY --from=builder /usr/local/bin/chromedriver /usr/local/bin/chromedriver
COPY --from=builder /app/main /app/main
COPY --from=builder /usr/local/bin/start-selenium /usr/local/bin/start-selenium
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Set the default executable and entrypoint
WORKDIR /app
ENTRYPOINT ["/usr/local/bin/start-selenium"]
CMD ["-port", "4444"]

# Expose the Selenium server port
EXPOSE 4444
EXPOSE 8080
