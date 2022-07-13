package https

import (
	"bytes"
	"context"
	"encoding/json"
	"goflow/commons"

	"errors"
	"fmt"
	"net/url"
	"reflect"
)

type JSONMarshallerUnmarshal interface {
	json.Marshaler
	json.Unmarshaler
}

func newHTTPEndpoint() *HTTPEndpoint {
	return &HTTPEndpoint{
		unformattedURL: "",
		url:            nil,
		queryParam:     make(url.Values),
		body:           nil,
		responseModel:  nil,
	}
}

type HTTPEndpoint struct {
	unformattedURL string
	url            *url.URL
	queryParam     url.Values
	body           *bytes.Buffer
	responseModel  JSONMarshallerUnmarshal
}

func (httpEndpoint *HTTPEndpoint) MyType() string {
	return reflect.TypeOf(httpEndpoint).String()
}

func (httpEndpoint *HTTPEndpoint) MyName() string {
	return ""
}

func (httpEndpoint *HTTPEndpoint) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["name"] = httpEndpoint.MyName()
	descriptionMap["unformatted url"] = httpEndpoint.unformattedURL
	descriptionMap["url"] = httpEndpoint.url.String()
	descriptionMap["query param"] = httpEndpoint.queryParam
	if httpEndpoint.body != nil {
		descriptionMap["body"] = httpEndpoint.body.String()
	}
	descriptionMap["response model"] = httpEndpoint.responseModel
	return descriptionMap
}

func (httpEndpoint *HTTPEndpoint) GetUnformattedURL() string {
	return httpEndpoint.unformattedURL
}

func (httpEndpoint *HTTPEndpoint) GetURL() *url.URL {
	return httpEndpoint.url
}

func (httpEndpoint *HTTPEndpoint) GetQueryParam() url.Values {
	return httpEndpoint.queryParam
}

func (httpEndpoint *HTTPEndpoint) GetBody() *bytes.Buffer {
	return httpEndpoint.body
}

func (httpEndpoint *HTTPEndpoint) GetResponseModel() JSONMarshallerUnmarshal {
	return httpEndpoint.responseModel
}

func NewHTTPEnpointBuilder(ctx context.Context) *httpEndpointBuilder {
	return &httpEndpointBuilder{
		nil,
		make([]string, 0, 3),
		make(map[string][]string),
		nil,
	}
}

type httpEndpointBuilder struct {
	body                  interface{}
	urlParams             []string
	queryParams           map[string][]string
	endPointResponseModel JSONMarshallerUnmarshal
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) MyType() string {
	return reflect.TypeOf(genericHTTPEndpointBuilder).String()
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) MyName() string {
	return ""
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) Describe() map[string]interface{} {
	descriptionMap := make(map[string]interface{})
	descriptionMap["name"] = genericHTTPEndpointBuilder.MyName()
	descriptionMap["url params"] = genericHTTPEndpointBuilder.urlParams
	descriptionMap["query param"] = genericHTTPEndpointBuilder.queryParams
	descriptionMap["response model id"] = genericHTTPEndpointBuilder.endPointResponseModel
	return descriptionMap
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) generateURL() (int, string, string, error) {
	return 0, "", "", nil
}

// NOTE: fill urlParams array keeping in mind the placeholders in url string
func (genericHTTPEndpointBuilder *httpEndpointBuilder) AddURLParam(value string) error {
	if genericHTTPEndpointBuilder.urlParams != nil {
		genericHTTPEndpointBuilder.urlParams = append(genericHTTPEndpointBuilder.urlParams, value)
	} else {
		return errors.New("nil url parameter slice for http endpoint builder")
	}
	return nil
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) AddQueryParam(key string, value string) error {
	if genericHTTPEndpointBuilder.queryParams != nil {
		queryArray, ok := genericHTTPEndpointBuilder.queryParams[key]
		if ok {
			queryArray = append(queryArray, value)
		} else {
			genericHTTPEndpointBuilder.queryParams[key] = make([]string, 0)
			genericHTTPEndpointBuilder.queryParams[key] = append(genericHTTPEndpointBuilder.queryParams[key], value)
		}
	} else {
		return errors.New("nil query parameter slice for http endpoint builder")
	}

	return nil
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) SetEndpointResponseModelType(endpointmodel JSONMarshallerUnmarshal) {
	genericHTTPEndpointBuilder.endPointResponseModel = endpointmodel
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) SetBody(body interface{}) {
	genericHTTPEndpointBuilder.body = body
}

func (genericHTTPEndpointBuilder *httpEndpointBuilder) Build() (*HTTPEndpoint, error) {
	httpEndpoint := newHTTPEndpoint()
	if genericHTTPEndpointBuilder.body != nil {
		marshalledBody, err := json.Marshal(genericHTTPEndpointBuilder.body)
		if err != nil {
			return nil, commons.WrapError(err, "unable to marshal body")
		}
		httpEndpoint.body = bytes.NewBuffer(marshalledBody)
	}

	for key, valueArr := range genericHTTPEndpointBuilder.queryParams {
		for _, value := range valueArr {
			if key != "" && value != "" {
				httpEndpoint.queryParam.Add(key, value)
			} else {
				return nil, fmt.Errorf("incorrect query param (key = %s | value = %s)", key, value)
			}
		}
	}

	httpEndpoint.responseModel = genericHTTPEndpointBuilder.endPointResponseModel
	noOfPleaceholder, urlStr, unformattedURL, err := genericHTTPEndpointBuilder.generateURL()
	if err != nil {
		return nil, commons.WrapError(err, "url generation failed")
	}

	numOfURLParams := len(genericHTTPEndpointBuilder.urlParams)
	if noOfPleaceholder != numOfURLParams {
		return nil, fmt.Errorf("number of url params(%d) are not equal to number of placeholders(%d)", noOfPleaceholder, numOfURLParams)
	}

	httpEndpoint.unformattedURL = unformattedURL
	for _, param := range genericHTTPEndpointBuilder.urlParams {
		urlStr = fmt.Sprintf(urlStr, param)
	}

	formattedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, commons.WrapError(err, fmt.Sprintf("failed to parse url string(%s)", urlStr))
	}

	httpEndpoint.url = formattedURL

	return httpEndpoint, nil
}
