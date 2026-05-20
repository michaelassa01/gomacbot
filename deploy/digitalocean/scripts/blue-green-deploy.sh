#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${K8S_NAMESPACE:-gomacbot}"
IMAGE="${IMAGE:?IMAGE is required}"
DEPLOYMENT_PREFIX="gomacbot-api"
SERVICE_ACTIVE="gomacbot-api"
HEALTH_PATH="/health"
HEALTH_RETRIES="${HEALTH_RETRIES:-30}"
HEALTH_INTERVAL="${HEALTH_INTERVAL:-5}"

get_active_slot() {
  kubectl get svc "$SERVICE_ACTIVE" -n "$NAMESPACE" \
    -o jsonpath='{.spec.selector.deploy-slot}' 2>/dev/null || echo ""
}

wait_for_rollout() {
  local deployment="$1"
  kubectl rollout status "deployment/${deployment}" -n "$NAMESPACE" --timeout=600s
}

health_check_slot() {
  local slot="$1"
  local svc="${DEPLOYMENT_PREFIX}-${slot}"
  local job="health-${slot}-$$"
  local attempt=1

  echo "Health check via service ${svc} (${HEALTH_PATH})..."
  while [ "$attempt" -le "$HEALTH_RETRIES" ]; do
    kubectl delete pod "$job" -n "$NAMESPACE" --ignore-not-found >/dev/null 2>&1 || true

    kubectl run "$job" \
      --restart=Never \
      -n "$NAMESPACE" \
      --image=curlimages/curl:8.5.0 \
      --command -- \
      curl -sf "http://${svc}.${NAMESPACE}.svc.cluster.local${HEALTH_PATH}"

    local waited=0
    while [ "$waited" -lt 60 ]; do
      phase="$(kubectl get pod "$job" -n "$NAMESPACE" -o jsonpath='{.status.phase}' 2>/dev/null || echo "")"
      case "$phase" in
        Succeeded)
          kubectl delete pod "$job" -n "$NAMESPACE" --ignore-not-found >/dev/null
          echo "Health check passed for slot: ${slot}"
          return 0
          ;;
        Failed|Unknown)
          kubectl logs "$job" -n "$NAMESPACE" 2>/dev/null || true
          kubectl delete pod "$job" -n "$NAMESPACE" --ignore-not-found >/dev/null
          break
          ;;
      esac
      sleep 2
      waited=$((waited + 2))
    done
    kubectl delete pod "$job" -n "$NAMESPACE" --ignore-not-found >/dev/null

    echo "Attempt ${attempt}/${HEALTH_RETRIES} failed; retrying in ${HEALTH_INTERVAL}s..."
    sleep "$HEALTH_INTERVAL"
    attempt=$((attempt + 1))
  done

  echo "Health check failed for slot: ${slot}"
  return 1
}

switch_traffic() {
  local new_active="$1"
  kubectl patch svc "$SERVICE_ACTIVE" -n "$NAMESPACE" --type merge -p \
    "{\"spec\":{\"selector\":{\"app\":\"gomacbot-api\",\"deploy-slot\":\"${new_active}\"}}}"
  echo "Traffic switched to slot: ${new_active}"
}

deploy_to_slot() {
  local slot="$1"
  local deployment="${DEPLOYMENT_PREFIX}-${slot}"

  kubectl set image "deployment/${deployment}" "api=${IMAGE}" -n "$NAMESPACE"
  wait_for_rollout "$deployment"
  health_check_slot "$slot"
  switch_traffic "$slot"
}

ACTIVE="$(get_active_slot)"
if [ -z "$ACTIVE" ]; then
  echo "No active slot found; bootstrapping blue as initial active slot."
  ACTIVE="blue"
fi

case "$ACTIVE" in
  blue) INACTIVE="green" ;;
  green) INACTIVE="blue" ;;
  *)
    echo "Unknown active slot '${ACTIVE}'; expected blue or green."
    exit 1
    ;;
esac

echo "Active slot: ${ACTIVE}"
echo "Deploying image to inactive slot: ${INACTIVE}"
echo "Image: ${IMAGE}"

deploy_to_slot "$INACTIVE"

echo "Blue-green deployment complete. Active slot is now: ${INACTIVE}"
