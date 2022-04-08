# http-revproxy-injector
Reverse proxy to inject http component (header, form, etc) on request or response

# Examples

## add header cookie

```
# curl -X POST http://revproxydvwa:4322/revpr0xyconfig -H 'Content-Type: application/json' -d '{"name":"Cookie","value":"PHPSESSID=jv2db8n2jvjbjs4t44me934570; security=low", "place": "header"}'
{"message":"Data is injected"}

# curl http://revproxydvwa:4322/revpr0xyconfig
{"message":"OK","data":{"Cookie":{"value":"PHPSESSID=jv2db8n2jvjbjs4t44me934570; security=low","place":"header"}}}
```