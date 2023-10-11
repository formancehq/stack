FROM ghcr.io/dexidp/dex:v2.37.0
ENV DEX_FRONTEND_DIR=/app/web
COPY --chown=root:root pkg/web /app/web
