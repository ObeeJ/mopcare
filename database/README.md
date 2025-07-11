# Database Setup

## Quick Setup

1. **Connect to your Supabase project**
2. **Go to SQL Editor**
3. **Run the initialization script:**

```sql
-- Copy and paste the contents of init.sql
```

## Schema Files

- `schema.sql` - Complete database schema
- `init.sql` - Initialization script with DROP/CREATE statements

## Tables Created

- `courses` - Course information with unique_id support
- `series` - Video series within courses  
- `users` - User profiles with location data
- `user_course_enrollments` - Enrollment tracking

## Environment Variables Required

```env
SUPABASE_DB_URL=postgresql://postgres.[ref]:[password]@aws-0-us-east-1.pooler.supabase.com:6543/postgres
```