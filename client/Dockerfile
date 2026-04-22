# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm install

# Copy source code and build the application
COPY . .
RUN npm run build

# Runtime stage
FROM node:22-alpine

WORKDIR /app

# Copy built application and production dependencies
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json

# If you use modules that are not bundled, you might need to install prod deps here
# However, for SvelteKit with adapter-node, the build folder usually contains everything needed 
# or you might need to copy node_modules. Let's install prod deps just in case.
RUN npm install --omit=dev

EXPOSE 3000

ENV NODE_ENV=production

CMD ["node", "build"]
