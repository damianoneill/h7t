FROM scratch
COPY h7t ./
COPY transformer ./plugins/csv/transformer
ENTRYPOINT ["./h7t"]
