#!/bin/bash
if [ -f app.pid ]; then
	kill -9 `cat app.pid`
	rm app.pid
fi
nohup ./app > app.log 2>&1 &
echo $! > app.pid

