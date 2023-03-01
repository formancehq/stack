#!/bin/sh

set -x;
set -e;

pip install -r requirements.txt;
python test.py
