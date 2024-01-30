MY_URL=http://localhost:38080

curl -I -X OPTIONS \
  -H "Origin: ${MY_URL}" \
  -H 'Access-Control-Request-Method: GET' \
  -H 'Cache-Control: no-cache' \
  "${MY_URL}/api/orders"
