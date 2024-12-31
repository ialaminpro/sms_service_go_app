#!/bin/bash

# Determine the directory where the script is located
APP_PATH="$(cd "$(dirname "$0")" && pwd)"
APP_NAME="go_sms_service_app"
LOG_FILE="$APP_PATH/${APP_NAME}.log"

# Stop the application
if [ -f "$APP_PATH/${APP_NAME}.pid" ]; then
    PID=$(cat $APP_PATH/${APP_NAME}.pid)
    kill $PID
    rm "$APP_PATH/${APP_NAME}.pid" 
    echo "$(date '+%Y/%m/%d %H:%M:%S') Application stopped." >> $LOG_FILE
    echo "Application stopped."
else
    echo "$(date '+%Y/%m/%d %H:%M:%S') No PID file found. Is the app running?" >> $LOG_FILE
    echo "No PID file found. Is the app running?"
fi
