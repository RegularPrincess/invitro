FROM golang
 
ADD . /go/src/github.com/RegularPrincess/invitro
RUN go get github.com/lib/pq
RUN go install github.com/RegularPrincess/invitro/invitro_model
RUN go get github.com/PuerkitoBio/goquery
RUN go get github.com/djimenez/iconv-go
RUN go install github.com/RegularPrincess/invitro/invitro_parser
RUN go install github.com/RegularPrincess/invitro/invitro_server
ENTRYPOINT /go/bin/invitro_server
 
EXPOSE 8080
