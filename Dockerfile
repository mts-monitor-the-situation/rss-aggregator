FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY mts-rss-aggregator .
USER 65532:65532

ENTRYPOINT [ "./mts-rss-aggregator" ]