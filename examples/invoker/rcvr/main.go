package main

import (
	"context"
	"log"

	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/invoker/receiver"
)

type Hello struct{}

func (*Hello) SayHello(ctx context.Context, req string, rsp *string) error {
	traceId := ctx.Value("traceId")
	log.Println(traceId)
	*rsp = req + " dodo"
	return nil
}

func NewInterceptor() invoker.Interceptor {
	return func(fn invoker.InvokeFunc) invoker.InvokeFunc {
		return func(ctx context.Context, mi invoker.Message,
			opts ...invoker.InvokeOption) (invoker.Message, error) {
			log.Println("begin invoke")
			mo, err := fn(ctx, mi, opts...)
			log.Println("end invoke")
			return mo, err
		}
	}
}

func main() {
	inv := receiver.NewInvoker(new(Hello), invoker.Intercept(NewInterceptor()))
	if err := inv.Init(); err != nil {
		log.Fatalln(err)
	}
	log.Println(inv.Name())
	mi := invoker.NewMessage()
	mi.SetFuncName("SayHello")
	mi.SetAttachment("traceId", "sdfasrpq92371-1274sfjxnm-xdjsalfo")
	var rsp string
	mi.SetParameters([]interface{}{context.Background(), "hello", &rsp})
	mo, err := inv.Invoke(context.Background(), mi)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(rsp)
	log.Println(mo)

}
