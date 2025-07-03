# Neureka

Eureka client library written in Go for `Blockman Forge` game

## Object Definitions

These are the key structs used by the library:

- **`AddressObject`**: Represents the host information of an application instance.
- **`AppObject`**: Represents application information (simplified version of the data pulled from Eureka, stored in cache).
- **`EurekaAppCache`**: A cache holding a list of applications and their available hosts for fast access.
- **`EurekaClientConfig`**: Client configuration struct — contains everything needed to connect to and interact with the Eureka server.
- **`AppResponse`**: Data structure representing the response returned by Eureka when requesting app data.
- **`EurekaAppInfo`**: Structure representing an application’s metadata when fetched from Eureka.
- **`EurekaAppInstance`**: Represents a full Eureka application instance.
- **`EurekaRegisterRequest`**: Data structure used when sending a registration request to Eureka.

## Implemented Methods

- **`EurekaRegist`**: Register your service/application with Eureka.
- **`EurekaHeartBeat`**: Send a heartbeat to renew the registration lease.
- **`EurekaGetApp`**: Retrieve the current state of an application from Eureka.
- **`EurekaDeleteApp`**: Remove/unregister an application instance from Eureka.

## Logic Flow

1. Initialize `EurekaClientConfig` with settings:
   - `EurekaServerAddress`: Address of the Eureka server.
   - `Authorization`: HTTP Auth header for the Eureka server (e.g., `Basic KWJDhaDAWIDhndwa=`).
   - `AppName`: Name of your Go service.
   - `InstanceDomain`: Domain name of your service instance. If empty, defaults to IP.
   - `InstanceIp`: The IP address of the service instance (external IP if applicable).
   - `InstancePort`: The port your Go service runs on.
   - `InstanceHealthCheckUrl`: Health check endpoint.
   - `RenewalIntervalInSecs`: Time in seconds for lease renewal.
   - `DurationInSecs`: How long Eureka retains the instance info after lease expiry.
   - `AppRefreshSecs`: How often (in seconds) to refresh the cache of known apps.

2. Clean up previous instance (optional but recommended):
   - Call `EurekaGetApp` to check if already registered.
   - Call `EurekaDeleteApp` to remove stale registration.

3. Register current instance using `EurekaRegist`.

4. Start background heartbeat to renew registration via `EurekaHeartBeat`.

5. Maintain a background cache of other registered applications using `EurekaAppCache`.

6. Use `GetAppUrl` to retrieve service addresses by app name.

## Quick Start

### Automatic Flow

Use the built-in startup method `Start()` to handle:
- Registration
- Heartbeat
- Service discovery cache

### Steps:

1. Prepare and populate `EurekaClientConfig`
2. Call `Start()` for single Eureka service instance
3. Or call `StartBatch()` if dealing with multiple services
