# Suite configuration for the {{ .Service }}-controller Chainsaw tests. test-infra passes
# this file to `chainsaw test`; JUnit reporting is set on the command line.
apiVersion: chainsaw.kyverno.io/v1alpha2
kind: Configuration
metadata:
  name: {{ .Service }}-controller
spec:
  timeouts:
    apply: 30s
    # Upper bounds, to catch a reconcile that never converges: a resource that
    # converges finishes long before them, and a CR that goes ACK.Terminal
    # fails its step at once. assert bounds the waits on a CR's conditions.
    # The slowest seen to converge took about 3 minutes (a transit gateway VPC
    # attachment).
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
