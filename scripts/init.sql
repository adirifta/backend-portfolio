-- Database initialization script for Portfolio Backend

-- Create database
CREATE DATABASE portfolio_db;

-- Connect to the database
\c portfolio_db;

-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'admin',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create about table
CREATE TABLE abouts (
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

-- Create portfolio table
CREATE TABLE portfolios (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(500),
    project_url VARCHAR(500),
    category VARCHAR(100),
    tags VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create skills table
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    level INTEGER CHECK (level >= 0 AND level <= 100),
    category VARCHAR(100),
    icon VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create qualifications table
CREATE TABLE qualifications (
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
VALUES ('admin', '$2a$10$rRyBsGS4G2NKhL2H2/1NE.6a.TL1xE1JX5nU7QKJz5V5KJz5V5KJ', 'admin');

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
);

-- Insert sample portfolio items
INSERT INTO portfolios (title, description, image_url, project_url, category, tags)
VALUES 
(
    'E-commerce Website', 
    'A fully functional e-commerce platform with payment integration', 
    'https://example.com/project1.jpg', 
    'https://example.com/project1', 
    'Web Development', 
    'React, Node.js, PostgreSQL'
),
(
    'Mobile Fitness App', 
    'A fitness tracking application with workout plans and progress monitoring', 
    'https://example.com/project2.jpg', 
    'https://example.com/project2', 
    'Mobile Development', 
    'React Native, Firebase'
);

-- Insert sample skills
INSERT INTO skills (name, level, category, icon)
VALUES 
('JavaScript', 90, 'Programming', 'fa-js'),
('React', 85, 'Frontend', 'fa-react'),
('Node.js', 80, 'Backend', 'fa-node'),
('PostgreSQL', 75, 'Database', 'fa-database'),
('Go', 70, 'Programming', 'fa-code');

-- Insert sample qualifications
INSERT INTO qualifications (type, institution, title, description, start_date, end_date, current)
VALUES 
(
    'education', 
    'University of Technology', 
    'Bachelor of Computer Science', 
    'Specialized in software engineering and web development', 
    '2015-09-01', 
    '2019-06-01', 
    false
),
(
    'experience', 
    'Tech Solutions Inc.', 
    'Senior Full Stack Developer', 
    'Led development of multiple web applications and mentored junior developers', 
    '2020-01-15', 
    NULL, 
    true
),
(
    'experience', 
    'Web Innovations Ltd.', 
    'Frontend Developer', 
    'Developed user interfaces for various client projects using React', 
    '2019-07-01', 
    '2019-12-20', 
    false
);

-- Create indexes for better performance
CREATE INDEX idx_portfolios_category ON portfolios(category);
CREATE INDEX idx_skills_category ON skills(category);
CREATE INDEX idx_qualifications_type ON qualifications(type);
CREATE INDEX idx_qualifications_institution ON qualifications(institution);