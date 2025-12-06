#!/bin/sh

# if health endpoint is disabled, this fallbacks to checking if http will answer at all
if [ "$FAILBOOK_HEALTH_ENABLED" = "false" ]; then
    exec wget --no-verbose --tries=1 --spider http://localhost:12001
fi

exec wget --no-verbose --tries=1 --spider http://localhost:12001/manage/health/live
