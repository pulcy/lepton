FROM alpine:3.4

ADD ./lepton /bin/

ENTRYPOINT ["/bin/lepton"]
