-- INSERT INTO author(avatar_url,name,company_name,twitter_username) VALUES ('url', 'name', 'company', '@twitter');
CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    avatar_url text,
    name varchar(50) UNIQUE NOT NULL,
    company_name varchar(50),
    twitter_username varchar(15)
);