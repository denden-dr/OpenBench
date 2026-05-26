FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/build ./build

RUN chown -R node:node /app
USER node

EXPOSE 5173
ENV PORT=5173
ENV HOST=0.0.0.0

CMD ["node", "build"]
