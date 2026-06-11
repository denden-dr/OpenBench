# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY package*.json ./

# Install all dependencies (including devDependencies)
RUN npm ci

# Copy source code
COPY . .

# Build the SvelteKit application (adapter-node outputs to ./build)
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Set ownership of workspace to node user
RUN chown node:node /app

# Switch to non-root user context
USER node

# Copy package manifests and install only production dependencies
COPY --chown=node:node package*.json ./
RUN npm ci --omit=dev

# Copy built application from builder
COPY --chown=node:node --from=builder /app/build ./build

# Expose SvelteKit node adapter port (default 3000)
EXPOSE 3000

ENV PORT=3000
ENV NODE_ENV=production

CMD ["node", "build"]
