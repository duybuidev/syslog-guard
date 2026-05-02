```text
SysLog-Guard (Distributed System)
├── 0. Infrastructure Layer (Docker & Resource Limits)
├── 1. Data Generation Layer (Mock Ecosystem)
│   ├── Auth Service (Java - 128MB RAM Limit)
│   ├── Order Service (Go - 32MB RAM Limit)
│   └── Shipping Service (Python/FastAPI - 64MB RAM Limit)
├── 2. Core Logic Layer (SysWatch - Work In Progress)
│   └── Container Monitor & Alerting Engine (Go)
└── 3. Storage & Analytics (WIP)
    └── Redis + Elasticsearch + Kibana
```


TECHSTACK
1. GO
2. FastAPI
3. Python
4. Docker
