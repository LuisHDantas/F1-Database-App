
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
