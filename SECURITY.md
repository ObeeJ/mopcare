# Security Guidelines

## Environment Variables

**CRITICAL:** Never commit `.env` files to version control!

### Setup Instructions:

1. Set environment variables in your deployment platform
2. Never commit `.env` files to version control
3. Use the deployment guide in `DEPLOYMENT.md`

### Required Environment Variables:

See `DEPLOYMENT.md` for complete list of required environment variables.

### Deployment:

- Use environment variables or secrets management in production
- Never hardcode credentials in source code
- Rotate keys regularly