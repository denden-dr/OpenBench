---
name: sveltekit-supabase-auth
description: Use when implementing authentication in the SvelteKit frontend — Supabase client setup with @supabase/ssr, cookie-based session management, server-side hooks, and secure API communication.
---

# SvelteKit + Supabase Secure Cookie Auth

## Overview

SvelteKit handles **user-facing auth** and manages the session in **HttpOnly cookies**. This prevents XSS-based token theft. The Go backend reads these cookies directly.

## Setup

```bash
npm install @supabase/ssr @supabase/supabase-js
```

## Client Initialization

```ts
// src/lib/supabase.ts (Browser Client)
import { createBrowserClient } from '@supabase/ssr'
import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public'

export const supabase = createBrowserClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY)
```

## Server Hooks (Cookie Management)

This is where the "magic" happens. SvelteKit manages the session in cookies.

```ts
// src/hooks.server.ts
import { createServerClient } from '@supabase/ssr'
import { type Handle } from '@sveltejs/kit'
import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public'

export const handle: Handle = async ({ event, resolve }) => {
    event.locals.supabase = createServerClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY, {
        cookies: {
            getAll: () => event.cookies.getAll(),
            setAll: (cookiesToSet) => {
                cookiesToSet.forEach(({ name, value, options }) => {
                    event.cookies.set(name, value, { ...options, path: '/' })
                })
            },
        },
    })

    event.locals.safeGetSession = async () => {
        const { data: { session } } = await event.locals.supabase.auth.getSession()
        if (!session) return { session: null, user: null }
        return { session, user: session.user }
    }

    return resolve(event, {
        filterSerializedResponseHeaders(name) {
            return name === 'content-range' || name === 'x-supabase-parse'
        },
    })
}
```

## API Client (Cookie-Based)

When calling the Go backend, **no manual token attachment is needed** if the API is on the same base domain.

```ts
// src/lib/api.ts
const API_BASE = import.meta.env.VITE_BACKEND_API_URL || '/api/v1'

export async function apiFetch(path: string, options: RequestInit = {}) {
    // browser attaches cookies automatically
    const res = await fetch(`${API_BASE}${path}`, { 
        ...options,
        // include credentials if cross-subdomain
        credentials: 'include' 
    })
    
    if (res.status === 401) {
        // Handle unauthorized (e.g., redirect to login)
    }
    
    return await res.json()
}
```

## Protected Routes (Server-Side)

```ts
// src/routes/(protected)/+layout.server.ts
import { redirect } from '@sveltejs/kit'

export const load = async ({ locals: { safeGetSession } }) => {
    const { session } = await safeGetSession()

    if (!session) {
        throw redirect(303, '/login')
    }

    return { session }
}
```

## Auth Actions

Use SvelteKit Form Actions for login/logout to ensure cookies are set/cleared correctly.

```ts
// src/routes/login/+page.server.ts
export const actions = {
    login: async ({ request, locals: { supabase } }) => {
        const formData = await request.formData()
        const email = formData.get('email')
        const password = formData.get('password')

        const { error } = await supabase.auth.signInWithPassword({ email, password })
        if (error) return { success: false, message: error.message }
        
        throw redirect(303, '/dashboard')
    }
}
```

## Quick Reference

| Concern | Where |
|---------|-------|
| Supabase Client | `src/lib/supabase.ts` |
| Cookie Management | `src/hooks.server.ts` |
| Route Guard | `src/routes/(protected)/+layout.server.ts` |
| Login Action | `src/routes/login/+page.server.ts` |

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Accessing `localStorage` for tokens | Use `event.locals.supabase` on the server |
| Not setting `credentials: 'include'` | Required if API is on a subdomain (e.g. `api.app.io`) |
| Client-side only login | Use Form Actions to ensure cookies are set via `Set-Cookie` |
| Forgetting `path: '/'` in cookie options | Session won't be available across all routes |
