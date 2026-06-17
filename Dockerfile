# Use a stable Node base image for build
FROM node:22-alpine AS builder
WORKDIR /app

# Install dependencies first to leverage Docker cache
COPY package.json package-lock.json* ./
COPY pnpm-lock.yaml* ./
COPY yarn.lock* ./
RUN npm install --production=false

# Copy source and build
COPY . .
RUN npm run build

# Production image
FROM node:22-alpine AS runner
WORKDIR /app
ENV NODE_ENV=production

# Copy only needed files from builder
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package.json ./package.json

EXPOSE 3000
CMD ["npm", "run", "start"]
