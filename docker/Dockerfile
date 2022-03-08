FROM gesquive/go-builder:latest AS builder

ENV APP=paperless-uploader

# This requires that `make release-snapshot` be called first
COPY dist/ /dist/
RUN copy-release
RUN chmod +x /app/paperless-uploader

COPY docker/config.yml /etc/paperless-uploader

# =============================================================================
FROM gesquive/docker-base:busybox
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

# Import from builder
COPY --from=builder /app/paperless-uploader /app/
COPY --from=builder /etc/paperless-uploader/ /etc/paperless-uploader/

WORKDIR /config
VOLUME /config
volume /watch

ENTRYPOINT ["run"]
CMD ["/app/paperless-uploader"]
