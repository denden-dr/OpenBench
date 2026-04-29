# User Database Design

This document contains the Entity-Relationship (ER) diagram for the local proxy `users` table, which serves to synchronize profile data from Supabase Auth.

```mermaid
erDiagram
    USERS {
        UUID id PK "Primary Key, matches Supabase auth.users id"
        VARCHAR email UK "Unique, Not Null"
        TEXT full_name "Nullable"
        TEXT avatar_url "Nullable"
        TIMESTAMPTZ updated_at "Not Null, Default NOW(), tracks last sync"
    }
```

### Table Details

- **`users`**: The core local representation of a user.
  - `id`: Used as the primary lookup and relational key for all local business entities. Directly mirrors the `sub` claim / `uuid` from the Supabase JWT.
  - `email`: Indexed for fast email-based lookups, guaranteed unique.
  - `full_name` & `avatar_url`: Supplementary profile information extracted from Supabase metadata during login.
  - `updated_at`: The timestamp tracking when the local database was last synchronized with Supabase's data payload.
