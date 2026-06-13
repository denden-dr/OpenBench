---
name: secure-containerization
description: Use when containerizing applications (Go, SvelteKit, Node.js) with Docker or Podman, ensuring non-root execution, pinned image tags, multi-stage builds, and dockerignore configurations. Do not use for virtual machine orchestration or Kubernetes configuration files.
version: 1.0.0
---

# Secure Containerization

## Overview
Container security is crucial for production deployments. Running containers as `root` or baking sensitive credentials into image layers creates severe security risks. Safe containerization relies on multi-stage builds, pinned stable base tags, non-root user execution, and strict ignore patterns.

## When to Use
- Writing or updating a Dockerfile for backend or frontend applications.
- Setting up or reviewing container build files.
- Configuring `.dockerignore` files to prevent secret leakage.

## Step-by-Step Instructions

1. **Configure `.dockerignore`**: Ensure a `.dockerignore` file is placed in the service context root, excluding local `.env` files, build caches, git history, and local dependencies.
2. **Implement Backend Containerization**: Read `assets/Dockerfile-backend.template` and construct the Go backend Dockerfile, utilizing multi-stage builds, static compilation, and user restriction mappings.
3. **Implement Frontend Containerization**: Read `assets/Dockerfile-frontend.template` and construct the Node/SvelteKit frontend Dockerfile, leveraging multi-stage assets building and node non-root group context.
4. **Enforce Non-Root Execution**: Verify the final Dockerfile switches execution context to a restricted non-root user (`USER appuser` or `USER node`) to prevent breakout vulnerabilities.

## Common Mistakes
- **No `.dockerignore`**: Leaving local `.env` files in the build folder, causing them to be copied by `COPY . .` and permanently baked into the image history.
- **Running as Root**: Leaving the default `root` user context. If a remote code execution vulnerability occurs, the attacker gains full root privileges inside the container, easing breakouts to the host.
- **Mutable Base Tags**: Using `latest` or generic tags (e.g., `FROM golang:latest` or `FROM node:alpine`). When base images update upstream, your build pipeline can break or inherit unexpected security bugs.
- **Baking Node Modules**: Copying local development `node_modules` instead of running `npm ci --omit=dev` inside the builder context, leading to bloated images and environment mismatch issues.
