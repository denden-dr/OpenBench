# Payload Contracts and Schema Verification

Always map structures between TypeScript and Go before implementing code to avoid missing parameters:
- **Go JSON tags**: Verify struct tags (e.g., `json:"user_id"` matches TypeScript's camelCase or snake_case).
- **Fallback safety**: When parsing API responses in TypeScript, do not assume cache parameters or fallbacks unless explicitly handled.

### Equivalent Mappings

| Svelte TypeScript Interface | Go Backend Struct Equivalent |
|-----------------------------|------------------------------|
| `userId: string`            | `UserID string json:"user_id"` |
| `email: string`             | `Email string json:"email"` |

### API Response Envelope Formatting
To prevent layout guards from reading `undefined` fields, always parse backend JSON responses through the standardized `data` envelope:

```typescript
// Svelte Client-side API Response Interface
export interface APIResponse<T> {
  code: number;
  message: string;
  data: T;
}

// Example fetch parse:
const response = await fetch('/api/v1/resource');
const body: APIResponse<MyDataType> = await response.json();
const actualData = body.data; // Access data inside the envelope (ST-07)
```

