-- init.sql (updated version)
-- Database initialization script for Portfolio Backend

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'admin',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create about table
CREATE TABLE IF NOT EXISTS abouts (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    title VARCHAR(255),
    description TEXT,
    email VARCHAR(255),
    phone VARCHAR(50),
    address TEXT,
    image_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create portfolio table (updated)
CREATE TABLE IF NOT EXISTS portfolios (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    tags VARCHAR(500),
    project_url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create portfolio_media table (new)
CREATE TABLE IF NOT EXISTS portfolio_media (
    id SERIAL PRIMARY KEY,
    portfolio_id INTEGER NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'image', 'video'
    url VARCHAR(500) NOT NULL,
    thumbnail VARCHAR(500),
    order_index INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (portfolio_id) REFERENCES portfolios(id) ON DELETE CASCADE
);

-- Create skills table
CREATE TABLE IF NOT EXISTS skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(255) NOT NULL,
    score INTEGER CHECK (score >= 0 AND score <= 100),
    category VARCHAR(100),
    icon VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create qualifications table
CREATE TABLE IF NOT EXISTS qualifications (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50) NOT NULL CHECK (type IN ('education', 'experience')),
    institution VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    current BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert initial admin user (password: admin123)
INSERT INTO users (username, password, role) 
VALUES ('admin', '$2a$10$rRyBsGS4G2NKhL2H2/1NE.6a.TL1xE1JX5nU7QKJz5V5KJz5V5KJ', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Insert sample about data
INSERT INTO abouts (name, title, description, email, phone, address, image_url)
VALUES (
    'John Doe', 
    'Full Stack Developer', 
    'Experienced developer with passion for creating innovative web applications.', 
    'john.doe@example.com', 
    '+1234567890', 
    '123 Main St, City, Country', 
    'https://example.com/profile.jpg'
)
ON CONFLICT (id) DO NOTHING;

-- Clean up old image_url column if exists
ALTER TABLE portfolios DROP COLUMN IF EXISTS image_url;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_portfolios_category ON portfolios(category);
CREATE INDEX IF NOT EXISTS idx_portfolio_media_portfolio_id ON portfolio_media(portfolio_id);
CREATE INDEX IF NOT EXISTS idx_skills_category ON skills(category);
CREATE INDEX IF NOT EXISTS idx_qualifications_type ON qualifications(type);
CREATE INDEX IF NOT EXISTS idx_qualifications_institution ON qualifications(institution);