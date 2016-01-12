#!/bin/sh

killall mapd
killall gated
killall entitiesd

ps aux | egrep "(mapd|gated|entitiesd)" | grep -v "grep"
