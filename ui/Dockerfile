FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine

# Create nginx user and set up directories with proper permissions
RUN addgroup -g 1000 -S moogie && \
    adduser -u 1000 -S moogie -G moogie && \
    mkdir -p /var/cache/nginx /var/run /var/log/nginx && \
    chown -R moogie:moogie /var/cache/nginx /var/run /var/log/nginx /etc/nginx/conf.d /usr/share/nginx/html && \
    chmod -R 755 /var/cache/nginx /var/run /var/log/nginx

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf

# Change ownership of the files
RUN chown -R moogie:moogie /usr/share/nginx/html

# Switch to non-root user
USER moogie

EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
