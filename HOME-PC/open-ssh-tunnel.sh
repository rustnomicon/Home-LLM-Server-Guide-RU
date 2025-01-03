#!/bin/sh

if [ -z "$SSHPASS" ] || [ -z "$SSH_USER" ] || [ -z "$SSH_HOST" ] || [ -z "$PROXY_PORT_APP" ] || [ -z "$PROXY_PORT_SERVER" ]; then
  echo "Ошибка: убедитесь, что переменные окружения SSHPASS, SSH_USER и SSH_HOST установлены."
  exit 1
fi

sshpass -e ssh -L "$PROXY_PORT_SERVER:localhost:$PROXY_PORT_APP" -v "$SSH_USER@$SSH_HOST"


