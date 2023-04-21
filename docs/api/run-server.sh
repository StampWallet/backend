#!/bin/sh
cp ../../api/schema.yaml .
python3 -m http.server 8000 --directory .
