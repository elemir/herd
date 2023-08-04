# Herd

Herd is an let another ECS framework for go language with Ebitengine integrations

# Design decisions

- Avoid dependency injection
- Use sparse sets as storage
- Use reflect for preparing a new query
- Use unsafe for access and add new entities and components on it