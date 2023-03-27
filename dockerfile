# Build stage
FROM amd64/golang:latest AS builder

RUN apt-get update && apt-get install -y ca-certificates
    
RUN apt-get upgrade -y ca-certificates

# Install Chrome
RUN apt-get update && \
    apt-get install -y wget gnupg2 && \
    wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
COPY blank ./blank
COPY genExports ./genExports
COPY scrape ./scrape
COPY search ./search
COPY storage ./storage
COPY main.go ./

RUN go build -o main .

# Final stage
FROM amd64/alpine:latest
COPY --from=builder /app/main /app/main

# Set the default executable and entrypoint
ENTRYPOINT ["/app/main"]
CMD ["google-chrome-stable", "--no-sandbox", "--disable-gpu", "--headless", "--remote-debugging-address=0.0.0.0", "--remote-debugging-port=9222", "--disable-dev-shm-usage"]

EXPOSE 8080
