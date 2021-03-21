package grafana

import (
    "encoding/json"
    "fmt"
    "github.com/valyala/fasthttp"
)

type Metric struct {
    Name  string  `json:"name"`
    Value float64 `json:"value"`
    Time  int64   `json:"time"`
    Interval int  `json:"interval"`
}

func PostMetrics(url string, apikey string, metrics []Metric) {
    postData, _ := json.Marshal(metrics)
    fmt.Println(string(postData))

    // Put these in the global scope so they don't get converted for each request.

    req := fasthttp.AcquireRequest()
    req.Header.SetMethod("POST")
    req.Header.SetContentType("application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
    req.SetRequestURI(url)
    req.SetBody(postData)

    res := fasthttp.AcquireResponse()
    if err := fasthttp.Do(req, res); err != nil {
        panic("handle error")
    }
    fasthttp.ReleaseRequest(req)

    body := res.Body()

    fmt.Println(string(body))
    // Do something with body.

    fasthttp.ReleaseResponse(res) // Only when you are done with body!
}
