#!/bin/sh

./mapd &
./gated -gid 1 &
./entitiesd -pid 1 &
ps aux | egrep "(mapd|gated|entitiesd)" | grep -v "grep"
