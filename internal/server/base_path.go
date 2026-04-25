package server

import (
	"net/http"
	"path"
	"strings"

	"github.com/labstack/echo/v5"
)

func normalizeBasePath(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || trimmed == "/" {
		return "/"
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}
	normalized := path.Clean(trimmed)
	if normalized == "." || normalized == "/" {
		return "/"
	}
	return normalized
}

func configuredBasePath(cfg *Config) string {
	if cfg == nil {
		return "/"
	}
	return normalizeBasePath(cfg.BasePath)
}

func stripBasePathMiddleware(basePath string) echo.MiddlewareFunc {
	basePath = normalizeBasePath(basePath)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if basePath == "/" {
				return next(c)
			}

			req := c.Request()
			strippedPath, ok := stripBasePath(req.URL.Path, basePath)
			if !ok {
				return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
			}

			cloned := req.Clone(req.Context())
			urlCopy := *req.URL
			urlCopy.Path = strippedPath
			urlCopy.RawPath = ""
			cloned.URL = &urlCopy
			cloned.RequestURI = strippedPath
			if urlCopy.RawQuery != "" {
				cloned.RequestURI += "?" + urlCopy.RawQuery
			}
			c.SetRequest(cloned)
			return next(c)
		}
	}
}

func stripBasePath(requestPath, basePath string) (string, bool) {
	basePath = normalizeBasePath(basePath)
	if basePath == "/" {
		if requestPath == "" {
			return "/", true
		}
		return requestPath, true
	}
	if requestPath == basePath {
		return "/", true
	}
	prefix := basePath + "/"
	if !strings.HasPrefix(requestPath, prefix) {
		return "", false
	}
	stripped := strings.TrimPrefix(requestPath, basePath)
	if stripped == "" {
		return "/", true
	}
	return stripped, true
}
