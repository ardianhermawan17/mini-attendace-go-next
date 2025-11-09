# Deployment Guide - Mini Attendance System

## 📋 Table of Contents

1. [Local Development](#local-development)
2. [Docker Deployment](#docker-deployment)
3. [Cloud Deployment](#cloud-deployment)
4. [Kubernetes Deployment](#kubernetes-deployment)
5. [Production Checklist](#production-checklist)
6. [Monitoring & Maintenance](#monitoring--maintenance)

## Local Development

### Prerequisites

- Go 1.21+
- PostgreSQL 16+
- Redis 7+
- Kafka 7.5+
- Docker & Docker Compose (optional)

### Setup Steps

1. **Install Go dependencies**
   ```bash
   cd backend
   go mod download
   ```

2. **Set up PostgreSQL**
   ```bash
   # Create database
   createdb mini_attendance
   
   # Or using Docker
   docker run -d --name postgres \
     -e POSTGRES_PASSWORD=postgres \
     -e POSTGRES_DB=mini_attendance \
     -p 5432:5432 \
     postgres:16-alpine
   ```

3. **Set up Redis**
   ```bash
   # Using Docker
   docker run -d --name redis \
     -p 6379:6379 \
     redis:7-alpine
   ```

4. **Set up Kafka**
   ```bash
   # Using Docker Compose
   docker-compose up -d kafka zookeeper
   ```

5. **Configure environment**
   ```bash
   cp .env.example .env
   # Edit .env with your settings
   ```

6. **Run application**
   ```bash
   go run cmd/attendance-api/main.go
   ```

## Docker Deployment

### Using Docker Compose (Recommended for Local)

1. **Start all services**
   ```bash
   docker-compose up -d
   ```

2. **Verify services**
   ```bash
   docker-compose ps
   ```

3. **View logs**
   ```bash
   docker-compose logs -f backend
   ```

4. **Stop services**
   ```bash
   docker-compose down
   ```

### Building Custom Docker Image

1. **Build image**
   ```bash
   cd backend
   docker build -t mini-attendance:latest .
   ```

2. **Run container**
   ```bash
   docker run -d \
     --name mini-attendance \
     -p 8080:8080 \
     -e DB_HOST=postgres \
     -e REDIS_HOST=redis \
     -e KAFKA_BROKERS=kafka:29092 \
     mini-attendance:latest
   ```

3. **Push to registry**
   ```bash
   docker tag mini-attendance:latest your-registry/mini-attendance:latest
   docker push your-registry/mini-attendance:latest
   ```

## Cloud Deployment

### AWS Deployment

#### Using ECS (Elastic Container Service)

1. **Create ECR repository**
   ```bash
   aws ecr create-repository --repository-name mini-attendance
   ```

2. **Push image**
   ```bash
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin <account-id>.dkr.ecr.us-east-1.amazonaws.com
   docker tag mini-attendance:latest <account-id>.dkr.ecr.us-east-1.amazonaws.com/mini-attendance:latest
   docker push <account-id>.dkr.ecr.us-east-1.amazonaws.com/mini-attendance:latest
   ```

3. **Create RDS PostgreSQL**
   ```bash
   aws rds create-db-instance \
     --db-instance-identifier mini-attendance-db \
     --db-instance-class db.t3.micro \
     --engine postgres \
     --master-username postgres \
     --master-user-password <password> \
     --allocated-storage 20
   ```

4. **Create ElastiCache Redis**
   ```bash
   aws elasticache create-cache-cluster \
     --cache-cluster-id mini-attendance-cache \
     --cache-node-type cache.t3.micro \
     --engine redis \
     --num-cache-nodes 1
   ```

5. **Create MSK Kafka cluster**
   ```bash
   aws kafka create-cluster \
     --cluster-name mini-attendance-kafka \
     --broker-node-group-info InstanceType=kafka.t3.small,ClientSubnets=<subnet-ids>
   ```

#### Using ECS Task Definition

```json
{
  "family": "mini-attendance",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "256",
  "memory": "512",
  "containerDefinitions": [
    {
      "name": "mini-attendance",
      "image": "<account-id>.dkr.ecr.us-east-1.amazonaws.com/mini-attendance:latest",
      "portMappings": [
        {
          "containerPort": 8080,
          "hostPort": 8080,
          "protocol": "tcp"
        }
      ],
      "environment": [
        {
          "name": "DB_HOST",
          "value": "<rds-endpoint>"
        },
        {
          "name": "REDIS_HOST",
          "value": "<elasticache-endpoint>"
        },
        {
          "name": "KAFKA_BROKERS",
          "value": "<msk-brokers>"
        }
      ],
      "logConfiguration": {
        "logDriver": "awslogs",
        "options": {
          "awslogs-group": "/ecs/mini-attendance",
          "awslogs-region": "us-east-1",
          "awslogs-stream-prefix": "ecs"
        }
      }
    }
  ]
}
```

### Google Cloud Deployment

#### Using Cloud Run

1. **Build and push image**
   ```bash
   gcloud builds submit --tag gcr.io/<project-id>/mini-attendance
   ```

2. **Deploy to Cloud Run**
   ```bash
   gcloud run deploy mini-attendance \
     --image gcr.io/<project-id>/mini-attendance \
     --platform managed \
     --region us-central1 \
     --set-env-vars DB_HOST=<cloudsql-ip>,REDIS_HOST=<memorystore-ip>
   ```

#### Using GKE (Google Kubernetes Engine)

See [Kubernetes Deployment](#kubernetes-deployment) section.

### Azure Deployment

#### Using Container Instances

```bash
az container create \
  --resource-group mini-attendance \
  --name mini-attendance-api \
  --image <registry>/mini-attendance:latest \
  --ports 8080 \
  --environment-variables \
    DB_HOST=<db-host> \
    REDIS_HOST=<redis-host> \
    KAFKA_BROKERS=<kafka-brokers>
```

#### Using App Service

```bash
az appservice plan create \
  --name mini-attendance-plan \
  --resource-group mini-attendance \
  --sku B1 --is-linux

az webapp create \
  --resource-group mini-attendance \
  --plan mini-attendance-plan \
  --name mini-attendance-api \
  --deployment-container-image-name <registry>/mini-attendance:latest
```

## Kubernetes Deployment

### Prerequisites

- kubectl installed
- Kubernetes cluster (1.24+)
- Helm (optional)

### Manual Deployment

1. **Create namespace**
   ```bash
   kubectl create namespace mini-attendance
   ```

2. **Create ConfigMap**
   ```bash
   kubectl create configmap app-config \
     --from-literal=APP_NAME=mini-attendance \
     --from-literal=ENVIRONMENT=production \
     -n mini-attendance
   ```

3. **Create Secrets**
   ```bash
   kubectl create secret generic app-secrets \
     --from-literal=DB_PASSWORD=<password> \
     --from-literal=JWT_SECRET=<secret> \
     -n mini-attendance
   ```

4. **Deploy PostgreSQL**
   ```yaml
   apiVersion: v1
   kind: PersistentVolumeClaim
   metadata:
     name: postgres-pvc
     namespace: mini-attendance
   spec:
     accessModes:
       - ReadWriteOnce
     resources:
       requests:
         storage: 10Gi
   ---
   apiVersion: apps/v1
   kind: StatefulSet
   metadata:
     name: postgres
     namespace: mini-attendance
   spec:
     serviceName: postgres
     replicas: 1
     selector:
       matchLabels:
         app: postgres
     template:
       metadata:
         labels:
           app: postgres
       spec:
         containers:
         - name: postgres
           image: postgres:16-alpine
           ports:
           - containerPort: 5432
           env:
           - name: POSTGRES_DB
             value: mini_attendance
           - name: POSTGRES_PASSWORD
             valueFrom:
               secretKeyRef:
                 name: app-secrets
                 key: DB_PASSWORD
           volumeMounts:
           - name: postgres-storage
             mountPath: /var/lib/postgresql/data
         volumeClaims:
         - metadata:
             name: postgres-storage
           spec:
             accessModes: [ "ReadWriteOnce" ]
             resources:
               requests:
                 storage: 10Gi
   ```

5. **Deploy Redis**
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: redis
     namespace: mini-attendance
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: redis
     template:
       metadata:
         labels:
           app: redis
       spec:
         containers:
         - name: redis
           image: redis:7-alpine
           ports:
           - containerPort: 6379
   ```

6. **Deploy Backend API**
   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: mini-attendance-api
     namespace: mini-attendance
   spec:
     replicas: 3
     selector:
       matchLabels:
         app: mini-attendance-api
     template:
       metadata:
         labels:
           app: mini-attendance-api
       spec:
         containers:
         - name: api
           image: <registry>/mini-attendance:latest
           ports:
           - containerPort: 8080
           env:
           - name: DB_HOST
             value: postgres
           - name: REDIS_HOST
             value: redis
           - name: DB_PASSWORD
             valueFrom:
               secretKeyRef:
                 name: app-secrets
                 key: DB_PASSWORD
           - name: JWT_SECRET
             valueFrom:
               secretKeyRef:
                 name: app-secrets
                 key: JWT_SECRET
           livenessProbe:
             httpGet:
               path: /health
               port: 8080
             initialDelaySeconds: 10
             periodSeconds: 10
           readinessProbe:
             httpGet:
               path: /health
               port: 8080
             initialDelaySeconds: 5
             periodSeconds: 5
           resources:
             requests:
               memory: "256Mi"
               cpu: "250m"
             limits:
               memory: "512Mi"
               cpu: "500m"
   ---
   apiVersion: v1
   kind: Service
   metadata:
     name: mini-attendance-api
     namespace: mini-attendance
   spec:
     type: LoadBalancer
     selector:
       app: mini-attendance-api
     ports:
     - protocol: TCP
       port: 80
       targetPort: 8080
   ```

### Using Helm

1. **Create Helm chart**
   ```bash
   helm create mini-attendance
   ```

2. **Deploy with Helm**
   ```bash
   helm install mini-attendance ./mini-attendance \
     --namespace mini-attendance \
     --values values.yaml
   ```

## Production Checklist

### Security

- [ ] Change default passwords
- [ ] Enable HTTPS/TLS
- [ ] Configure firewall rules
- [ ] Set up VPN/private network
- [ ] Enable database encryption
- [ ] Configure backup encryption
- [ ] Set up API rate limiting
- [ ] Enable audit logging
- [ ] Configure CORS properly
- [ ] Use environment variables for secrets

### Performance

- [ ] Configure database connection pooling
- [ ] Set up Redis replication
- [ ] Configure Kafka replication
- [ ] Enable caching headers
- [ ] Set up CDN for static assets
- [ ] Configure load balancing
- [ ] Set up auto-scaling
- [ ] Monitor query performance

### Reliability

- [ ] Set up automated backups
- [ ] Configure backup retention
- [ ] Test disaster recovery
- [ ] Set up health checks
- [ ] Configure alerting
- [ ] Set up monitoring
- [ ] Configure log aggregation
- [ ] Set up error tracking

### Compliance

- [ ] Review data retention policies
- [ ] Configure audit logging
- [ ] Set up compliance monitoring
- [ ] Document security procedures
- [ ] Configure access controls
- [ ] Set up encryption at rest
- [ ] Set up encryption in transit

## Monitoring & Maintenance

### Health Monitoring

```bash
# Check API health
curl https://api.example.com/health

# Check database
psql -h <host> -U postgres -d mini_attendance -c "SELECT 1"

# Check Redis
redis-cli -h <host> ping

# Check Kafka
kafka-broker-api-versions.sh --bootstrap-server <host>:9092
```

### Log Monitoring

```bash
# View logs in Grafana
# Navigate to http://grafana.example.com
# Query: {app="mini-attendance"}

# View logs in Loki
curl -G -s "http://loki.example.com/loki/api/v1/query" \
  --data-urlencode 'query={app="mini-attendance"}'
```

### Backup & Recovery

```bash
# Backup PostgreSQL
pg_dump -h <host> -U postgres mini_attendance > backup.sql

# Restore PostgreSQL
psql -h <host> -U postgres mini_attendance < backup.sql

# Backup Redis
redis-cli -h <host> BGSAVE

# Restore Redis
redis-cli -h <host> --rdb /path/to/dump.rdb
```

### Scaling

```bash
# Scale Kubernetes deployment
kubectl scale deployment mini-attendance-api \
  --replicas=5 \
  -n mini-attendance

# Check scaling status
kubectl get deployment mini-attendance-api -n mini-attendance
```

### Updates & Rollback

```bash
# Update image
kubectl set image deployment/mini-attendance-api \
  api=<registry>/mini-attendance:v1.1.0 \
  -n mini-attendance

# Rollback to previous version
kubectl rollout undo deployment/mini-attendance-api \
  -n mini-attendance

# Check rollout status
kubectl rollout status deployment/mini-attendance-api \
  -n mini-attendance
```

## Troubleshooting

### Common Issues

**Database Connection Error**
```bash
# Check database connectivity
psql -h <host> -U postgres -d mini_attendance -c "SELECT 1"

# Check connection string
echo $DB_HOST $DB_PORT $DB_USER $DB_NAME
```

**Redis Connection Error**
```bash
# Check Redis connectivity
redis-cli -h <host> ping

# Check Redis memory
redis-cli -h <host> info memory
```

**Kafka Connection Error**
```bash
# Check Kafka brokers
kafka-broker-api-versions.sh --bootstrap-server <host>:9092

# Check topics
kafka-topics.sh --list --bootstrap-server <host>:9092
```

**High Memory Usage**
```bash
# Check container memory
docker stats

# Check database connections
psql -h <host> -U postgres -c "SELECT count(*) FROM pg_stat_activity"

# Check Redis memory
redis-cli -h <host> info memory
```

## Support & Documentation

- [README.md](./README.md) - Project overview
- [backend/README.md](./backend/README.md) - Backend documentation
- [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) - Architecture details
- [QUICKSTART.md](./QUICKSTART.md) - Quick start guide

---

**Last Updated**: January 2024
**Version**: 1.0.0
