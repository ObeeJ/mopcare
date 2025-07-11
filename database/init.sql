-- Database initialization script
-- Run this script in your Supabase SQL editor or PostgreSQL client

-- Drop existing tables if they exist (for clean setup)
DROP TABLE IF EXISTS user_course_enrollments CASCADE;
DROP TABLE IF EXISTS series CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS courses CASCADE;

-- Create tables with proper schema
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    overview_video_url TEXT,
    cover_image_url TEXT,
    unique_id VARCHAR(255) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE series (
    id SERIAL PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    video_url TEXT NOT NULL,
    thumbnail_url TEXT,
    duration INTEGER NOT NULL DEFAULT 0,
    is_free_preview BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    total_amount_paid DECIMAL(10,2) DEFAULT 0.00,
    enrolled_courses TEXT,
    completed_courses_count INTEGER DEFAULT 0,
    state VARCHAR(100),
    city VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE user_course_enrollments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('enrolled', 'completed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, course_id)
);

-- Create indexes
CREATE INDEX idx_series_course_id ON series(course_id);
CREATE INDEX idx_enrollments_user_id ON user_course_enrollments(user_id);
CREATE INDEX idx_enrollments_course_id ON user_course_enrollments(course_id);
CREATE INDEX idx_users_email ON users(email);