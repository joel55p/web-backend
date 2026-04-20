package main // El paquete "main" es el punto de entrada de la aplicación en Go. Es necesario para ejecutar el programa.

import ( // Este es un servidor HTTP básico en Go que escucha en el puerto 8080 y responde con "Hello World" a cualquier solicitud.

	"database/sql" // "database/sql" se utiliza para interactuar con una base de datos SQL, aunque en este código no se muestra su uso específico, es común en aplicaciones web para manejar datos persistentes.

	"log"          // "log" se utiliza para registrar errores y mensajes informativos en la consola.
	"net"          // "net" se utiliza para crear un servidor TCP que escuche en el puerto 8080 y acepte conexiones entrantes. uno de os paquetes más importantes para la comunicación de red en Go.
	_ "modernc.org/sqlite" // Este es un import anónimo que se utiliza para registrar el controlador de SQLite con el paquete "database/sql". Esto permite que el programa utilice SQLite como base de datos sin necesidad de importar explícitamente el paquete en el código. El guion bajo (_) indica que el paquete se importa solo por sus efectos secundarios, es decir, para registrar el controlador de la base de datos, sin utilizar directamente ninguna función o tipo del paquete en el código.
)

func main() { // La función "main" es el punto de entrada del programa. Aquí se configura el servidor TCP y se maneja la lógica principal para aceptar conexiones.
	// se necesita hacer primero una conexion a la base de datos
	db, err := sql.Open("sqlite", "series.db") //se abre conexion una vez en el main, se hace con sql.open que recibe el tipo de base de datos y el nombre del archivo
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()                              //se cierra cuando se termine el main, es decir, cuando se cierre el servidor, se cierra la conexion a la base de datos para liberar recursos. es importante cerrar la conexion a la base de datos para evitar fugas de memoria y asegurar que los recursos se liberen adecuadamente cuando el programa termine su ejecución.

	listener, err := net.Listen("tcp", ":8080") // "net.Listen" crea un servidor TCP que escucha en el puerto 8080. Si hay un error al crear el servidor, se registra y se termina el programa. si funciona se guarda el listener para aceptar conexiones entrantes.
	if err != nil {                             // Si ocurre un error al intentar escuchar en el puerto, se registra el error y se termina el programa.
		log.Fatal(err) // "log.Fatal" registra el error y termina la ejecución del programa. Esto es útil para asegurarse de que el servidor no intente continuar si no puede escuchar en el puerto especificado.
	}
	defer listener.Close()              // "defer" asegura que el listener se cerrará cuando la función "main" termine, lo que es importante para liberar recursos. y main va a esperar a que se cierre el listener antes de finalizar la ejecución del programa.
	log.Print("Listening on port 8080") // Se registra un mensaje en la consola indicando que el servidor está escuchando en el puerto 8080.

	for { // Este es un bucle infinito que acepta conexiones entrantes. Cada vez que se acepta una conexión, se maneja en una goroutine separada para permitir que el servidor continúe aceptando otras conexiones mientras se procesa la actual.
		conn, err := listener.Accept() // "listener.Accept" espera a que llegue una conexión entrante y la acepta. Si hay un error al aceptar la conexión, se registra el error y se continúa con el siguiente ciclo del bucle para esperar otra conexión.
		if err != nil {                // Si ocurre un error al aceptar la conexión, se registra el error y se continúa con el siguiente ciclo del bucle para esperar otra conexión.
			log.Print("Error accepting:", err) // "log.Print" registra el error que ocurrió al aceptar la conexión, pero no termina el programa. Esto permite que el servidor siga funcionando y acepte otras conexiones a pesar de los errores ocasionales.
			continue                           // "continue" se utiliza para saltar el resto del código en el bucle actual y pasar a la siguiente iteración, lo que permite que el servidor siga aceptando conexiones incluso si ocurre un error al aceptar una conexión específica.
		}
		go handle(conn, db) // "go handle(conn)" inicia una nueva goroutine para manejar la conexión aceptada. Esto permite que el servidor procese múltiples conexiones simultáneamente sin bloquear el bucle principal que acepta conexiones. Cada conexión se maneja de forma independiente en su propia goroutine, lo que mejora la capacidad de respuesta del servidor.
	}
} //cabe mencionar que conn, err := listener.Accept() lo que hace es que si acepta se guarda en la variable conn, y si no acepta se guarda el error en la variable err, y se maneja el error con un if para evitar que el programa se caiga.

// osea que if err !=nil lo que significa es que si si hay un error, entonces se ejecuta el bloque de código dentro del if, que en este caso es log.Print("Error accepting:", err) y continue para seguir aceptando conexiones. Si no hay error, entonces se ejecuta go handle(conn) para manejar la conexión aceptada en una goroutine separada.
// una gorutine es una función que se ejecuta de manera concurrente con otras funciones. Es una forma ligera de manejar múltiples tareas al mismo tiempo sin bloquear el programa principal. En este caso, cada vez que se acepta una conexión, se inicia una nueva goroutine para manejar esa conexión específica, lo que permite que el servidor siga aceptando otras conexiones mientras se procesa la actual.

	

	

	
	



