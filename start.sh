#!/bin/bash

# Determine the directory where the script is located
APP_PATH="$(cd "$(dirname "$0")" && pwd)"
APP_NAME="go_sms_service_app"
LOG_FILE="$APP_PATH/${APP_NAME}.log"

# Start the application using nohup
nohup $APP_PATH/${APP_NAME} > $LOG_FILE 2>&1 &
echo $! > $APP_PATH/${APP_NAME}.pid
echo "Application started, PID: $(cat $APP_PATH/${APP_NAME}.pid)"
