FROM alpine:latest as builder

RUN apk update && apk add --no-cache libcap ca-certificates  && update-ca-certificates

# Create appuser
ENV USER=scratchuser
ENV UID=10001
ENV GID=10005
RUN addgroup --gid "$GID" "${USER}"

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --ingroup "${USER}" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
RUN mkdir -p /var/lib/ledis
RUN mkdir -p /var/lib/acme
# Copy executable to builder to setcap (allow to run on 80 and 443 if necessary)
# Default port still 5000
COPY pass-go /pass-go
RUN setcap 'cap_net_bind_service=+ep' /pass-go

############################

FROM scratch
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable
COPY --from=builder /pass-go /pass-go
# Ledis DB data dir
COPY --from=builder --chown="${USER}":"${USER}" /var/lib/ledis /var/lib/ledis
# Letsencrypt certificate cache
COPY --from=builder --chown="${USER}":"${USER}" /var/lib/acme /var/lib/acme

VOLUME ["/var/lib/ledis", "/var/lib/acme"]

# Use an unprivileged user.
USER "${USER}"
WORKDIR /var/lib/ledis

ENTRYPOINT ["/pass-go"]