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






-- CONSTRUCTORS

-- Returns the number of victories for a given constructor
CREATE OR REPLACE FUNCTION get_constructor_victories(constructor_name VARCHAR)
RETURNS INTEGER AS $$
DECLARE
    victories_count INTEGER;
BEGIN
    SELECT COUNT(*)
    INTO victories_count
    FROM RESULTADOS R
    JOIN CONSTRUTORES C ON R.idconstrutor = C.ID
    WHERE C.Nome = constructor_name
      AND R.Posicaofinal = 1;  -- Assuming "1" indicates a victory

    RETURN victories_count;
END;
$$ LANGUAGE plpgsql;

-- Returns the number of drivers that have driven for a given constructor
CREATE OR REPLACE FUNCTION get_unique_driver_count(constructor_name VARCHAR)
RETURNS INTEGER AS $$
DECLARE
    driver_count INTEGER;
BEGIN
    SELECT COUNT(DISTINCT R.nomepiloto) 
    INTO driver_count
    FROM RESULTADOS R
    JOIN CONSTRUTORES C ON R.idconstrutor = C.ID
    WHERE C.Nome = constructor_name;

    RETURN driver_count;
END;
$$ LANGUAGE plpgsql;

-- Returns the data range for a given constructor
CREATE OR REPLACE FUNCTION get_constructor_year_range(constructor_name VARCHAR)
RETURNS TABLE (first_year SMALLINT, last_year SMALLINT) AS $$
BEGIN
    RETURN QUERY
    SELECT MIN(CO.Ano) AS first_year, MAX(CO.Ano) AS last_year
    FROM RESULTADOS R
    JOIN CONSTRUTORES C ON R.IDConstrutor = C.ID
    JOIN CORRIDAS CO ON R.IDCorrida = CO.ID
    WHERE C.Nome = constructor_name;
END;
$$ LANGUAGE plpgsql;

-- DRIVERS

-- Returns the data range for a given driver
CREATE OR REPLACE FUNCTION get_driver_year_range(pilot_name VARCHAR)
RETURNS TABLE (first_year SMALLINT, last_year SMALLINT) AS $$
BEGIN
    RETURN QUERY
    SELECT MIN(CO.Ano) AS first_year, MAX(CO.Ano) AS last_year
    FROM RESULTADOS R
    JOIN PILOTOS P ON R.nomepiloto = P.nome
    JOIN CORRIDAS CO ON R.IDCorrida = CO.ID
    WHERE P.Nome = pilot_name;
END;
$$ LANGUAGE plpgsql;

-- Returns points and victories for a given driver by year
CREATE OR REPLACE FUNCTION get_driver_performance_by_year(driver_name VARCHAR)
RETURNS TABLE (
    year SMALLINT,
    total_points NUMERIC,
    total_victories INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        CO.Ano AS year,
        SUM(R.qtdpontos)::NUMERIC AS total_points,  -- Sum points for each year
        COUNT(CASE WHEN R.posicaoFinal = 1 THEN 1 END)::INT AS total_victories  -- Count victories (1st place)
    FROM 
        RESULTADOS R
    JOIN 
        PILOTOS P ON R.nomepiloto = P.Nome
    JOIN 
        CORRIDAS CO ON R.idcorrida = CO.ID
    WHERE 
        P.Nome = driver_name
    GROUP BY 
        CO.Ano
    ORDER BY 
        CO.Ano;
END;
$$ LANGUAGE plpgsql;

-- Returns points and victories for a given driver by circuit
CREATE OR REPLACE FUNCTION get_driver_performance_by_circuit(driver_name VARCHAR)
RETURNS TABLE (
    circuit_name VARCHAR,
    total_points NUMERIC,
    total_victories INT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        CI.Nome AS circuit_name,
        SUM(R.qtdpontos)::NUMERIC AS total_points,  -- Sum points for each circuit
        COUNT(CASE WHEN R.posicaoFinal = 1 THEN 1 END)::INT AS total_victories  -- Count victories (1st place)
    FROM 
        RESULTADOS R
    JOIN 
        PILOTOS P ON R.nomepiloto = P.Nome
    JOIN 
        CORRIDAS CO ON R.idcorrida = CO.ID
    JOIN 
        CIRCUITOS CI ON CO.nomecircuito = CI.Nome
    WHERE 
        P.Nome = driver_name
    GROUP BY 
        CI.Nome
    ORDER BY 
        CI.Nome;
END;
$$ LANGUAGE plpgsql;




-- Returns the number of times each status appears in the RESULTS table
CREATE OR REPLACE FUNCTION admin_report_status_counts()
RETURNS TABLE (
    status_name VARCHAR,
    status_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        descricaostatus AS status_name,
        COUNT(*) AS status_count
    FROM 
        RESULTADOS
    GROUP BY 
        descricaostatus
    ORDER BY 
        status_count DESC;
END;
$$ LANGUAGE plpgsql;

-- Create indexes
CREATE INDEX idx_resultados_constructor_pilot_position ON RESULTADOS (idconstrutor, nomepiloto, posicaofinal);
CREATE INDEX idx_construtores_nome ON CONSTRUTORES (Nome);


-- Returns each driver's name and the number of victories they have for a given constructor
CREATE OR REPLACE FUNCTION get_constructor_driver_wins(constructor_name VARCHAR)
RETURNS TABLE (
    driver_name VARCHAR,
    win_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        P.Nome AS driver_name,
        COUNT(CASE WHEN R.posicaofinal = 1 THEN 1 END) AS win_count
    FROM 
        PILOTOS P
    JOIN 
        RESULTADOS R ON P.Nome = R.nomepiloto
    JOIN 
        CONSTRUTORES C ON R.idconstrutor = C.ID
    WHERE 
        C.Nome = constructor_name
    GROUP BY 
        P.Nome
    ORDER BY 
        win_count DESC;
END;
$$ LANGUAGE plpgsql;

-- Returns the number of times each status appears in the RESULTS table for a given constructor
CREATE OR REPLACE FUNCTION get_constructor_status_count(constructor_name VARCHAR)
RETURNS TABLE (
    status_description VARCHAR,
    status_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        R.descricaostatus AS status_description,
        COUNT(*) AS status_count
    FROM 
        RESULTADOS R
    JOIN 
        CONSTRUTORES C ON R.idconstrutor = C.ID
    WHERE 
        C.Nome = constructor_name
    GROUP BY 
        R.descricaostatus
    ORDER BY 
        status_count DESC;
END;
$$ LANGUAGE plpgsql;



-- Create index
CREATE INDEX idx_resultados_pilot_victories ON RESULTADOS (nomepiloto, posicaofinal, idcorrida);

-- Returns a summary of the victories of a given driver
CREATE OR REPLACE FUNCTION get_driver_victories_summary(driver_name VARCHAR)
RETURNS TABLE (
    year NUMERIC,
    race_name VARCHAR,
    victory_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        EXTRACT(YEAR FROM Corrida.datacorrida) AS year,
        Corrida.Nome AS race_name,
        COUNT(*) AS victory_count
    FROM 
        RESULTADOS R
    JOIN 
        CORRIDAS Corrida ON R.idcorrida = Corrida.ID    
    WHERE 
        R.nomepiloto = driver_name
        AND R.posicaofinal = 1
    GROUP BY 
        ROLLUP(EXTRACT(YEAR FROM Corrida.datacorrida), Corrida.Nome)
    ORDER BY 
        year, race_name;
END;
$$ LANGUAGE plpgsql;


-- Create index
CREATE INDEX idx_resultados_pilot_status ON RESULTADOS (nomepiloto, descricaostatus);


-- Returns the number of times each status appears in the RESULTS table for a given driver
CREATE OR REPLACE FUNCTION get_driver_results_by_status(driver_name VARCHAR)
RETURNS TABLE (
    status_description VARCHAR,
    result_count BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        R.descricaostatus AS status_description,
        COUNT(*) AS result_count
    FROM 
        RESULTADOS R
    WHERE 
        R.nomepiloto = driver_name
    GROUP BY 
        R.descricaostatus
    ORDER BY 
        result_count DESC;
END;
$$ LANGUAGE plpgsql;


-- Create trigger to insert a new user in the USERS table when a new constructor is inserted
CREATE OR REPLACE FUNCTION trg_after_constructor_insert()
RETURNS TRIGGER AS $$
DECLARE
    existing_user_count INTEGER;
    new_login VARCHAR(255);
BEGIN
    -- Define the new login based on a standard format
    new_login := NEW.nome || '_c';
    
    -- Check if a user with this login already exists
    SELECT COUNT(*)
    INTO existing_user_count
    FROM USERS
    WHERE login = new_login;
    
    -- If the user already exists, raise an error and cancel the insertion
    IF existing_user_count > 0 THEN
        RAISE EXCEPTION 'User with this login already exists.';
    ELSE
        -- Insert the new user in the USERS table
        INSERT INTO USERS (login, password, tipo, id_original)
        VALUES (new_login, md5(NEW.nome), 'Escuderia', NEW.nome);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_after_constructor_insert
AFTER INSERT ON CONSTRUTORES
FOR EACH ROW
EXECUTE FUNCTION trg_after_constructor_insert();


-- Create trigger to insert a new user in the USERS table when a new driver is inserted
CREATE OR REPLACE FUNCTION trg_after_driver_insert()
RETURNS TRIGGER AS $$
DECLARE
    existing_user_count INTEGER;
    new_login VARCHAR(255);
BEGIN
    -- Define the new login based on a standard format
    new_login := NEW.nome || '_d';
    
    -- Check if a user with this login already exists
    SELECT COUNT(*)
    INTO existing_user_count
    FROM USERS
    WHERE USERS.login = new_login;
    
    -- If the user already exists, raise an error and cancel the insertion
    IF existing_user_count > 0 THEN
        RAISE EXCEPTION 'User with this login already exists.';
    ELSE
        -- Insert the new user in the USERS table with the driver's name as the password
        INSERT INTO USERS (login, password, tipo, id_original)
        VALUES (new_login, md5(NEW.nome), 'Piloto', NEW.nome);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_after_driver_insert
AFTER INSERT ON PILOTOS
FOR EACH ROW
EXECUTE FUNCTION trg_after_driver_insert();


-- Create index
CREATE INDEX idx_city_name_coordinates ON CIDADES (Nome, Lat, Long);
