CREATE TABLE games (
    id SERIAL PRIMARY KEY,
    titulo       VARCHAR(150) NOT NULL,
    descripcion  VARCHAR(255) NOT NULL,
    categoria    VARCHAR(50) NOT NULL,
    fecha        DATE NOT NULL,
    estado       VARCHAR(20) CHECK (estado IN ('none', 'deseado', 'comprado')) NOT NULL,
    imagen      VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);