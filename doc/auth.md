# Auth

There are two authentication flows currently supported in Weave GitOps Core:
1. OIDC (requires an OIDC provider)
1. Cluster user

## OIDC

```mermaid
sequenceDiagram
    participant User
    participant Browser
    participant API
    participant OIDC Provider
    autonumber
    User->>Browser: navigates to dashboard
    Browser->>API: GET /oauth2/userinfo with no cookie 
    API->>Browser: 400 Bad Request
    Browser->>User: show login screen
    User->>Browser: login with OIDC provider
    Browser->>API:GET /oauth2?return_url=/
    Note over Browser,API: return_url set to dashboard
    API->>Browser: 303 See Other with Location set to OIDC Provider
    Note over API,Browser: state cookie is set, HttpOnly (can we also set Secure to true?)
    Browser->>OIDC Provider: GET /auth?..
    Note over Browser,OIDC Provider: Passing client_id,redirect_uri,response_type,scope and state
    OIDC Provider->>Browser: 200 OK
    Browser->>User: show OIDC login screen
    User->>Browser: select Gitlab login
    Browser->>OIDC Provider: GET /auth/gitlab?..
    Note over Browser,OIDC Provider: Passing client_id,redirect_uri,response_type,scope and state
    OIDC Provider->>Browser: 200 OK
    Browser->>User: show Gitlab login / grant access screen
    User->>Browser: login to Gitlab / grant access
    Browser->>OIDC Provider: POST /approval
    OIDC Provider->>Browser: 303 See Other with Location set to /oauth2/callback
    Browser->>API: GET /oauth2/callback
    API->>Browser: 303 See Other with Location set to dashboard
    Note over API,Browser: Set id_token cookie
    Browser->>User: show dashboard
```

Endpoints:
- 2: /oauth2/userinfo calls UserInfo() handler
- 6: /oauth2 calls OAuth2Flow() handler
- 18: /oauth2/callback calls Callback() handler

Notes:
- 3: This should return 401 Unauthorized

## Cluster user