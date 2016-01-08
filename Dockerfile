FROM alpine:3.2
RUN apk add -U ca-certificates git openssh curl perl && rm -rf /var/cache/apk/*
ADD gister /bin/
ENTRYPOINT ["/bin/gister"]
