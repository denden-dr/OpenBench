---
trigger: model_decision
description: Backend Folder Rules
---

# Common Guide
- DTO should never reference domain model / sensitive attribute
- Service and Repository object should always reference through interface
- API response should wrapped as object with message, status, data

# Naming Guide
- DTO files `something_dto`
- Service files `something_service`
- Repositoy files `something_repo`
- Handler files `something_handler`
- Mock files `mock_something`
- Test files `something_test`