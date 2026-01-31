# Deployment Guide

Guide for deploying the Finance App backend to various platforms.

## Table of Contents

1. [Docker](#docker)
2. [Fly.io](#flyio)
3. [Render](#render)
4. [Local VPS (Linux)](#local-vps)

---

## Docker

### Build Docker Image

```bash
# Build the image
docker build -t finance-app:latest .

# Tag for registry (optional)
docker tag finance-app:latest huxxnainali/finance-app:latest
```

### Run Locally with Docker

```bash
# Create .env file with your settings
cp .env.example .env

# Run with Docker Compose (recommended)
docker-compose up -d

# Or run standalone container
docker run -p 3000:3000 \
  --env-file .env \
  finance-app:latest
```

### Push to Docker Registry

```bash
# Login to Docker Hub
docker login

# Push image
docker push huxxnainali/finance-app:latest
```

---

## Fly.io

### Prerequisites

- Fly.io account
- `flyctl` CLI installed

### Setup & Deploy

#### 1. Initialize Fly App

```bash
flyctl launch
```

During initialization:

- Choose app name (e.g., `finance-app`)
- Choose region (closest to your users)
- Add MongoDB external service or Fly Postgres

#### 2. Configure Secrets

```bash
flyctl secrets set \
  MONGODB_URI="your-mongodb-uri" \
  JWT_SECRET="your-secret-key" \
  DATABASE_NAME="finance_app" \
  JWT_EXPIRY_HOURS="24"
```

#### 3. Configure fly.toml

```toml
app = "finance-app"
primary_region = "lax" # Or your region

[build]
  image = "finance-app:latest"
  dockerfile = "./Dockerfile"

[env]
  PORT = "3000"

[[services]]
  protocol = "tcp"
  internal_port = 3000
  processes = ["app"]

  [services.http_checks]
    enabled = true
    uri = "/health"

  [[services.ports]]
    port = 80
    handlers = ["http"]
    force_https = true

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]
```

#### 4. Deploy

```bash
flyctl deploy
```

#### 5. Monitor

```bash
# View logs
flyctl logs

# Check app status
flyctl status

# SSH into running instance
flyctl ssh console
```

### Database on Fly.io

Option A: External MongoDB (MongoDB Atlas recommended)

- Create free tier on MongoDB Atlas
- Use connection string in `MONGODB_URI`

Option B: Postgres (if switching from MongoDB)

- `flyctl postgres create` (requires refactoring)

---

## Render

### Prerequisites

- Render account (free or paid)
- GitHub repository with code

### Setup & Deploy

#### 1. Connect GitHub

1. Go to Render dashboard
2. Click "New Web Service"
3. Select "Build and deploy from a Git repository"
4. Connect your GitHub repository

#### 2. Configure Service

```
Name: finance-app
Environment: Docker
Build Command: (left blank, uses Dockerfile)
Start Command: (left blank, uses EXPOSE 3000)
```

#### 3. Add Environment Variables

```
MONGODB_URI=your-mongodb-uri
JWT_SECRET=your-secret-key
DATABASE_NAME=finance_app
JWT_EXPIRY_HOURS=24
PORT=3000
```

#### 4. Create Database

Option A: MongoDB Atlas (recommended)

1. Create free tier cluster on MongoDB Atlas
2. Get connection string
3. Add to environment variables

Option B: Render PostgreSQL

```
Name: finance-app-db
PostgreSQL version: 15
```

#### 5. Deploy

- Push code to GitHub
- Render automatically builds and deploys
- Monitor deployment in Render dashboard

#### 6. Custom Domain (optional)

1. Go to "Settings" → "Custom Domain"
2. Add your domain
3. Update DNS records according to Render instructions

---

## Local VPS

### Prerequisites

- Linux server (Ubuntu 20.04+ recommended)
- SSH access
- 2GB+ RAM
- Public IP or domain

### Setup

#### 1. Install Dependencies

```bash
# SSH into server
ssh user@your-server-ip

# Update system
sudo apt update && sudo apt upgrade -y

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install MongoDB
curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | sudo apt-key add -
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list
sudo apt update
sudo apt install -y mongodb-org
sudo systemctl start mongod
sudo systemctl enable mongod

# Install Nginx (reverse proxy)
sudo apt install -y nginx

# Install Certbot (SSL certificates)
sudo apt install -y certbot python3-certbot-nginx
```

#### 2. Clone & Build Application

```bash
# Clone repository
git clone <your-repo> /opt/finance-app
cd /opt/finance-app

# Create .env
cp .env.example .env
# Edit .env with your settings

# Build
go build -o finance-app cmd/server/main.go
```

#### 3. Create Systemd Service

```bash
sudo nano /etc/systemd/system/finance-app.service
```

Paste:

```ini
[Unit]
Description=Finance App Backend
After=network.target mongodb.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/finance-app
EnvironmentFile=/opt/finance-app/.env
ExecStart=/opt/finance-app/finance-app
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable finance-app
sudo systemctl start finance-app

# Check status
sudo systemctl status finance-app
```

#### 4. Configure Nginx

```bash
sudo nano /etc/nginx/sites-available/finance-app
```

Paste:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/finance-app /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

#### 5. SSL Certificate

```bash
sudo certbot --nginx -d your-domain.com
```

#### 6. Monitoring & Logs

```bash
# Check app logs
sudo journalctl -u finance-app -f

# Check Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log

# Monitor system
top
htop  # Install if needed: sudo apt install htop
```

---

## Pre-Deployment Checklist

- [ ] All tests passing
- [ ] Environment variables configured
- [ ] MongoDB connection tested
- [ ] JWT secret changed from default
- [ ] .env not committed to Git
- [ ] CORS configured if needed
- [ ] Database indexes created
- [ ] Error handling complete
- [ ] Rate limiting configured (optional)
- [ ] Health check endpoint working

---

## Post-Deployment

### Verify Deployment

```bash
# Check health endpoint
curl https://your-domain.com/health

# Test signup
curl -X POST https://your-domain.com/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'
```

### Monitor

- Set up application monitoring (DataDog, New Relic, etc.)
- Enable database backups
- Set up log aggregation
- Configure alerts for errors

### Scaling

- Horizontal scaling: Deploy multiple instances behind load balancer
- Vertical scaling: Increase server resources
- Database: Connection pooling, replication for MongoDB

---

## Troubleshooting

### Common Issues

**Port Already in Use**

```bash
# Find process using port 3000
lsof -i :3000
# Kill it
kill -9 <PID>
```

**MongoDB Connection Failed**

```bash
# Check MongoDB status
sudo systemctl status mongod
# Restart if needed
sudo systemctl restart mongod
```

**Permission Denied Errors**

```bash
# Fix file permissions
sudo chown -R user:user /opt/finance-app
chmod -R 755 /opt/finance-app
```

**Out of Memory**

```bash
# Check memory usage
free -h
# Increase swap if needed
sudo fallocate -l 4G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

---

## Security Hardening

1. **Update Dependencies**

   ```bash
   go get -u ./...
   ```

2. **Firewall Rules**

   ```bash
   sudo ufw allow 22/tcp
   sudo ufw allow 80/tcp
   sudo ufw allow 443/tcp
   sudo ufw enable
   ```

3. **SSH Security**
   - Disable password authentication
   - Use SSH keys only
   - Limit login attempts

4. **Database Security**
   - Use MongoDB authentication
   - Restrict network access
   - Enable encryption at rest

5. **Application Security**
   - Keep dependencies updated
   - Rotate secrets regularly
   - Implement rate limiting
   - Add request validation

---

## Backup & Recovery

### MongoDB Backup

```bash
# Full backup
mongodump --out /backup/$(date +%Y%m%d_%H%M%S)

# Restore
mongorestore /backup/20240131_120000
```

### Application Backup

```bash
# Backup source code
tar -czf finance-app-backup.tar.gz /opt/finance-app

# Upload to storage (S3, etc.)
```

---

## Further Reading

- [Fly.io Docs](https://fly.io/docs/)
- [Render Docs](https://render.com/docs)
- [MongoDB Atlas Docs](https://www.mongodb.com/docs/atlas/)
- [Nginx Docs](https://nginx.org/en/docs/)
- [Let's Encrypt](https://letsencrypt.org/)
