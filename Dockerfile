FROM gcr.io/distroless/static
COPY kaspa_exporter /usr/local/bin/kaspa_exporter
ENTRYPOINT [ "/usr/local/bin/kaspa_exporter" ]