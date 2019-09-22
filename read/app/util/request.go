package util

import (
	"crypto/tls"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"read/app/util/jaeger_service"
)

func HttpGet(url string) (string, error) {

	span, _ := opentracing.StartSpanFromContext(
		jaeger_service.ParentContext,
		"call Http Get",
		opentracing.Tag{Key: string(ext.Component), Value: "HTTP"},
		ext.SpanKindRPCClient,
	)
	span.Finish()

	tr := &http.Transport{
		TLSClientConfig : &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout   : time.Second * 5, //默认5秒超时时间
		Transport : tr,
	}

	req, err := http.NewRequest("GET", url,nil)
	if err != nil {
		return "", err
	}

	injectErr := jaeger_service.Tracer.Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
	if injectErr != nil {
		log.Fatalf("%s: Couldn't inject headers", err)
	}

	resp ,err :=  client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}
	return string(content), err
}