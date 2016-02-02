#!/bin/sh
#
# /etc/init.d/mysystem
# Subsystem file for "mirror" server
#
# chkconfig: 2345 95 05	(1)
# description: mirror server daemon
#
# processname:mirror
# config: /etc/mirror/mirror.conf
# config: /etc/sysconfigmirror/
# pidfile: /var/run/mirror.pid

NAME="mirror"
: ${MIRROR_LOGFILE:=/var/log/mirror.log}
: ${DAEMONOPTS:="--insecure"}
PIDFILE=/var/run/mirror.pid

# source function library
[ -f /etc/rc.d/init.d/functions ] && . /etc/rc.d/init.d/functions

# pull in sysconfig settings
[ -f /etc/sysconfig/mirror ] && . /etc/sysconfig/mirror

RETVAL=0

start() {
	echo -n $"Starting $NAME:"
	mirror daemon ${DAEMONOPTS} >> "${MIRROR_LOGFILE}" 2>&1 &
	PID=$(ps -ef | grep mirror | grep daemon | awk -F" " '{print $2}')
  if [ -z $PID ]; then
      printf "%s\n" "Fail"
  else
      echo $PID > $PIDFILE
      printf "%s\n" "Ok"
  fi
}

reload() {
	echo -n "Reloading $NAME:"
	restart
}

stop() {
	kill $(cat $PIDFILE)
}

restart() {
	if check; then
	    stop && sleep 1 && start
	else
	    start
	fi
}

check() {
	[ -f $PIDFILE ] && ps -A -o pid | grep "^\s*$(cat $PIDFILE)$" > /dev/null 2>&1
}

status() {
    if check; then
        echo 'Mirror daemon is running'
        exit 0
    else
        echo 'Mirror daemon is not running'
        exit 1
    fi
}

case "$1" in
	start)
		start
		;;
	stop)
		stop
		;;
	restart)
		stop
		start
		;;
	reload)
		reload
		;;
	condrestart)
		if [ -f /var/lock/subsys/$NAME ] ; then
			stop
			# avoid race
			sleep 3
			start
		fi
		;;
	status)
		status
		;;
	*)
		echo "Usage: $0 {start|stop|restart|reload|condrestart|status}"
		RETVAL=1
esac
exit $RETVAL
