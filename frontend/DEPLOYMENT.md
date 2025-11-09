# Frontend Deployment Guide

This guide covers deploying the Mini Attendance Frontend application to various environments.

## Table of Contents

1. [Local Development](#local-development)
2. [Docker Deployment](#docker-deployment)
3. [Cloud Deployment](#cloud-deployment)
4. [Environment Configuration](#environment-configuration)
5. [Troubleshooting](#troubleshooting)

## Local Development

### Prerequisites

- Node.js 18 or higher
- npm or yarn

### Setup

1. Install dependencies:

```bash
npm install
```

2. Create `.env.local`:

```bash
cp .env.example .env.local
```

3. Update environment variables:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

4. Start development server:

```bash
npm run dev
```

5. Open [http://localhost:3000](http://localhost:3000)

## Docker Deployment

### Build Docker Image

```bash
docker build -t mini-attendance-frontend:latest .
```

### Run Container Locally

```bash
docker run -p 3000:3000 \
  -e NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1 \
  mini-attendance-frontend:latest
```

### Docker Compose

Use the root `docker-compose.yml` to run the entire stack:

```bash
cd ..
docker-compose up -d
```

This will start:
- Frontend (port 3000)
- Backend (port 8080)
- PostgreSQL (port 5432)
- Redis (port 6379)
- Kafka (port 9092)
- Grafana (port 3001)
- Loki (port 3100)

## Cloud Deployment

### AWS ECS

1. Create ECR repository:

```bash
aws ecr create-repository --repository-name mini-attendance-frontend
```

2. Build and push image:

```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com

docker build -t mini-attendance-frontend:latest .

docker tag mini-attendance-frontend:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/mini-attendance-frontend:latest

docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/mini-attendance-frontend:latest
```

3. Create ECS task definition and service

### Google Cloud Run

1. Build and push image:

```bash
gcloud builds submit --tag gcr.io/<project-id>/mini-attendance-frontend
```

2. Deploy to Cloud Run:

```bash
gcloud run deploy mini-attendance-frontend \
  --image gcr.io/<project-id>/mini-attendance-frontend \
  --platform managed \
  --region us-central1 \
  --set-env-vars NEXT_PUBLIC_API_URL=https://api.example.com/api/v1
```

### Azure Container Instances

1. Build and push image:

```bash
az acr build --registry <registry-name> --image mini-attendance-frontend:latest .
```

2. Deploy:

```bash
az container create \
  --resource-group <resource-group> \
  --name mini-attendance-frontend \
  --image <registry-name>.azurecr.io/mini-attendance-frontend:latest \
  --ports 3000 \
  --environment-variables NEXT_PUBLIC_API_URL=https://api.example.com/api/v1
```

### Vercel (Recommended for Next.js)

1. Push code to GitHub

2. Connect repository to Vercel

3. Set environment variables in Vercel dashboard:

```
NEXT_PUBLIC_API_URL=https://api.example.com/api/v1
```

4. Deploy automatically on push

## Environment Configuration

### Development

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_LOG_LEVEL=debug
NEXT_PUBLIC_ENABLE_ANALYTICS=false
```

### Staging

```env
NEXT_PUBLIC_API_URL=https://api-staging.example.com/api/v1
NEXT_PUBLIC_LOG_LEVEL=info
NEXT_PUBLIC_ENABLE_ANALYTICS=true
```

### Production

```env
NEXT_PUBLIC_API_URL=https://api.example.com/api/v1
NEXT_PUBLIC_LOG_LEVEL=warn
NEXT_PUBLIC_ENABLE_ANALYTICS=true
NEXT_PUBLIC_ENABLE_SENTRY=true
```

## Performance Optimization

### Build Optimization

```bash
npm run build
```

This creates an optimized production build with:
- Code splitting
- Image optimization
- CSS minification
- JavaScript minification

### Caching Strategy

- Static assets: 1 year cache
- HTML: No cache (always fresh)
- API responses: Cached via RTK Query

### CDN Configuration

For production, use a CDN like CloudFlare:

1. Point domain to CDN
2. Configure cache rules
3. Enable compression
4. Enable HTTP/2

## Monitoring

### Health Check

The application exposes a health check endpoint:

```bash
curl http://localhost:3000/api/health
```

### Logging

Logs are output to stdout and can be collected by:
- Docker logs
- CloudWatch
- Stackdriver
- ELK Stack

### Error Tracking

Configure Sentry for error tracking:

1. Create Sentry project
2. Set `NEXT_PUBLIC_SENTRY_DSN` environment variable
3. Errors are automatically reported

## Security

### HTTPS

Always use HTTPS in production:

```bash
# Redirect HTTP to HTTPS
docker run -p 80:80 -p 443:443 \
  -e NEXT_PUBLIC_API_URL=https://api.example.com/api/v1 \
  mini-attendance-frontend:latest
```

### CORS

CORS is configured in `next.config.js`. Update for production:

```javascript
headers: async () => {
  return [
    {
      source: '/api/:path*',
      headers: [
        { key: 'Access-Control-Allow-Origin', value: 'https://example.com' },
      ],
    },
  ];
};
```

### Environment Variables

Never commit `.env.local` or `.env.production.local`:

```bash
echo ".env.local" >> .gitignore
echo ".env.production.local" >> .gitignore
```

## Scaling

### Horizontal Scaling

Deploy multiple instances behind a load balancer:

```bash
docker run -p 3000:3000 mini-attendance-frontend:latest
docker run -p 3001:3000 mini-attendance-frontend:latest
docker run -p 3002:3000 mini-attendance-frontend:latest
```

### Load Balancing

Use Nginx or HAProxy:

```nginx
upstream frontend {
  server localhost:3000;
  server localhost:3001;
  server localhost:3002;
}

server {
  listen 80;
  location / {
    proxy_pass http://frontend;
  }
}
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 3000
lsof -i :3000

# Kill process
kill -9 <PID>
```

### API Connection Issues

1. Check `NEXT_PUBLIC_API_URL` is correct
2. Verify backend is running
3. Check CORS configuration
4. Check network connectivity

### Build Failures

```bash
# Clear cache
rm -rf .next node_modules

# Reinstall
npm install

# Rebuild
npm run build
```

### Memory Issues

Increase Node.js memory:

```bash
NODE_OPTIONS=--max-old-space-size=4096 npm run build
```

## Rollback

### Docker

```bash
# Tag previous version
docker tag mini-attendance-frontend:v1.0.0 mini-attendance-frontend:latest

# Run previous version
docker run -p 3000:3000 mini-attendance-frontend:v1.0.0
```

### Vercel

Automatic rollback available in Vercel dashboard

## Support

For deployment issues, check:
- Application logs
- Backend API logs
- Network connectivity
- Environment variables
- Docker/container logs
