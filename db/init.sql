-- Crear usuario que usará la API
CREATE USER userdb WITH PASSWORD 'admin';
GRANT ALL PRIVILEGES ON DATABASE tpwebdb TO userdb;

\c tpwebdb

-- Permisos sobre esquema
GRANT USAGE ON SCHEMA public TO userdb;

-- Defaultara futuros objetos
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT
  SELECT, INSERT, UPDATE, DELETE ON TABLES TO userdb;

ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT
  USAGE, SELECT, UPDATE ON SEQUENCES TO userdb;

-- Crear tabla (creará games_id_seq)
CREATE TABLE IF NOT EXISTS public.games (
    id SERIAL PRIMARY KEY,
    titulo       VARCHAR(150) NOT NULL,
    descripcion  VARCHAR(255) NOT NULL,
    categoria    VARCHAR(50) NOT NULL,
    fecha        DATE NOT NULL,
    estado       VARCHAR(20) CHECK (estado IN ('none', 'deseado', 'comprado')) NOT NULL,
    imagen       VARCHAR(50) NOT NULL,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- (Opcional) Transferir ownership para simplificar
ALTER TABLE public.games OWNER TO userdb;
ALTER SEQUENCE public.games_id_seq OWNER TO userdb;

-- Ahora sí: GRANT sobre TODO lo que ya existe
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO userdb;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO userdb;

-- Datos iniciales
INSERT INTO public.games (titulo, descripcion, categoria, fecha, estado, imagen) VALUES
('FIFA25','Simulador de Fútbol','Deporte','2024-09-10','comprado','img/fifa25.png'),
('Call of Duty','Juego de disparos','Accion','2021-06-21','none','img/cod.png'),
('FIFA26','Simulador de Fútbol','Deporte','2025-09-15','deseado','img/fifa26.png'),
('Battlefield 5','Juego de disparos','Accion','2023-06-21','none','img/btf5.png');