#!/bin/bash
pid=$(ps aux | grep -i urlencoder | awk {'print $2'} )
if ps -p $pid >/dev/null
then
   echo "$pid is running"
   kill -9 $pid
fi
ulimit -n 614400;
nohup ./urlencoder.linux.x386 > urlencoder.log &