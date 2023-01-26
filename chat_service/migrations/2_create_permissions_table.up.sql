CREATE TABLE IF NOT EXISTS permissions(
    id SERIAL PRIMARY KEY,
    user_type VARCHAR CHECK ("user_type" IN('superadmin', 'user')) NOT NULL,
    resource VARCHAR NOT NULL,
    action VARCHAR NOT NULL,
    UNIQUE(user_type, resource, action)   
);
