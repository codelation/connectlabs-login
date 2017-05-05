#!/bin/sh

mac=`echo "$QUERY_STRING" | grep -oE "(^|[?&])mac=[^&]+" | cut -f 2 -d "="`
session=`echo "$QUERY_STRING" | grep -oE "(^|[?&])session=[^&]+" | cut -f 2 -d "="`
email=`echo "$QUERY_STRING" | grep -oE "(^|[?&])email=[^&]+" | cut -f 2 -d "="`
ip=`echo "$QUERY_STRING" | grep -oE "(^|[?&])ip=[^&]+" | cut -f 2 -d "="`
response=`udshape -n ssid1 -p login -m ${mac} -U ${email} -P ${session}`
response2=`udsplash -n ssid1 -a ${ip}`

echo "Content-type: text/html"
echo

cat << EOF
<html>
  <head>
   <title></title>
  </head>
  <body>
session = ${session}
<br />
email = ${email}
<br />
response = ${response}
<br />
response2 = ${response2}
<br />
  </body>
</html>
EOF
