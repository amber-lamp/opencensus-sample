package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

func main() {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "jaeger-agent:6831",
		Process: jaeger.Process{
			ServiceName: "demo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer exporter.Flush()

	trace.RegisterExporter(exporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	log.Println("Start demo server...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
	
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(context.Background(), "/handler")
	defer span.End()

	omikuji := omikuji(ctx)
	fmt.Fprintf(w, "運勢は[%s]です", omikuji)
}

func omikuji(ctx context.Context) string{
	_, childSpan := trace.StartSpan(ctx, "/omikuji")
	defer childSpan.End()
	t := time.Now()

	var omikuji string
	var msg string
	if (t.Month() == 1 && t.Day() >= 1 && t.Day() <= 3){
		omikuji = "大吉"
		msg = "お正月は大吉"
	} else {
		t := t.UnixNano()
		rand.Seed(t)
		s := rand.Intn(6)
		switch s + 1 {
		case 1:
			omikuji = "凶"
			msg = "残念でした"
			time.Sleep(time.Second)
		case 2, 3:
			omikuji = "吉"
			msg = "そこそこでした"
		case 4, 5:
			omikuji = "中吉"
			msg = "まあまあでした"
		case 6:
			omikuji = "大吉"
			msg = "いいですね"
		}
	}
	childSpan.Annotate([]trace.Attribute{
		trace.StringAttribute("omikuji", omikuji),
	}, msg)
	return omikuji
}