FROM ubuntu:latest

# apt-get stuff
RUN mkdir -p /opt/build/src
WORKDIR /opt/build/

RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y nodejs npm wget
# firefox stuff
RUN apt-get install -y wget bzip2 libxtst6 libgtk-3-0 libx11-xcb-dev libdbus-glib-1-2 libxt6 libpci-dev && rm -rf /var/lib/apt/lists/*
RUN wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin
# RUN ls /usr/local/go/bin/
RUN go version
RUN node --version

# Build frontend
COPY src/ /opt/build/src/
WORKDIR /opt/build/src/frontend
RUN npm run build

# Build golang
WORKDIR /opt/build/src/
RUN go build

# copy go to opt/app
RUN mkdir -p /opt/app
RUN cp .env /opt/app
RUN cp ./ulysses /opt/app
# copy frontend to opt/app
RUN cp -r /opt/build/src/frontend /opt/app/frontend
WORKDIR /opt/app

CMD ["./ulysses"]