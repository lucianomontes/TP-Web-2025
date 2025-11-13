module tp_web_logica

go 1.22.2

require (
	github.com/lib/pq v1.10.9
	tp_web_datos v0.0.0
)

replace tp_web_datos => ../tp-web-datos
