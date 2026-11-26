# DataHub API
The API's main feature is to receive messages from message brokers and redirect them to a proper external service.
For achieving this, the API allows users to connect their services, informing the event types they want to receive and to which service each message should be sent.  
Internally we connect dynamically to distinct message brokers, associating the user to their respective events, and control the encryption keys all through the API.

From the necessity of data security regulations(LGPD and PCI), the data sent through the application should always be encrypted, so we provide resources for creating a key per customer, and for exchanging based on RSA and AES encryption algorithms. 

## Architecture

### Overview for API and connections 

![Simplified Architecture](resources/public/SimplifiedArchitecture.png)

The image displays an overview on how the API controls and keep connections alive, it does not reflect directly on how the code is structured. Through the API we can add new connections, and trigger commands to connect or disconnect to each or all of them. Also, on the API is where on the users are set up internally, and there are public resources (with required authentication) where the available events can be linked to their own APIs.

The cluster control is the part of the application which will maintain all nodes equalized with all enabled connections. When creating a connection through the API, it will broadcast a message through the application's internal RabbitMQ so all nodes can recognize the new data and connect to it.

Connecting to a broker, may cause issues related to disconnecting due to network problems or broker maintenance, some others related to guaranteeing a message will be successfully read. They are handled in the AMQP Connection Management layer.

### Message delivery

The application's core feature is to connect to an AMQP broker then read, encrypt and deliver its messages. The AMQP service orchestrates this whole behavior, the image below displays how it is supposed to work. 

![AMQP Service](resources/public/AMQPService.png)

Allowing creating, reading, updating and deleting connections, the **AMQP Service** main features are handling connection management and accessing data from a repository. The **connection manager** is simply a map that stores active connections, its structure is represented by a `map[string]Connector`.

A **Connector** is an interface responsible for connecting, disconnecting, automatic reconnection, reading messages, handling acknowledgements. The business logic for reading a message is handled by another interface called **Message Processor**. For delivering WebHooks, there's the need for creating a standardized structure, and obtaining some of its required data can vary depending on the message's source, the **Message Processor** interface manipulates message for many sources, and it's final task is to send a message to the WebHook Deliverer. 

The **WebHook Deliverer** is a structure for handling a internal RabbitMQ queue, receiving a specific DTO. It has the responsibility of performing HTTP calls. If a message is not delivered correctly due to HTTP errors, it should be rescheduled, with an exponential delay based on how many attempts were made.

## Infrastructure


For creating a new instance of this application it's necessary to setup a PostgresSQL and a RabbitMQ Server.  
The applications relies on a API Management tool to handle the consumer certificates, such as Azure's API Management.
On the API Manager's side will be handled both authentication and which resources will be available to the issuers. On the API's side, only the issuer's identifier is needed for authentication, usually passed on HTTP Headers.  
Remember to update the following properties:

```yaml
swagger:
  ip: localhost
  port: 8080

splunk:
  dial: "localhost:4455"

rabbit:
  dial: "amqp://guest:guest@127.0.0.1:5672"
  virtualHost: "hooks"

database:
  dialect: "postgres"
  connectionUri: "host=127.0.0.1 port=5432 dbname=pier_hooks sslmode=enable"
  user: "postgres"
  password: "adminpass"
```

