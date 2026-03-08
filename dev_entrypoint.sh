#!/bin/sh
set -a
. /.env.dev
set +a
exec /bin/muxly-msg-subscriber
