# Request X-URL Traefik Plugin

A Traefik middleware plugin that adds the **full original request URL**
(`scheme + host + path + query`) to an HTTP header before forwarding the request
to the backend service.

This plugin is useful for legacy backends, logging, auditing, or applications
that need access to the original client-facing URL.

---

## Features

- Adds full request URL to a configurable HTTP header
- Preserves query string
- Supports HTTP and HTTPS
- Uses standard `X-Forwarded-*` headers
- Lightweight (no external dependencies)
- Compatible with Traefik v2 and v3

---

## How It Works

The plugin constructs the full request URL using:

- `X-Forwarded-Proto` (fallback: TLS detection or `http`)
- `X-Forwarded-Host` (fallback: `Host`)
- `RequestURI` (path + query string)

Resulting value example:

```

[https://example.com/api/v1/test?id=123](https://example.com/api/v1/test?id=123)

````

This value is then injected into the configured request header.

---

## Configuration

### Static Configuration

Enable the plugin in Traefik static configuration:

```yaml
experimental:
  plugins:
    requesturl:
      moduleName: github.com/akotlyar/plugin-request-x-url
      version: v1.0.0
````

---

### Dynamic Configuration (Middleware)

```yaml
http:
  middlewares:
    add-full-request-url:
      plugin:
        requesturl:
          headerName: X-Real-Url
```

Attach the middleware to a router:

```yaml
http:
  routers:
    my-service:
      rule: "Host(`example.com`)"
      entryPoints:
        - websecure
      middlewares:
        - add-full-request-url
      service: my-service
```

---

## Configuration Options

| Option       | Type   | Default         | Description                       |
| ------------ | ------ | --------------- | --------------------------------- |
| `headerName` | string | `X-Request-Url` | HTTP header to store the full URL |

