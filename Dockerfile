FROM golang AS pipeline
#RUN mkdir -p /go/src/website
WORKDIR /go/src/module20
ADD main.go .
ADD go.mod .
RUN go install .

FROM alpine:latest
LABEL version="1.0.0"
LABEL maintainer="Alex Jevo<loisapandohva@gmail.com>"
WORKDIR /root/
COPY --from=pipeline /go/bin/pipeline .
ENTRYPOINT ./pipeline
EXPOSE 8080