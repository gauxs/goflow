package https

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goflow/commons"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"time"
)

const (
	HTTPAuthorizationHeaderKey        = "Authorization"
	HTTPXAuthorizationHeaderKey       = "X-Authorization"
	HTTPEncodingFormatHeaderKey       = "Content-Encoding"
	HTTPContentTypeHeaderKey          = "Content-Type"
	HTTPContentTypeHeaderDefaultValue = "application/json"
)

type HttpContentEncodingFormat uint8

const (
	UndefinedContentEncoding HttpContentEncodingFormat = iota
	Gzip
	Compress
	Deflate
	Br
)

func (httpContentEncodingFormat HttpContentEncodingFormat) String() string {
	switch httpContentEncodingFormat {
	case UndefinedContentEncoding:
		return "undefined"
	case Gzip:
		return "gzip"
	case Compress:
		return "compress"
	case Deflate:
		return "identity"
	case Br:
		return "br"
	}

	return ""
}

type HttpMethod uint8

const (
	UndefinedHTTPMethod HttpMethod = iota
	Get
	Put
	Post
)

func (httpMethod HttpMethod) String() string {
	switch httpMethod {
	case UndefinedHTTPMethod:
		return "undefined"
	case Get:
		return "Get"
	case Put:
		return "Put"
	case Post:
		return "Post"
	}

	return ""
}

type JSONMarshallerUnmarshaller interface {
	json.Marshaler
	json.Unmarshaler
}

func NewHTTPHandler() (*HTTPHandler, error) {
	return &HTTPHandler{
		UndefinedHTTPMethod,
		make(http.Header),
		nil,
		nil,
	}, nil
}

type HTTPHandler struct {
	method            HttpMethod
	header            http.Header
	httpEndpoint      *HTTPEndpoint
	constructionLogic commons.ActorConstructionLogic
}

func (httpHandler *HTTPHandler) MyType() string {
	return reflect.TypeOf(httpHandler).String()
}

func (httpHandler *HTTPHandler) MyName() string {
	return ""
}

func (httpHandler *HTTPHandler) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["name"] = httpHandler.MyName()
	descriptionMap["http method"] = httpHandler.method.String()
	descriptionMap["header details"] = httpHandler.header
	if httpHandler.httpEndpoint != nil {
		descriptionMap["endpoint details"] = httpHandler.httpEndpoint.Describe()
	}
	descriptionMap["actor knows its construction"] = fmt.Sprintf("%t", (httpHandler.constructionLogic != nil))
	return descriptionMap
}

func (httpHandler *HTTPHandler) GetHTTPMethod() HttpMethod {
	return httpHandler.method
}

func (httpHandler *HTTPHandler) SetHTTPMethod(httpmethod HttpMethod) {
	httpHandler.method = httpmethod
}

func (httpHandler *HTTPHandler) GetHeader() http.Header {
	return httpHandler.header
}

func (httpHandler *HTTPHandler) SetHeader(header http.Header) {
	authValue := header.Get(HTTPAuthorizationHeaderKey)
	xAuthValue := header.Get(HTTPXAuthorizationHeaderKey)
	contentTypeValue := header.Get(HTTPContentTypeHeaderKey)
	if contentTypeValue == "" {
		contentTypeValue = HTTPContentTypeHeaderDefaultValue
	}

	httpHandler.AddToHeader(HTTPAuthorizationHeaderKey, authValue)
	httpHandler.AddToHeader(HTTPXAuthorizationHeaderKey, xAuthValue)
	httpHandler.AddToHeader(HTTPContentTypeHeaderKey, contentTypeValue)
}

func (httpHandler *HTTPHandler) AddToHeader(key string, value string) error {
	if httpHandler.header != nil {
		httpHandler.header.Add(key, value)
	} else {
		return errors.New("http request struct has nil header")
	}

	return nil
}

func (httpHandler *HTTPHandler) GetEndpoint() *HTTPEndpoint {
	return httpHandler.httpEndpoint
}

func (httpHandler *HTTPHandler) SetEndpoint(httpEndpoint *HTTPEndpoint) {
	httpHandler.httpEndpoint = httpEndpoint
}

func (httpHandler *HTTPHandler) SetConstructionLogic(logic commons.ActorConstructionLogic) {
	httpHandler.constructionLogic = logic
}

func (httpHandler *HTTPHandler) MakeRequest(ctx context.Context) error {
	var err error
	var httpURL *url.URL
	var request *http.Request
	var httpMethod HttpMethod
	var httpHeader http.Header
	var responseBodyByte []byte
	var httpEndpoint *HTTPEndpoint

	if httpHandler.method == UndefinedHTTPMethod {
		return errors.New("undefined http method when trying to make a http request")
	}

	httpMethod = httpHandler.GetHTTPMethod()
	httpEndpoint = httpHandler.GetEndpoint()

	if httpEndpoint == nil {
		return errors.New("nil http endpoint when trying to make a http request")
	}

	httpURL = httpEndpoint.GetURL()

	if httpURL == nil {
		return errors.New("nil url when trying to make a http request")
	}

	body := httpEndpoint.GetBody()

	if body != nil {
		request, err = http.NewRequest(httpMethod.String(), httpURL.String(), body)
	} else {
		request, err = http.NewRequest(httpMethod.String(), httpURL.String(), nil)
	}
	if err != nil {
		return commons.WrapError(err, "unable to Get a new http request using http package")
	}

	queryParam := httpEndpoint.GetQueryParam()
	if queryParam != nil {
		request.URL.RawQuery = queryParam.Encode()
	}

	httpHeader = httpHandler.GetHeader()
	request.Header = httpHeader

	// just for representation,
	// can be implemented better,
	// by using resource pooling
	var httpClient *http.Client
	httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5,
		},
		Timeout: 120 * time.Second,
	}

	// externalSegment := newrelic.StartExternalSegment(newRelicTransaction, request)
	response, err := httpClient.Do(request)
	// externalSegment.Response = response
	// externalSegment.End()

	if err != nil {
		return commons.WrapError(err, "unable to Do a http reuqest")
	}
	defer response.Body.Close()

	if response.StatusCode >= 400 {
		// handle 4xx - 5xx errors
	}

	responseEncodingFormat := response.Header.Get(HTTPEncodingFormatHeaderKey)
	switch responseEncodingFormat {
	case Gzip.String():
		var responseReader io.ReadCloser
		if responseReader, err = gzip.NewReader(response.Body); err != nil {
			return commons.WrapError(err, fmt.Sprintf("%s encoding format, unable to prepare response body", responseEncodingFormat))
		}
		defer responseReader.Close()
		responseBodyByte, err = ioutil.ReadAll(responseReader)
	default:
		responseBodyByte, err = ioutil.ReadAll(response.Body)
	}
	if err != nil {
		return commons.WrapError(err, fmt.Sprintf("%s encoding format, unable to read response body", responseEncodingFormat))
	}

	respModel := httpEndpoint.GetResponseModel()
	if respModel == nil {
		return errors.New("nil endpoint response model when going to unmrshal response into it")
	}

	err = respModel.UnmarshalJSON(responseBodyByte)
	if err != nil {
		return commons.WrapError(err, "unable to unmarshal response into model")
	}

	return nil
}

func (httpHandler *HTTPHandler) Act(ctx context.Context, sourceDataStore commons.OperableDataStore, destinationDataStore commons.OperableDataStore) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// log runtime errors here
			err = errors.New(string(debug.Stack()))
		}
	}()

	err = httpHandler.MakeRequest(ctx)
	if err != nil {
		// log http request errors here
		return err
	}

	// will never be nil as we have already checked
	endpoint := httpHandler.GetEndpoint()
	// will never be nil as we have already checked
	endpointModel := endpoint.GetResponseModel()

	// Transfer the required data to a datastore
	destinationDataStore.Push(endpointModel)
	return nil
}

func (httpHandler *HTTPHandler) Construct(ctx context.Context, sourceDataStore commons.OperableDataStore) (err error) {
	if httpHandler.constructionLogic == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			// log runtime errors here
			err = errors.New(string(debug.Stack()))
		}
	}()

	err = httpHandler.constructionLogic(ctx, httpHandler, sourceDataStore)
	if err != nil {
		// log contruction errors here
		return err
	}

	return nil
}
