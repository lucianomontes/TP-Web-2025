PRE REQUISITOS:
-docker compose v2

1. Ejecutar el siguiente comando para levantar todos los servicios:

docker compose up --build -d

-El comando debería:
    -descargar imagenes si es necesario
    -levantar servicios (web y database)
    -correr script sql crear usuario con permisos y precargar la db

-Los logs se pueden ver con:

docker compose logs web
docker compose logs database

----------------------------------------------------------------------------------------------------------
2. Acceder a la web en el siguiente link:
http://localhost:8080
