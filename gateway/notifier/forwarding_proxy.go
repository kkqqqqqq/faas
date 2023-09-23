// Copyright (c) OpenFaaS Author(s). All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package notifier

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"sync"


	"github.com/openfaas/faas/gateway/pkg/middleware"
	"github.com/openfaas/faas/gateway/types"
	"github.com/openfaas/faas/gateway/essemble"
)


// MakeForwardingProxyHandler create a handler which forwards HTTP requests
func MakeForwardingProxyHandler(proxy *types.HTTPClientReverseProxy,
	notifiers []HTTPNotifier,
	baseURLResolver middleware.BaseURLResolver,
	urlPathTransformer middleware.URLPathTransformer,
	serviceAuthInjector middleware.AuthInjector) http.HandlerFunc {

	writeRequestURI := false
	if _, exists := os.LookupEnv("write_request_uri"); exists {
		writeRequestURI = exists
	}

	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := baseURLResolver.Resolve(r)	
		originalURL := r.URL.String()
		requestURL := urlPathTransformer.Transform(r)

		for _, notifier := range notifiers {
			notifier.Notify(r.Method, requestURL, originalURL, http.StatusProcessing, "started", time.Second*0)
		}

		start := time.Now()
		//猜测这的时返回给用户的，里面的upstreamrequest是给函数的
		_,statusCode, err := forwardRequest(w, r, proxy.Client, baseURL, requestURL, proxy.Timeout, writeRequestURI, serviceAuthInjector)
		seconds := time.Since(start)

		if err != nil {
			log.Printf("error with upstream request to: %s, %s\n", requestURL, err.Error())
		}


		for _, notifier := range notifiers {
			notifier.Notify(r.Method, requestURL, originalURL, statusCode, "completed", seconds)
		}

		}

}

func MakeEssembleForwardingProxyHandler(proxy *types.HTTPClientReverseProxy,
	notifiers []HTTPNotifier,
	baseURLResolver middleware.BaseURLResolver,
	urlPathTransformer middleware.URLPathTransformer,
	serviceAuthInjector middleware.AuthInjector) http.HandlerFunc {

	writeRequestURI := false
	if _, exists := os.LookupEnv("write_request_uri"); exists {
		writeRequestURI = exists
	}

	
	return func(w http.ResponseWriter, r *http.Request) {
		baseURL := baseURLResolver.Resolve(r)	 //     http://127.0.0.1:8081
		originalURL := r.URL.String()	//   /classification
		requestURL := urlPathTransformer.Transform(r)  //     /classification

		for _, notifier := range notifiers {
			notifier.Notify(r.Method, requestURL, originalURL, http.StatusProcessing, "started", time.Second*0)
		}

		log.Printf("classificaion:from gateway:--baseURL:%s --requestURL:%s --originalURL:%s\n",baseURL, requestURL, originalURL)
		
		model_selected := essemble.ModelSelection(1 ,1,"efaas") 
		
		time_start := time.Now()
		
		_,statusCode, err:= EssembleForwardRequest(model_selected,w, r, proxy.Client, baseURL, proxy.Timeout, writeRequestURI, serviceAuthInjector)
		
		if err != nil {
			log.Printf("error with upstream request to: %s, %s\n", requestURL, err.Error())
		}

		time_cost := time.Since(time_start)
		log.Printf("all time cost: %.4f seconds\n",time_cost/1000000000)

		
		for _, notifier := range notifiers {
			notifier.Notify(r.Method,requestURL, originalURL, statusCode, "completed", time_cost)
		}

	}
}


func buildUpstreamRequest(r *http.Request, baseURL string, requestURL string) *http.Request {
	url := baseURL + requestURL
	
	if len(r.URL.RawQuery) > 0 {
		url = fmt.Sprintf("%s?%s", url, r.URL.RawQuery)
	}
	//log.Printf("buildupstreamRequest:"+url)
	//NewRequest 使用给定的Method，URL 和可选的BODY参数，返回一个新的Request，
	upstreamReq, _ := http.NewRequest(r.Method, url, nil)
	//设置给函数的请求的头
	copyHeaders(upstreamReq.Header, &r.Header)
	// Hop-by-hop headers. These are removed when sent to the backend.
	deleteHeaders(&upstreamReq.Header, &hopHeaders)

	if len(r.Host) > 0 && upstreamReq.Header.Get("X-Forwarded-Host") == "" {
		upstreamReq.Header["X-Forwarded-Host"] = []string{r.Host}
	}

	if upstreamReq.Header.Get("X-Forwarded-For") == "" {
		upstreamReq.Header["X-Forwarded-For"] = []string{r.RemoteAddr}
	}

	if r.Body != nil {
		upstreamReq.Body = r.Body
	}
	return upstreamReq
}

func forwardRequest (w http.ResponseWriter,
	r *http.Request,
	proxyClient *http.Client,
	baseURL string,
	requestURL string,
	timeout time.Duration,
	writeRequestURI bool,
	serviceAuthInjector middleware.AuthInjector) ([]byte,int, error) {

	//	log.Printf("line 209")
	// 创建了一个请求？
	upstreamReq := buildUpstreamRequest(r, baseURL, requestURL)
	//defer关键字修饰的语句会推迟到执行return语句或函数执行完毕以及发生错误之后
	//才会执行。defer语句常用于成对的操作，例如打开和关闭，连接和断开，加锁和解锁。
	if upstreamReq.Body != nil {
		defer upstreamReq.Body.Close()
	}
	//log.Printf("line 214")

	if serviceAuthInjector != nil {
		serviceAuthInjector.Inject(upstreamReq)
		
	}

	if writeRequestURI {
		log.Printf("writeRequestURI")
		log.Printf("forwardRequest: %s %s\n", upstreamReq.Host, upstreamReq.URL.String())
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	defer cancel()

	//客户端，发送请求，收到回复
	res, resErr := proxyClient.Do(upstreamReq.WithContext(ctx))
	//error
	if resErr != nil {
		badStatus := http.StatusBadGateway
		w.WriteHeader(badStatus)
		badbody := []byte("wrong") 
		return badbody,badStatus,resErr
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	copyHeaders(w.Header(), &res.Header)

	
	w.WriteHeader(res.StatusCode)
	
	if res.Body != nil {
		// Copy the body over
		io.CopyBuffer(w, res.Body, nil)
	}
	
	body,_:= ioutil.ReadAll(res.Body)
	return body,res.StatusCode, nil
}




func EssembleForwardRequest (model_selected []string,
	w http.ResponseWriter,
	r *http.Request,
	proxyClient *http.Client,
	baseURL string,
	timeout time.Duration,
	writeRequestURI bool,
	serviceAuthInjector middleware.AuthInjector) ([]byte,int, error) {
	
	n := len(model_selected)
	

	log.Printf(" EssembleForwardRequest:--model_selected：%s" ,model_selected)
	log.Printf(" EssembleForwardRequest:--number of selected models：%d" ,n )
	log.Printf(" EssembleForwardRequest:--baseURL:%s" ,baseURL)

	//用go协程创建多个upstreamRequest
	wg := &sync.WaitGroup{}

	models := make(chan string)
	results := make(chan string)

	go func() {
		for result := range results{
			if result == ""{
				close(results)
			}
			// else{
			// 	fmt.Println("result:", result)
			// }
		}
	}()

	for i:= 0; i <= n; i++{
		wg.Add(1)
		go worker(timeout,serviceAuthInjector,proxyClient,r,writeRequestURI, baseURL,wg, models, results)
	}

	for _, model := range model_selected{
		models <- model
	}

	models <- "" //worker结束标志
	wg.Wait() //同步结束
	results <- "" // result结束标志
	log.Printf("worker done")


	StatusCode := http.StatusOK
	body :="yes"

	return []byte(body),StatusCode, nil

}


func worker(
	timeout time.Duration,
	serviceAuthInjector middleware.AuthInjector,
	proxyClient *http.Client,
	r *http.Request,
	writeRequestURI bool,
	baseURL string,
	group *sync.WaitGroup, 
	models chan string, 
	result chan string) {
	for model := range models{
		log.Printf("this is "+baseURL)

		if model == ""{
			close(models)
		}else{
			model= "/function/"+model;
			log.Printf(model)
			upstreamReq := buildUpstreamRequest(r, baseURL, model)
			if upstreamReq.Body != nil {
				defer upstreamReq.Body.Close()
			}
			if serviceAuthInjector != nil {
				serviceAuthInjector.Inject(upstreamReq)
				
			}
		
			if writeRequestURI {
				log.Printf("corker:writeRequestURI")
				log.Printf("forwardRequest: %s %s\n", upstreamReq.Host, upstreamReq.URL.String())
			}
		
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			res, resErr := proxyClient.Do(upstreamReq.WithContext(ctx))

			if resErr !=nil{
				result<-"wrong"
			}else{
				if res.Body != nil {
					defer res.Body.Close()
				}
				body,_:= ioutil.ReadAll(res.Body)
				result<-string(body)
			}
			
		}
	}
	group.Done()

}


func copyHeaders(destination http.Header, source *http.Header) {
	for k, v := range *source {
		vClone := make([]string, len(v))
		copy(vClone, v)
		(destination)[k] = vClone
	}
}


func deleteHeaders(target *http.Header, exclude *[]string) {
	for _, h := range *exclude {
		target.Del(h)
	}
}

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.
// Copied from: https://golang.org/src/net/http/httputil/reverseproxy.go
var hopHeaders = []string{
	"Connection",
	"Proxy-Connection", // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",      // canonicalized version of "TE"
	"Trailer", // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",
	"Upgrade",
}
