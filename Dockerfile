FROM scratch

ENV LISTEN="127.0.0.1:26641"
ENV PROXY=""

COPY relaygram /usr/bin/relaygram

ENTRYPOINT ["/usr/bin/relaygram"]
