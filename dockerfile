# Use the official Golang image as the base image
FROM golang:latest

# Install necessary dependencies
RUN apt-get update && \
    apt-get install -y wget unzip curl default-jdk

# Install Google Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable

# Install ChromeDriver
RUN CHROMEDRIVER_VERSION=`curl -sS chromedriver.storage.googleapis.com/LATEST_RELEASE` && \
    mkdir -p /opt/chromedriver-$CHROMEDRIVER_VERSION && \
    curl -sS -o /tmp/chromedriver_linux64.zip http://chromedriver.storage.googleapis.com/$CHROMEDRIVER_VERSION/chromedriver_linux64.zip && \
    unzip -qq /tmp/chromedriver_linux64.zip -d /opt/chromedriver-$CHROMEDRIVER_VERSION && \
    rm /tmp/chromedriver_linux64.zip && \
    chmod +x /opt/chromedriver-$CHROMEDRIVER_VERSION/chromedriver && \
    ln -fs /opt/chromedriver-$CHROMEDRIVER_VERSION/chromedriver /usr/local/bin/chromedriver

# Install Selenium standalone server
RUN mkdir /opt/selenium && \
    curl -sS -o /opt/selenium/selenium-server-standalone.jar http://selenium-release.storage.googleapis.com/3.141/selenium-server-standalone-3.141.59.jar

# Set the working directory
WORKDIR /app

# Copy your Go program's source code
COPY . .

# Build your Go program
RUN go build -o main .

# Expose the Selenium standalone server port (default is 4444)
EXPOSE 4444
EXPOSE 4445

# Start the Selenium standalone server and your Go program
CMD java -jar /opt/selenium/selenium-server-standalone.jar & sleep 10 && ./main
