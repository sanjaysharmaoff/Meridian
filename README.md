# System Architecture

1. **Ingress Translation:** API Gateway acts as the K8s Ingress, converting external REST/HTTP requests into typed binary gRPC calls for fast internal Go-to-Go communication.

2. **Database per Service:** Trip Service owns its MongoDB ride state. Other services do not access its database directly.

3. **Async Choreography:** Trip Service publishes to `find_available_drivers` and immediately responds, avoiding blocking while searching for a driver.

4. **Targeted Driver Routing:** Driver Service consumes the request, performs geospatial driver matching, and communicates via `driver_cmd_trip_request` and `driver_trip_response` queues.

5. **Webhook Payment Flow:** Stripe processes payment externally, calls the Payment Service webhook, which publishes `payment_success` to RabbitMQ to finalize the ride.

6. **DLX/DLQ Resilience:** Failed or unacknowledged messages are routed through RabbitMQ's **DLX/DLQ** for safe storage and later reprocessing.

7. **Container Scaling:** Docker/K8s enables independent service scaling, e.g. scale Driver Service from `2 → 20` pods without scaling Payment Service.

8. **Distributed Tracing:** Go `context` propagates trace IDs across **HTTP → gRPC → RabbitMQ → services**, while **Jaeger** visualizes the complete ride flow and helps identify latency and failures.

```mermaid
%%{init: {
  "theme": "dark",
  "themeVariables": {
    "background": "#070B14",
    "primaryTextColor": "#F8FAFC",
    "secondaryTextColor": "#E2E8F0",
    "lineColor": "#CBD5E1",
    "textColor": "#F8FAFC",
    "mainBkg": "#111827",
    "clusterBkg": "#0F172A",
    "clusterBorder": "#475569",
    "edgeLabelBackground": "#1E293B",
    "fontFamily": "Arial, sans-serif"
  }
}}%%

flowchart TD

    %% Styles
    classDef client fill:#6D28D9,color:#FFFFFF,stroke:#C4B5FD,stroke-width:3px
    classDef gateway fill:#1D4ED8,color:#FFFFFF,stroke:#93C5FD,stroke-width:3px
    classDef service fill:#047857,color:#FFFFFF,stroke:#6EE7B7,stroke-width:3px
    classDef queue fill:#C2410C,color:#FFFFFF,stroke:#FDBA74,stroke-width:3px
    classDef database fill:#0F766E,color:#FFFFFF,stroke:#5EEAD4,stroke-width:3px
    classDef tracing fill:#A16207,color:#FFFFFF,stroke:#FDE047,stroke-width:3px
    classDef external fill:#334155,color:#FFFFFF,stroke:#CBD5E1,stroke-width:3px

    %% External Actors
    Client["<b>CLIENT APPS</b><br/>Rider & Driver UI"]:::client
    Stripe["<b>STRIPE API</b><br/>External Checkout"]:::external

    %% Kubernetes Cluster
    subgraph K8s["KUBERNETES / DOCKER CLUSTER"]
        direction TB

        GW["<b>API GATEWAY</b><br/>Ingress Controller"]:::gateway

        %% Go Services
        subgraph Services["gRPC MICROSERVICES - GO"]
            direction LR

            Trip["<b>TRIP SERVICE</b><br/>Ride State Manager"]:::service
            Driver["<b>DRIVER SERVICE</b><br/>Geospatial Matching"]:::service
            Payment["<b>PAYMENT SERVICE</b><br/>Transaction Engine"]:::service
        end

        %% RabbitMQ
        subgraph Broker["RABBITMQ EVENT BUS"]
            direction TB

            Q1[/"<b>find_available_drivers</b><br/>Queue"/]:::queue
            Q2[/"<b>driver_cmd_trip_request</b><br/>Queue"/]:::queue
            Q3[/"<b>driver_trip_response</b><br/>Queue"/]:::queue
            Q4[/"<b>payment_success</b><br/>Queue"/]:::queue
            DLX[/"<b>DLX & DLQ</b><br/>Dead Letter Exchange"/]:::queue
        end

        %% Storage and Observability
        MongoDB[("<b>TRIP DB</b><br/>MongoDB")]:::database
        Jaeger["<b>JAEGER</b><br/>Distributed Tracing"]:::tracing
    end

    %% External Routing
    Client -->|"HTTP / REST"| GW
    Stripe -->|"WEBHOOKS"| Payment

    %% gRPC
    GW ==>|"gRPC CHANNEL"| Trip
    GW ==>|"gRPC CHANNEL"| Driver
    GW ==>|"gRPC CHANNEL"| Payment

    %% Event Choreography
    Trip -.->|"PUBLISHES"| Q1
    Q1 -.->|"CONSUMES"| Driver

    Driver -.->|"PUBLISHES"| Q2
    Driver -.->|"PUBLISHES"| Q3

    Q2 -.->|"CONSUMES"| Trip
    Q3 -.->|"CONSUMES"| Trip

    Payment -.->|"PUBLISHES"| Q4
    Q4 -.->|"CONSUMES"| Trip

    %% Dead Letter
    Q1 -.->|"FAILED MESSAGE"| DLX
    Q2 -.->|"FAILED MESSAGE"| DLX
    Q3 -.->|"FAILED MESSAGE"| DLX
    Q4 -.->|"FAILED MESSAGE"| DLX

    %% Database
    Trip ===|"DATABASE"| MongoDB

    %% Tracing
    GW -.->|"SPANS & CONTEXT"| Jaeger
    Trip -.->|"SPANS & CONTEXT"| Jaeger
    Driver -.->|"SPANS & CONTEXT"| Jaeger
    Payment -.->|"SPANS & CONTEXT"| Jaeger

    %% Link Colors
    linkStyle 2 stroke:#60A5FA,stroke-width:3px
    linkStyle 3 stroke:#60A5FA,stroke-width:3px
    linkStyle 4 stroke:#60A5FA,stroke-width:3px
```


<img width="2095" height="331" alt="image" src="https://github.com/user-attachments/assets/1b186404-19e3-4fe7-9b3b-1817a851c7e7" />



