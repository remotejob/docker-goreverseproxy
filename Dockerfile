FROM scratch
EXPOSE 8080

COPY server /
ENTRYPOINT ["/server"]
