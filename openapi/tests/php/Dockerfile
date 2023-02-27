FROM php:8.1.16
WORKDIR /src
RUN apt update && apt install git -y
RUN curl -sS https://getcomposer.org/installer -o /tmp/composer-setup.php
RUN php /tmp/composer-setup.php --install-dir=/usr/local/bin --filename=composer

COPY sdks/php sdks/php
COPY openapi/tests/php openapi/tests/php

WORKDIR openapi/tests/php
ENTRYPOINT ["/src/openapi/tests/php/entrypoint.sh"]
