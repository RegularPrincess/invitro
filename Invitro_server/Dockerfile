FROM golang
 
ADD . /go/src/invitro_server
RUN go get github.com/lib/pq
RUN go install invitro_server
ENTRYPOINT /go/bin/invitro_server
 
EXPOSE 8080
