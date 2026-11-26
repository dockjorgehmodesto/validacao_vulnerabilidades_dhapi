# Available Environment Configuration

## Feature toggling
* `features.disableMessageReading [bool]` Disables the application's ability to read messages for processing
* `features.disableWebhookDelivery [bool]` Disables the application's ability to deliver webhooks

## Application's descriptions and Exposing
* `app.profile [string]` Application's profile. Set to "prod" when going live.
* `app.name [string]` Application's name.
* `app.description [string]` Application's description. Used on swagger docs. 
* `server.port [string]` Application's internal port.
* `swagger.ip [string]` Application's external address.
* `swagger.port [numeric]` Application's external port.
* `swagger.customJson.enabled [bool]` Enables routing swagger's json to a customized URL.
* `swagger.customJson.url [string]` URL for swagger's json.
* `swagger.hideInternals [string]` Hides all resources tagged with "Internal".
* `swagger.internalsToken [string]` Token for displaying internal resources.

## Logs and Splunk
* `log.level [string]` Type of log for printing in the console.
* `log.splunk [bool]` Toggles splunk integration.
* `log.stdOut.enable [bool]` Pretty print logs sent to splunk to stdout.
* `splunk.dial [string]` Splunk's host and port for TCP connection.
* `log.billing.enabled [bool]` Toggles integration with billing's middleware. This will make the middleware inactive, thus stops sending messages to logstash.
* `log.billing.tcp.enabled [bool]` Toggles delivery to TCP.
* `log.billing.tcp.connectionstring [string]` Target TCP connection with host and port for delivering billing's logs.  

## Database, Pagination and Caching
* `database.showSQL [bool]` Toggles script logging. 
* `database.dialect [string]` Type of database dialect.
* `database.connectionUri [string]` Connection string for the datasource. Eg.: "host=127.0.0.1 port=5432 dbname=pier_hooks sslmode=disable".
* `database.user [string]` User for database logon.
* `database.password [string]` Password for database logon. 
* `database.migrationPrincipal.enabled [bool]` Toggles the database migration behavior for the application's database.
* `database.migrationPrincipal.folder [string]` Directory containing all SQL scripts. 
* `database.migrationPrincipal.minimumIdle [numeric]` Amount of connections that should be always open.
* `database.migrationPrincipal.maximumPoolSize [numeric]` Connection limit.
* `database.migrationIssuers.enabled [bool]` Toggles the database migration behavior for all issuers' database.
* `database.migrationIssuers.folder [string]` Directory containing all SQL scripts.
* `database.migrationIssuers.minimumIdle [numeric]` Amount of connections that should be always open.
* `database.migrationIssuers.maximumPoolSize [numeric]` Connection limit.
* `pagination.pageDefault [numeric]` Default page index. 
* `pagination.limitDefault [numeric]` Default limit per page.
* `pagination.sortDefault [string]` Default sorting for database queries.
* `cache.debug [bool]` Toggles console logs for cache access.
* `cache.enabled [bool]` Main toggle for caching.
* `cache.entities.tenant [bool]` Toggles caching for Tenants.

## Core RabbitMQ

Controls for the application's core RabbitMQ connection.

* `rabbit.dial [string]` AMQP Connection string
* `rabbit.virtualHost [string]` Virtual host for the application
* `rabbit.sleep [numeric]` Time to check connection and reconnect
* `rabbit.headers.attempts [numeric]` Max attempts for sending a message

## Dynamic AMQP connections

### Application's configuration
* `amqp.reconnectTimeoutSeconds [numeric]` Timeout in seconds for reconnecting, currently used only by ServiceBus connections. 
* `amqp.dynamicRabbit.restartTimeoutSeconds [numeric]` Timeout in seconds for reconnecting to RabbitMQ Queues.

### Parallelism handling
Controls parallelism when fetching data from a queue or subscription.  
Each listening connection creates *N* goroutines based on the values. 
*Be careful* when increasing this value, as it can overload database connections, http calls(requests can be canceled if too many are open) or even connections with the AMQP broker itself. Lowering it can slow message reads.

* `amqp.workerCount.rabbitMqConnector [numeric]` Amount of goroutines used by each RabbitMQ connector.  
* `amqp.workerCount.servicebusConnector [numeric]` Amount of goroutines used by each ServiceBus connector.
* `amqp.workerCount.webhookDeliverer [numeric]` Amount of goroutines used by the webhook deliverer.

### ServiceBus properties
* `amqp.deadletter.contextTimeoutInMinutes [numeric]` Timeout once used for reading data from a dead letter queue, will soon be removed.
* `amqp.deadletter.workerCount [numeric]` Similar to previous topic, but in the dead letter context, will soon be removed.

### RabbitMQ internal properties
* `amqp.dynamicRabbit.defaultHeartbeatSeconds [numeric]`  RabbitMQ's internal property for health checking.
* `amqp.dynamicRabbit.defaultLocale [string]` RabbitMQ's internal property for localization.
* `amqp.dynamicRabbit.exchangesDeclaration.kind [string]` Default exchange type when declaring or connecting to a new one.
* `amqp.dynamicRabbit.exchangesDeclaration.durable [bool]` If false, the exchange is deleted when nothing is connected to it, this can mean data loss.
* `amqp.dynamicRabbit.exchangesDeclaration.autoDelete [bool]` If true, the exchange is deleted when there isn't any messages in it.
* `amqp.dynamicRabbit.exchangesDeclaration.internal [bool]` If true, only RabbitMQ itself can communicate with the exchange.
* `amqp.dynamicRabbit.exchangesDeclaration.noWait [bool]` If true, the server will not acknowledge interactions.
* `amqp.dynamicRabbit.queuesDeclaration.durable [bool]` Same behavior presented on durable exchanges, but for queues.
* `amqp.dynamicRabbit.queuesDeclaration.autoDelete [bool]` Same behavior presented on durable exchanges, but for queues.
* `amqp.dynamicRabbit.queuesDeclaration.exclusive [bool]` If true, only one connection is allowed on the queue.
* `amqp.dynamicRabbit.queuesDeclaration.noWait [bool]` If true, the server will not acknowledge interactions.
* `amqp.dynamicRabbit.queueBindings.noWait [bool]` If true, the server will not acknowledge interactions.
* `amqp.dynamicRabbit.consumer.autoAck [bool]` When true, the server will not wait for clients validation on message reading.
* `amqp.dynamicRabbit.consumer.exclusive [bool]` If true, only one connection is allowed to consume messages.
* `amqp.dynamicRabbit.consumer.noLocal [bool]` If true, the host that sent a message cannot read it.
* `amqp.dynamicRabbit.consumer.noWait [bool]` If true, the server will not acknowledge interactions.

## HTTP Client

Properties for creating an HTTP client. They are setup with the default values used by Golang's HTTP client. They can be changed for handling either parallelism related issues, or limiting timeouts.

Some properties were not described, once they were exported to environment variables just for flexibility, they are also following the default values for a HttpClient.
 

* `http.timeoutSecs [numeric]` Timeout in seconds for a http connection. Overrides the timeouts bellow.
* `http.transport.dialContext.timeoutSecs [numeric]` -
* `http.transport.dialContext.keepAliveSecs [numeric]` - 
* `http.transport.maxIdleconns [numeric]` Limit for parallel open http connections. This value is application-wide, as in, it is not client or listener specific.
* `http.transport.idleConnTimeoutSecs [numeric]` -
* `http.transport.tlsHandshakeTimeoutSecs [numeric]` -
* `http.transport.expectContinueTimeoutSecs [numeric]` -

### Misc
* `encryption.aesKey.monthsToExpire [numeric]` Default amount in months to expire each issuer AES Key
