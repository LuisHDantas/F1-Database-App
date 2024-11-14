-- Create the USERS table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password TEXT NOT NULL, 
    tipo VARCHAR(20) CHECK (tipo IN ('Administrador', 'Escuderia', 'Piloto')),
    id_original VARCHAR(100)  
);

-- Log Table to track login activity
CREATE TABLE log_table (
    log_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
    login_date DATE DEFAULT CURRENT_DATE,
    login_time TIME DEFAULT CURRENT_TIME
);

-- Add existing data to the USERS table
INSERT INTO users (login, password, tipo, id_original)
SELECT 
    pilotos.Nome || '_d' AS login,
    MD5(pilotos.Nome) AS password,
    'Piloto' AS tipo,
    pilotos.Nome AS id_original
FROM pilotos;

INSERT INTO users (login, password, tipo, id_original)
SELECT 
    construtores.Nome || '_c' AS login,
    MD5(construtores.Nome) AS password,
    'Escuderia' AS tipo,
    construtores.Nome AS id_original
FROM construtores;

INSERT INTO users (login, password, tipo, id_original)
VALUES ('admin', MD5('admin'), 'Administrador', 'Admin');
