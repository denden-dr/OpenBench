# Route Registration Example

Register the struct methods directly to Fiber's router routes:

```go
authHandler := auth.NewHandler(authService, cfg.JWTAccessExpiry, isDev)
app.Post("/api/v1/auth/signin", authHandler.SignIn)
```
