export interface User {
  id: string
  email: string
  role: string
}

export interface LoginCredentials {
  email: string
  password: string
}

export interface LoginResponse {
  data: {
    access_token: string
    expires_at: string
    user: User
  }
}

export interface UserResponse {
  data: User
}

export interface ApiProblemDetails {
  type?: string
  title?: string
  status?: number
  detail?: string
  instance?: string
}
