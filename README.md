This package provides three types of logging for api services. (From an old project)
- Keyval logger
- JSON logger
- Nil logger

Check out the example folder on how to use.

```go
func main() {
  svcLogger := log.NewJSONLogger()
  if os.Getenv("ENV") == "DEV" {
    svcLogger = log.NewKeyvalLogger(log.ColorYellow)
  }
  
  // Replace the global default logger with a configured logger.
  svcLogger = svcLogger.With("service", "log-example", "version", 1.0)
  log.SetLogger(svcLogger)
  
  http.HandleFunc("/", withLogger(svcLogger, stubHandler))
  http.HandleFunc("/admin", withLogger(svcLogger, protectedRoute(stubHandler)))
  
  // Use the global logger to log.
  log.With("port", "8080").Print("listening")
  _ = http.ListenAndServe(":8080", nil)
}
```

Output (JSONLogger)
```json
{"time":"2021-04-19T01:49:57Z","message":"listening","port":"8080","service":"log-example","version":1}
{"time":"2021-04-19T01:49:59Z","message":"Request.","requestID":"78fad37cd3dc62d1466af9e15c3e9835","requestURI":"/","method":"GET","service":"log-example","version":1}
{"time":"2021-04-19T01:49:59Z","message":"Response.","code":200,"duration":"0ms","requestID":"78fad37cd3dc62d1466af9e15c3e9835","requestURI":"/","method":"GET","service":"log-example","version":1}
```

Output (KeyvalLogger)
```shell
time="2021-04-19T01:50:22Z" message="listening" port="8080" service="log-example" version=1
time="2021-04-19T01:50:27Z" message="Request." requestID="78fad37cd3dc62d1466af9e15c3e9835" requestURI="/" method="GET" service="log-example" version=1
time="2021-04-19T01:50:27Z" message="Response." code=200 duration="0ms" requestID="78fad37cd3dc62d1466af9e15c3e9835" requestURI="/" method="GET" service="log-example" version=1
```
