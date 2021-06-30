// Code generated by go-swagger; DO NOT EDIT.

package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"

	"github.com/jakub-dzon/k4e-operator/restapi/operations"
	"github.com/jakub-dzon/k4e-operator/restapi/operations/devices"
)

type contextKey string

const AuthKey contextKey = "Auth"

//go:generate mockery -name DevicesAPI -inpkg

/* DevicesAPI  */
type DevicesAPI interface {
	GetDeviceConfiguration(ctx context.Context, params devices.GetDeviceConfigurationParams) middleware.Responder

	RegisterDevice(ctx context.Context, params devices.RegisterDeviceParams) middleware.Responder
}

// Config is configuration for Handler
type Config struct {
	DevicesAPI
	Logger func(string, ...interface{})
	// InnerMiddleware is for the handler executors. These do not apply to the swagger.json document.
	// The middleware executes after routing but before authentication, binding and validation
	InnerMiddleware func(http.Handler) http.Handler

	// Authorizer is used to authorize a request after the Auth function was called using the "Auth*" functions
	// and the principal was stored in the context in the "AuthKey" context value.
	Authorizer func(*http.Request) error

	// AuthAgentAuth Applies when the "X-Secret-Key" header is set
	AuthAgentAuth func(token string) (interface{}, error)

	// Authenticator to use for all APIKey authentication
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// Authenticator to use for all Bearer authentication
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// Authenticator to use for all Basic authentication
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator
}

// Handler returns an http.Handler given the handler configuration
// It mounts all the business logic implementers in the right routing.
func Handler(c Config) (http.Handler, error) {
	h, _, err := HandlerAPI(c)
	return h, err
}

// HandlerAPI returns an http.Handler given the handler configuration
// and the corresponding *Kube4EdgeManagement instance.
// It mounts all the business logic implementers in the right routing.
func HandlerAPI(c Config) (http.Handler, *operations.Kube4EdgeManagementAPI, error) {
	spec, err := loads.Analyzed(swaggerCopy(SwaggerJSON), "")
	if err != nil {
		return nil, nil, fmt.Errorf("analyze swagger: %v", err)
	}
	api := operations.NewKube4EdgeManagementAPI(spec)
	api.ServeError = errors.ServeError
	api.Logger = c.Logger

	if c.APIKeyAuthenticator != nil {
		api.APIKeyAuthenticator = c.APIKeyAuthenticator
	}
	if c.BasicAuthenticator != nil {
		api.BasicAuthenticator = c.BasicAuthenticator
	}
	if c.BearerAuthenticator != nil {
		api.BearerAuthenticator = c.BearerAuthenticator
	}

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.AgentAuthAuth = func(token string) (interface{}, error) {
		if c.AuthAgentAuth == nil {
			return token, nil
		}
		return c.AuthAgentAuth(token)
	}

	api.APIAuthorizer = authorizer(c.Authorizer)
	api.DevicesGetDeviceConfigurationHandler = devices.GetDeviceConfigurationHandlerFunc(func(params devices.GetDeviceConfigurationParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.DevicesAPI.GetDeviceConfiguration(ctx, params)
	})
	api.DevicesRegisterDeviceHandler = devices.RegisterDeviceHandlerFunc(func(params devices.RegisterDeviceParams, principal interface{}) middleware.Responder {
		ctx := params.HTTPRequest.Context()
		ctx = storeAuth(ctx, principal)
		return c.DevicesAPI.RegisterDevice(ctx, params)
	})
	api.ServerShutdown = func() {}
	return api.Serve(c.InnerMiddleware), api, nil
}

// swaggerCopy copies the swagger json to prevent data races in runtime
func swaggerCopy(orig json.RawMessage) json.RawMessage {
	c := make(json.RawMessage, len(orig))
	copy(c, orig)
	return c
}

// authorizer is a helper function to implement the runtime.Authorizer interface.
type authorizer func(*http.Request) error

func (a authorizer) Authorize(req *http.Request, principal interface{}) error {
	if a == nil {
		return nil
	}
	ctx := storeAuth(req.Context(), principal)
	return a(req.WithContext(ctx))
}

func storeAuth(ctx context.Context, principal interface{}) context.Context {
	return context.WithValue(ctx, AuthKey, principal)
}
