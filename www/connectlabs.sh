#!/bin/sh

cl_clean_mac() {
  echo `echo "$1"|awk '{gsub("-",":",$0); print tolower($0)}'`
}

cl_get_ip_from_mac() {
  MAC=`cl_clean_mac $1`
  IP=`cat /var/dhcp.leases | grep -oE "^[0-9]* ${MAC} ([^ ]*)"| cut -f 3 -d " "`
  echo "$IP"
}

cl_send_status() {
  MAC=`cl_clean_mac $1`
  IP=`cl_get_ip_from_mac ${MAC}`
  echo `udsplash -n ssid1 -a ${IP}`
}

cl_send_login() {
  MAC=`cl_clean_mac $1`
  EMAIL=$2
  COMMAND="udshape -n ssid1 -p login -m ${MAC} -U ${EMAIL} -P ${MAC}"
  RESPONSE=`${COMMAND}`
  if [ $? == 0 ]; then
    echo "Success"
    RESPONSE=`cl_send_status ${MAC}`
  else
    echo "$COMMAND"
    echo "$RESPONSE"
  fi
}
