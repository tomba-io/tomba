FROM alpine
COPY tomba /usr/bin/tomba
ENTRYPOINT ["/usr/bin/tomba"]