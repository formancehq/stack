FROM caddy:2.6.2-builder as builder

# TODO(paul): fix commit to main version
# NOTE: Using a `=` in a `--with` argument will create a replace directive in
# the go.mod file.
RUN xcaddy build v2.6.2 \
    --with github.com/formancehq/stack/components/gateway/pkg/plugins@082579ea \
    --with github.com/ggicci/caddy-jwt=github.com/formancehq/caddy-jwt@c8d9423

FROM caddy:2.6.2

COPY --from=builder /usr/bin/caddy /usr/bin/caddy