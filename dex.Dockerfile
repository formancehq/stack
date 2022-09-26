FROM ghcr.io/dexidp/dex:v2.34.0
ENV DEX_FRONTEND_DIR=/app/web
COPY --chown=root:root pkg/web /app/web