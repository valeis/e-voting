package web

import (
	"bytes"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"net/url"
	"rest-api-go/internal/repository"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
	UserRepo     repository.UserRepository
	RedisClient  any
}

type HandlerFactory interface {
	CreateHandler(action string) http.HandlerFunc
}

type OrgHandlerFactory struct {
	setups OrgSetup
}

func (o OrgHandlerFactory) CreateHandler(action string) http.HandlerFunc {
	switch action {
	case "query":
		return o.setups.QueryHandler
	case "invoke":
		return o.setups.Invoke
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid action", http.StatusNotFound)
		}
	}
}

func Serve(setups OrgSetup, r *fiber.App) {
	factory := OrgHandlerFactory{setups: setups}

	r.Get("/query", func(c *fiber.Ctx) error {
		handler := factory.CreateHandler("query")
		adaptHandlerFuncToFiber(handler, c)
		return nil
	})

	r.Post("/invoke", func(c *fiber.Ctx) error {
		handler := factory.CreateHandler("invoke")
		adaptHandlerFuncToFiber(handler, c)
		return nil
	})

	//http.HandleFunc("/query", factory.CreateHandler("query"))
	//http.HandleFunc("/invoke", factory.CreateHandler("invoke"))
	fmt.Println("Listening (http://localhost:3000/)...")
	//if err := http.ListenAndServe(":3000", nil); err != nil {
	//	fmt.Println(err)
	//}
}

func adaptHandlerFuncToFiber(handler http.HandlerFunc, c *fiber.Ctx) {
	// Create a fake ResponseWriter for fiber
	resWriter := &responseWriterAdapter{ctx: c}

	//Call the standard library handler with adapted response writer and Fiber request
	req := toHTTPRequest(c)
	handler(resWriter, req)

	//Ensure the response is flushed
	c.Send(resWriter.body.Bytes())
}

// Adapter for ResponseWrite
type responseWriterAdapter struct {
	ctx  *fiber.Ctx
	body *bytes.Buffer
}

func (w *responseWriterAdapter) Header() http.Header {
	return http.Header{}
}

func (w *responseWriterAdapter) Write(data []byte) (int, error) {
	if w.body == nil {
		w.body = &bytes.Buffer{}
	}
	return w.body.Write(data)
}

func (w *responseWriterAdapter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
}

// Convert Fiber's request to an http.Request

func toHTTPRequest(c *fiber.Ctx) *http.Request {
	req := new(http.Request)
	req.Method = string(c.Request().Header.Method())
	req.URL = &url.URL{
		Path:     string(c.Request().URI().Path()),
		RawQuery: string(c.Request().URI().QueryString()),
	}
	req.Header = http.Header{}
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Add(string(key), string(value))
	})
	req.Body = io.NopCloser(bytes.NewReader(c.Body()))
	return req
}
