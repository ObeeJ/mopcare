# Deployment Guide

## Environment Variables Setup

### Required Environment Variables:
Set these in your deployment platform (Render/Heroku/etc):

```
SUPABASE_DB_URL=postgresql://postgres.[ref]:[password]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
SUPABASE_PROJECT_URL=https://[project-id].supabase.co
SUPABASE_SERVICE_KEY=[your-service-key]
GATEWAY_PORT=9090
```

### Database Setup:
1. Run the SQL script from `database/init.sql` in your Supabase SQL editor
2. This creates all required tables and indexes

### Render Deployment:
1. Connect your GitHub repository to Render
2. Set environment variables in Render dashboard
3. Deploy using `render.yaml` configuration

### Docker Deployment:
```bash
docker-compose up --build -d
```

## Security Notes:
- Never commit `.env` files to version control
- Use environment variables in production
- Rotate database credentials regularly