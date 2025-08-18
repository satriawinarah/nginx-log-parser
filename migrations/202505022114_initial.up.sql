create table visitor (
    id UUID primary key,
    host_name VARCHAR(45) not null,
    request_ip VARCHAR(255) not null,
    request_ip_location VARCHAR(255) not null,
    request_time TIMESTAMP DEFAULT now() not null,
    request_method TEXT not null,
    request_uri TEXT not null,
    user_agent TEXT not null,
    response_status VARCHAR(100) not null,
    created_by VARCHAR(255) not null,
    updated_by VARCHAR(255),
    created_at TIMESTAMP DEFAULT now() not null,
    updated_at TIMESTAMP DEFAULT now(),
    is_deleted BOOLEAN DEFAULT false not null
);