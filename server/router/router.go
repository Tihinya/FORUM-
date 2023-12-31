package router

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Middleware func(http.Handler) http.Handler
type route struct {
	regex           *regexp.Regexp
	handler         http.Handler
	localMiddleware []Middleware
	middlewareAfter []Middleware
}
type router struct {
	routes           map[string][]route
	globalMiddleware []Middleware
}
type CORS struct {
	Origin      string
	Methods     []string
	Headers     []string
	Credentials bool
}

func NewRouter() *router {
	return &router{
		routes:           map[string][]route{},
		globalMiddleware: []Middleware{},
	}
}

// NewRouter creates new route and appends it to Router,
// method specifies the method that is allowed,
// regexp must contain a named group,
// parsed value will be accessable through handler's context via GetField function
// as an optional argument NewRoute can take any number of Middleware and apply them to itself
//
// regexp example -> https://regex101.com/r/84S9iL/1
func (r *router) NewRoute(method, regexpString string, handler http.HandlerFunc, middleware ...Middleware) {
	regex := regexp.MustCompile("^" + regexpString + "$")
	method = strings.ToUpper(method)
	r.routes[method] = append(r.routes[method], route{
		regex,
		handler,
		middleware,
		make([]Middleware, 0),
	})
}

// AddGlobalMiddleware takes any amount of user defined middlewares and
// applies them consecutively in all routes
func (r *router) AddGlobalMiddleware(middleware ...Middleware) {
	r.globalMiddleware = append(r.globalMiddleware, middleware...)
}

// If there is a exact match Serve redirects to the associated handler function.
//
// When matched, regular expression groups are used as key value pairs
// accessible in handler's context via GetField function
func (r *router) Serve(w http.ResponseWriter, req *http.Request) {
	r.serve(w, req)
}

// Same usage as in Serve function, but also adds specified CORS headers
func (r *router) ServeWithCORS(c CORS) http.HandlerFunc {
	headers := make(map[string]string)
	if c.Origin != "" {
		headers["Access-Control-Allow-Origin"] = c.Origin
	}
	if len(c.Methods) > 0 {
		headers["Access-Control-Allow-Methods"] = strings.Join(c.Methods, ", ")
	}
	if len(c.Headers) > 0 {
		headers["Access-Control-Allow-Headers"] = strings.Join(c.Headers, ", ")
	}
	if c.Credentials {
		headers["Access-Control-Allow-Credentials"] = "true"
	}
	return func(w http.ResponseWriter, req *http.Request) {
		for header, value := range headers {
			w.Header().Set(header, value)
		}
		if req.Method == "OPTIONS" {
			return
		}
		r.serve(w, req)
	}
}
func (r *router) serve(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes[strings.ToUpper(req.Method)] {
		match := route.regex.FindStringSubmatch(req.URL.Path)
		if len(match) > 0 {
			matchMap := make(map[string]string)
			groupName := route.regex.SubexpNames()
			// map group name(key) to submatched result
			// these arrays have one to one relationship
			for i := 1; i < len(match); i++ {
				matchMap[groupName[i]] = match[i]
			}
			ctx := context.WithValue(req.Context(), struct{}{}, matchMap)
			req = req.WithContext(ctx)
			handler := route.handler
			for i := len(route.localMiddleware) - 1; i >= 0; i-- {
				handler = route.localMiddleware[i](handler)
			}
			for i := len(r.globalMiddleware) - 1; i >= 0; i-- {
				handler = r.globalMiddleware[i](handler)
			}
			handler.ServeHTTP(w, req)
			return
		}
	}
}

// Returns the string value of the given key from matched URL variables
func GetFieldString(r *http.Request, name string) (string, error) {
	fields, ok := r.Context().Value(struct{}{}).(map[string]string)
	if !ok {
		return "", fmt.Errorf("internal error: no fileds in context")
	}
	field, exist := fields[name]
	if !exist {
		return "", fmt.Errorf("no such variable in request: %v", name)
	}
	return field, nil
}

// Returns the integer value of the given key from matched URL variables
func GetFieldInt(r *http.Request, name string) (int, error) {
	field, err := GetFieldString(r, name)
	if err != nil {
		return 0, err
	}
	fieldInt, err := strconv.Atoi(field)
	if err != nil {
		return 0, err
	}
	return fieldInt, nil
}
