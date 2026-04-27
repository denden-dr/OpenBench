# Authentication Flow Diagram

This diagram illustrates the sequence of operations for authenticating a request and synchronizing the user profile with the local database.

```mermaid
sequenceDiagram
    autonumber
    participant Client
    participant MW as AuthMiddleware (Fiber)
    participant AS as AuthService
    participant JWKS as JWKS Cache (lestrrat-go/jwx)
    participant UR as UserRepository (sqlx)
    participant DB as PostgreSQL

    Note over Client, DB: Request Lifecycle

    Client->>MW: Request [Authorization: Bearer <JWT>]
    
    rect rgb(240, 240, 240)
        Note right of MW: Token Extraction
        MW->>MW: Extract raw JWT from header
        alt is malformed or missing
            MW-->>Client: 401 Unauthorized (JSON Error)
        end
    end

    MW->>AS: VerifyAndSync(ctx, rawToken)

    rect rgb(245, 245, 255)
        Note right of AS: JWT Verification
        AS->>JWKS: Get KeySet (Local Cache)
        JWKS-->>AS: KeySet
        AS->>AS: jwt.Parse (Verify Signature, Exp, Iss)
        alt is invalid or expired
            AS-->>MW: error (Invalid Token)
            MW-->>Client: 401 Unauthorized
        end
    end

    rect rgb(245, 255, 245)
        Note right of AS: User Synchronization (JIT)
        AS->>UR: UpsertFromAuth(ctx, sub, email, metadata)
        UR->>DB: INSERT ... ON CONFLICT (id) DO UPDATE
        DB-->>UR: rows (User Record)
        UR-->>AS: domain.User
    end

    AS-->>MW: domain.User
    
    rect rgb(240, 240, 240)
        Note right of MW: Context Injection
        MW->>MW: c.Locals("user", user)
    end

    MW->>MW: c.Next()
    Note over MW: Execution proceeds to Handler
    MW-->>Client: 200 OK / Response
```

## Flow Description

1.  **Request Initiation**: The client sends an HTTP request with the `Authorization: Bearer <JWT>` header.
2.  **Middleware Extraction**: The Fiber middleware extracts the token. If missing or the scheme is not `Bearer`, it returns an immediate `401`.
3.  **Service Delegation**: The middleware calls `AuthService.VerifyAndSync`.
4.  **Signature Verification**: The service uses cached JWKS keys (fetched from Supabase at startup/refreshed periodically) to verify the token's cryptographic signature.
5.  **Claim Validation**: The service validates standard claims (`exp`, `iss`).
6.  **Just-In-Time (JIT) Sync**: The service extracts the user identity (from the `sub` claim) and updates/creates the local user record in PostgreSQL via the `UserRepository`.
7.  **Success Path**: The local `User` object is returned to the middleware.
8.  **Context Storage**: The middleware stores the `User` object in Fiber's `c.Locals`, making it available to all subsequent handlers in the request chain.
9.  **Next**: The middleware calls `c.Next()`, passing control to the business logic handler.
