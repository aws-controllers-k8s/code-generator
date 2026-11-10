# Suite configuration for the {{ .Service }}-controller Chainsaw tests. test-infra passes
# this file to `chainsaw test`; JUnit reporting is set on the command line.
#
# ack-generate chainsaw wrote this file once, as a starting point. It's this
# controller's to edit, and isn't regenerated: tune the timeouts to the
# service's slowest resources, and parallel to the test account's quotas.
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: {{ .Service }}-controller
spec:
  timeouts:
    apply: 30s
    # Upper bounds, to catch a reconcile that never converges: a resource that
    # converges finishes long before them, and a CR that goes ACK.Terminal
    # fails its step at once. assert bounds the waits on a CR's conditions,
    # so it must outlast the slowest resource's create or update.
    assert: 5m
    # Deleting a test's namespace waits for the controller to delete every
    # AWS resource in it.
    cleanup: 20m
    delete: 15m
    # `error` returns as soon as the resource doesn't match, so this only
    # bounds how long a matching (failing) check is retried.
    error: 10s
    # Scripts: AWS reads run once, but waiting for an AWS resource to be gone
    # after delete polls, bounded by this
    exec: 10m
  execution:
    # Tests running at once. Each creates its own dependencies, so keep it
    # within the account's quotas, e.g. VPCs per region.
    parallel: 4
  # Manifests can use bindings, e.g. ($namespace)
  templating:
    enabled: true
  error:
    # Collected for any failing step, in addition to the step's own catch
    catch:
    - events: {}
    - podLogs:
        namespace: ack-system
        selector: app.kubernetes.io/name=ack-{{ .Service }}-controller
        tail: 200
