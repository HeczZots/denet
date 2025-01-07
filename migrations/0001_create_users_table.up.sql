CREATE TABLE users (
    uid SERIAL PRIMARY KEY,        
    username VARCHAR(50) UNIQUE NOT NULL, 
    password_hash VARCHAR(255) NOT NULL, 
    refer_uid INTEGER DEFAULT 0 NOT NULL,              
    points INTEGER DEFAULT 0 NOT NULL      
);