package trace

import (
	instana "../go-sensor"
	"context"
	"fmt"
	opentrace "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

var GoSensor *instana.Sensor
var OTTracer opentrace.Tracer

/*func init() {

	// Instana
	GoSensor = instana.NewSensorWithOption(&GoSensorOpts)

	// OT tracing
	opentrace.InitGlobalTracer(instana.NewTracerWithOptions(&GoSensorOpts))

}*/

func SetupTracer(serviceName string) {


	var options = instana.Options{
		Service:  serviceName,
		LogLevel: instana.Debug,
	}

	// Instana
	GoSensor = instana.NewSensorWithOption(&options)

	// OT tracing
	opentrace.InitGlobalTracer(instana.NewTracerWithOptions(&options))
}

func GetTracerFromContext(ctx context.Context) opentrace.Tracer {
	parentSpan, ok := ctx.Value("parentSpan").(opentrace.Span)

	if ok == false {
		log.Panic("Unable to obtain parent span")
	} else {
		fmt.Printf("Retrieved tracer from context %T", parentSpan.Tracer())
	}
	return parentSpan.Tracer()
}

func GetTracersAndParentSpanFromContext(ctx context.Context) (opentrace.Tracer, opentrace.Span) {

	parentSpan, ok := ctx.Value("parentSpan").(opentrace.Span)

	if ok == false {
		log.Panic("Unable to obtain parent span")
	} else {
		fmt.Printf("Retrieved tracer from context %T", parentSpan.Tracer())
	}
	return parentSpan.Tracer(), parentSpan
}

func TraceHTTPClientGet(r *http.Request, url string, comments string) opentrace.Span {

	// Tracer from context
	tracer, parentSpan := GetTracersAndParentSpanFromContext(r.Context())
	childSpan := tracer.StartSpan("SQL", opentrace.ChildOf(parentSpan.Context()))
	childSpan.SetTag(string(ext.SpanKind), "client")
	childSpan.SetTag(string(ext.HTTPMethod), "GET")
	childSpan.SetTag(string(ext.HTTPUrl), url)
	childSpan.SetTag(string(ext.HTTPStatusCode), 200)
	childSpan.SetBaggageItem("Comments", comments)

	return childSpan

}

func TraceSQLExecution(r *http.Request, sql string, comments string) opentrace.Span {

	// Tracer from context
	tracer, parentSpan := GetTracersAndParentSpanFromContext(r.Context())

	childSpan := tracer.StartSpan("SQL", opentrace.ChildOf(parentSpan.Context()))
	childSpan.SetTag(string(ext.SpanKind), "client")
	childSpan.SetTag(string(ext.DBType), "MySQL")
	childSpan.SetTag(string(ext.DBInstance), "people")
	childSpan.SetTag(string(ext.DBUser), "pedro")
	childSpan.SetTag(string(ext.DBStatement), sql)
	childSpan.SetBaggageItem("Comments", comments)

	return childSpan

}

func TraceDBConnection(r *http.Request, comments string) opentrace.Span {

	// Tracer from context
	tracer, parentSpan := GetTracersAndParentSpanFromContext(r.Context())

	childSpan := tracer.StartSpan("Connection", opentrace.ChildOf(parentSpan.Context()))
	childSpan.SetTag(string(ext.SpanKind), "client")
	childSpan.SetTag(string(ext.DBType), "MySQL")
	childSpan.SetTag(string(ext.DBInstance), "people")
	childSpan.SetTag(string(ext.DBUser), "pedro")
	childSpan.SetTag(string(ext.DBStatement), "connect")
	childSpan.SetBaggageItem("Comments", comments)

	return childSpan
}

func TraceFunctionExecution(r *http.Request, f func(), comments string) {

	// Tracer from context
	tracer, parentSpan := GetTracersAndParentSpanFromContext(r.Context())

	childSpan := tracer.StartSpan("method", opentrace.ChildOf(parentSpan.Context()))
	childSpan.SetTag(string(ext.SpanKind), "intermediate")
	childSpan.SetTag(string(ext.Component), "method")
	childSpan.SetBaggageItem("Comments", comments)

	// Execute function
	f()

	childSpan.Finish()

}

func TraceError(r *http.Request, errorCode int, message string) {
	// Parent span
	parentSpan, ok := r.Context().Value("parentSpan").(opentrace.Span)

	if ok == false {
		OTTracer = nil
		log.Panic("Unable to obtain parent span")
	}

	parentSpan.SetTag(string(ext.Error), true)
	parentSpan.SetTag(string(ext.HTTPStatusCode), errorCode)
	parentSpan.SetTag("message", message)
}

func TracePanic(r *http.Request, message string) {
	// Parent span
	parentSpan, ok := r.Context().Value("parentSpan").(opentrace.Span)

	if ok == false {
		OTTracer = nil
		log.Panic("Unable to obtain parent span")
	}

	parentSpan.SetTag(string(ext.Error), true)
	parentSpan.SetTag("message", message)
	stackTrace := string(debug.Stack())

	stackTraceLines := strings.Split(stackTrace, "\n")
	parentSpan.SetBaggageItem("Stack Trace", "Error details")
	for i, line := range stackTraceLines {
		fmt.Println(line)
		parentSpan.SetBaggageItem(strconv.Itoa(i), line)
	}

	panic(message)
}
