#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://127.0.0.1:8081/api}"
USERNAME="${USERNAME:-admin}"
PASSWORD="${PASSWORD:-admin123}"
TARGET_TENANT_ID="${TARGET_TENANT_ID:-1}"

echo "[1/6] Login"
LOGIN_RESP="$(curl -sS "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"${USERNAME}\",\"password\":\"${PASSWORD}\"}")"
TOKEN="$(echo "${LOGIN_RESP}" | sed -n 's/.*"accessToken":"\([^"]*\)".*/\1/p')"
if [[ -z "${TOKEN}" ]]; then
  echo "login failed: ${LOGIN_RESP}"
  exit 1
fi

AUTH_HEADER="Authorization: Bearer ${TOKEN}"

echo "[2/6] Access codes"
curl -sS "${BASE_URL}/auth/codes" -H "${AUTH_HEADER}" >/dev/null

echo "[3/6] Tenant-scoped user list"
curl -sS "${BASE_URL}/user/list" \
  -H "${AUTH_HEADER}" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: ${TARGET_TENANT_ID}" \
  -d '{"page":1,"pageSize":10}' >/dev/null

echo "[4/6] Tenant-scoped role list"
curl -sS "${BASE_URL}/role/list" \
  -H "${AUTH_HEADER}" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: ${TARGET_TENANT_ID}" \
  -d '{"page":1,"pageSize":10}' >/dev/null

echo "[5/6] Tenant-scoped menu list"
curl -sS "${BASE_URL}/menu/list" \
  -H "${AUTH_HEADER}" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: ${TARGET_TENANT_ID}" \
  -d '{"page":1,"pageSize":10}' >/dev/null

echo "[6/6] Forbidden check (role permission assignment without code should fail for limited users)"
FORBIDDEN_RESP="$(curl -sS "${BASE_URL}/role/permission" \
  -H "${AUTH_HEADER}" \
  -H "Content-Type: application/json" \
  -d '{"roleId":"1","menuIds":["1"]}')"
echo "${FORBIDDEN_RESP}" | sed -n '1p'

echo "multitenant regression finished"
