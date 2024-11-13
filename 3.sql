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

















-----




-- Create an index
CREATE INDEX idx_resultados_constructor_pilot_position ON RESULTADOS (idconstrutor, nomepiloto, posicaofinal);
CREATE INDEX idx_construtores_nome ON CONSTRUTORES (Nome);

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


----

CREATE INDEX idx_resultados_pilot_victories ON RESULTADOS (nomepiloto, posicaofinal, idcorrida);

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



CREATE INDEX idx_resultados_pilot_status ON RESULTADOS (nomepiloto, descricaostatus);

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


---

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




-- AEROPORTOS

CREATE INDEX idx_city_name_coordinates ON CIDADES (Nome, Lat, Long);
