-- INSERT INTO author(avatar_url,name,company_name,twitter_username) VALUES ('url', 'name', 'company', '@twitter');
CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    avatar_url text NOT NULL,
    name varchar(50) NOT NULL,
    company_name varchar(50) NOT NULL,
    twitter_username varchar(16)
);