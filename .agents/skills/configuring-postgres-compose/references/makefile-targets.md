# Makefile Targeted Teardown (Podman / Docker)

Standard `compose down` can fail or lock if network bridges are in use. Stop and force-remove containers by name before running down:

```make
test-integration:
	podman-compose -f docker-compose-test.yml up -d postgres-test
	@until [ "$$(podman inspect --format='{{.State.Health.Status}}' openbench-postgres-test 2>/dev/null)" = "healthy" ]; do sleep 1; done
	cd apps/backend && APP_ENV=test go test -count=1 -tags=integration ./...
	podman stop openbench-postgres-test || true
	podman rm -f -v openbench-postgres-test || true

clean:
	@podman stop openbench-postgres-dev openbench-postgres-test || true
	@podman rm -f -v openbench-postgres-dev openbench-postgres-test || true
	podman-compose --env-file apps/backend/.env down || true
	podman-compose -f docker-compose-test.yml down || true
```
