# Technical Debt Markdown Table Entry Template

| ID | Category | Title | Description / Context | Impact | Remediation Plan | Effort |  Status |
|----|----------|-------|-----------------------|--------|------------------|--------|--------|
| TD-001 | Database | Ephemeral Test Volumes | Postgres-test runs on ephemeral storage without volume mounts. | Test data is wiped on container down. Cannot debug persistent issues between runs. | Set up separate persistent test volume if state persistence becomes necessary. | Low | Active |
