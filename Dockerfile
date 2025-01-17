FROM alpine
COPY stream-forward /usr/bin/example
ENTRYPOINT ["/usr/bin/example"]