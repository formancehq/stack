#!/bin/sh

set -x;
set -e;

composer update;
php index.php;
