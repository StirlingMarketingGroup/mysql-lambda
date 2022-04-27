# Run the following command to create the AWS layer zip in the same directory as this Dockerfile
# docker build -o . .

FROM ubuntu:20.04 as build-stage

RUN apt-get update && \
  apt-get install -y libmysqlclient-dev \
  gcc \
  wget

RUN wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz && \
  rm -rf /usr/local/go && \
  tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

COPY main.go go.mod go.sum ./
RUN go build -buildmode=c-shared -o mysql_lambda.so

FROM scratch AS export-stage
COPY --from=build-stage /mysql_lambda.so /mysql_lambda.so