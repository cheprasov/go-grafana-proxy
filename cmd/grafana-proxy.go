package main

import (
    "flag"
    "fmt"
    "go-grafana-proxy/pkg/config"
    "go-grafana-proxy/pkg/grafana"
    "log"
    "strconv"
    "time"

    "github.com/valyala/fasthttp"
)

var cfg config.Config

func init() {
    configFilePointer := flag.String("config", "", "Path to config.json file")
    flag.Parse()

    var err error
    cfg, err = config.ReadConfig(*configFilePointer)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(cfg)
}

func main() {
    flag.Parse()

    h := requestHandler
    if cfg.Compress {
        h = fasthttp.CompressHandler(h)
    }

    if err := fasthttp.ListenAndServe(cfg.Listen, h); err != nil {
        log.Fatalf("Error in ListenAndServe: %s", err)
    }
}

func requestHandler(ctx *fasthttp.RequestCtx) {

    fmt.Println(ctx, &ctx.Request)

    switch string(ctx.RequestURI()) {
    case "/grafana-proxy":
        go routeGrafanaProxy(ctx)
    default:
        go routeDefault(ctx)
    }

    // If the file doesn't exist, create it, or append to the file
    //f, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    //if err != nil {
    //    log.Fatal(err)
    //}
    //defer f.Close()
    //if _, err := fmt.Fprint(f, &ctx.Request); err != nil {
    //    log.Fatal(err)
    //}
}

func routeDefault(ctx *fasthttp.RequestCtx) {
    //fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
    //fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
    //fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
    //fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
    //fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
    //fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
    //fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
    //fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
    //fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
    //fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

    fmt.Fprint(ctx, &ctx.Request)
    ctx.SetContentType("text/plain; charset=utf8")

    // Set arbitrary headers
    // ctx.Response.Header.Set("X-My-Header", "my-header-value")

    // Set cookies
    //var c fasthttp.Cookie
    //c.SetKey("cookie-name")
    //c.SetValue("cookie-value")
    //ctx.Response.Header.SetCookie(&c)
}

func routeError(ctx *fasthttp.RequestCtx) {
    ctx.Response.SetStatusCode(400)
}

func routeGrafanaProxy(ctx *fasthttp.RequestCtx) {
    if string(ctx.Method()) != "POST" {
        routeError(ctx)
        return
    }
    authKey := string(ctx.Request.Header.Peek("Authorization"))
    if string(authKey) != cfg.Authorization {
        routeError(ctx)
        return
    }

    name := string(ctx.PostArgs().Peek("name"))

    if len(name) == 0 {
        routeError(ctx)
        return
    }

    fmt.Fprintf(ctx, "Request method is %q\n")

    var metrics = make([]grafana.Metric, 0, 2)

    temperature := string(ctx.PostArgs().Peek("temperature"))
    if len(temperature) != 0 {
        value, err := strconv.ParseFloat(temperature, 64)
        if err == nil {
            metrics = append(metrics, grafana.Metric{
                Name:     "data.flat." + name + ".temperature",
                Value:    value,
                Time:     time.Now().Unix(),
                Interval: 60,
            })
        }
    }

    humidity := string(ctx.PostArgs().Peek("humidity"))
    if len(humidity) != 0 {
        value, err := strconv.ParseFloat(humidity, 64)
        if err == nil {
            metrics = append(metrics, grafana.Metric{
                Name:     "data.flat." + name + ".humidity",
                Value:    value,
                Time:     time.Now().Unix(),
                Interval: 60,
            })
        }
    }

    if len(metrics) == 0 {
        routeError(ctx)
        return
    }
    grafana.PostMetrics(cfg.PostUrl, cfg.ApiKey, metrics)
}
