FROM ubuntu:20.04
RUN apt-get update
RUN apt-get install -y ca-certificates
ADD bin/fdic /
CMD ["/fdic"]
