# http-revproxy-injector
Reverse proxy to inject http component (header, form, etc) on request or response

# To use config

```
curl -X POST -d "name=Cookie&value=12345&place=test" http://192.168.1.220:4321/revpr0xyconfig

curl -X DELETE -d "name=Cookie" http://192.168.1.220:4321/revpr0xyconfig

curl http://192.168.1.220:4321/revpr0xyconfig
```