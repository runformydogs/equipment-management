CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'viewer',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS network_nodes (
                                             id SERIAL PRIMARY KEY,
                                             name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id INT REFERENCES network_nodes(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS devices (
                                       id SERIAL PRIMARY KEY,
                                       type VARCHAR(100) NOT NULL,
    vendor VARCHAR(100),
    model VARCHAR(100),
    serial VARCHAR(100) UNIQUE NOT NULL,
    location VARCHAR(255),
    status VARCHAR(50) DEFAULT 'active',
    network_node_id INT REFERENCES network_nodes(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );