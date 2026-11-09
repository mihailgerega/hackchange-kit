#!/usr/bin/env bash
# Static submission checklist for the Hack&Change AI-buddy case.
# Usage: bash check.sh [repo-root]   (run from the project; prints ✅/❌/⚠️ lines, exit 1 on any ❌)
set -uo pipefail
ROOT=${1:-$(git rev-parse --show-toplevel 2>/dev/null || pwd)}
cd "$ROOT" || exit 2
fail=0
ok()   { printf '✅ %s\n' "$1"; }
bad()  { printf '❌ %s\n' "$1"; fail=1; }
warn() { printf '⚠️  %s\n' "$1"; }
has()  { [ -e "$1" ]; }
any()  { local p; for p in "$@"; do [ -e "$p" ] && return 0; done; return 1; }  # any existing path among (expanded) globs

echo "== Admission (any ❌ here means 0 points)"
has README.md && ok "README.md exists" || bad "README.md missing"
for s in "Запуск|Quick start|Как запустить" "Тест|Tests" "учётн|credentials|Тестовые пользователи" "Проверк|How to check|Сценари"; do
  grep -qiE "$s" README.md 2>/dev/null && ok "README mentions: $s" || bad "README lacks a section matching: $s"
done
grep -qiE 'docker compose|make up' README.md 2>/dev/null && ok "README has a run command (docker compose / make up)" || bad "README has no run command"
has docker-compose.yml || has compose.yaml && ok "compose file present" || bad "no docker-compose.yml"
if command -v docker >/dev/null; then docker compose config -q 2>/dev/null && ok "docker compose config is valid" || bad "docker compose config fails"; fi
compose=$(ls docker-compose.y*ml compose.y*ml 2>/dev/null | head -1)
if [ -n "$compose" ] && grep -qE '^\s*(env_file:\s*\.env\s*$|-\s*\.env\s*$)' "$compose" && ! grep -q 'required: false' "$compose"; then
  bad "$compose requires .env (env_file) — a clean clone has no .env; use \${VAR:-default} in environment"
else ok "compose starts without a .env file"; fi
ls *.pdf docs/*.pdf presentation/*.pdf >/dev/null 2>&1 && ok "presentation PDF found" || bad "no presentation PDF (root, docs/ or presentation/)"

echo "== Code, docs, versions"
grep -qiE 'go 1\.[0-9]+|node|версии|versions' README.md 2>/dev/null && ok "README lists dependency versions" || bad "README lacks dependency versions"
grep -qiE 'LLM_PROVIDER|ANTHROPIC_API_KEY|модел|model' README.md 2>/dev/null && ok "README documents the LLM service/model/env var" || bad "README must name LLM service, model, env var, key setup"
has .env.example && ok ".env.example present" || bad ".env.example missing"
if has .env.example && grep -qE '(sk-ant-[A-Za-z0-9_-]{10,}|sk-[A-Za-z0-9]{20,}|AKIA[0-9A-Z]{16})' .env.example; then bad ".env.example contains something that looks like a real key"; else ok ".env.example has no real-looking keys"; fi
git ls-files --error-unmatch .env >/dev/null 2>&1 && bad ".env is committed" || ok ".env not committed"
if git grep -nE 'sk-ant-[A-Za-z0-9_-]{20,}' -- . ':!*.md' >/dev/null 2>&1; then bad "possible API key committed (git grep sk-ant-)"; else ok "no Anthropic keys in tracked files"; fi

echo "== SDD specs"
n_specs=$(find openspec/specs -name spec.md 2>/dev/null | wc -l)
[ "$n_specs" -gt 0 ] && ok "openspec/specs: $n_specs capability spec(s)" || bad "no OpenSpec specs in openspec/specs (archive your changes)"
n_active=$(find openspec/changes -mindepth 1 -maxdepth 1 -type d ! -name archive 2>/dev/null | wc -l)
[ "$n_active" -eq 0 ] && ok "no unarchived changes" || warn "$n_active change(s) not archived — archive finished ones (/opsx:archive)"
if command -v openspec >/dev/null; then openspec validate --all >/dev/null 2>&1 && ok "openspec validate --all passes" || bad "openspec validate --all fails"; fi
grep -rqiE 'MUST|SHALL' openspec/specs 2>/dev/null && ok "specs use SHALL/MUST requirements" || warn "specs have no SHALL/MUST"
grep -rqiE 'прав|роль|role|permission|доступ' openspec/specs 2>/dev/null && ok "specs cover permissions" || bad "no permission scenarios in specs"

echo "== Tests and coverage"
if has coverage/REPORT.md; then
  ok "coverage/REPORT.md present"
  awk -F'|' '/%/ {gsub(/[^0-9.]/,"",$3); if ($3 != "" && $3+0 <= 80) exit 1}' coverage/REPORT.md && ok "all parts > 80% in coverage/REPORT.md" || bad "a part is ≤ 80% in coverage/REPORT.md"
else bad "coverage/REPORT.md missing — run make cover-check and commit the report"; fi
has scripts/coverage-check.sh && ok "coverage script present (reproducible)" || bad "scripts/coverage-check.sh missing"
ls frontend/e2e/*.spec.ts >/dev/null 2>&1 && ok "e2e tests present" || warn "no e2e tests in frontend/e2e"

echo "== What Avito reviewers look for (tech-internship reviews; warnings)"
any backend/.golangci.y*ml .golangci.y*ml && ok "golangci-lint config in repo" || warn "no .golangci.yml — Avito expects linter config in the repo"
any frontend/eslint.config.* frontend/.eslintrc* && ok "ESLint config in frontend" || warn "no ESLint config in frontend/"
n_mig=$(find . -path ./node_modules -prune -o -path '*/migrations/*.sql' -print 2>/dev/null | wc -l)
[ "$n_mig" -gt 1 ] && ok "migrations: $n_mig SQL files" || warn "migrations should be separate files (golang-migrate), not one init.sql"
for s in "Принятые решения|Decisions|Компромисс|Trade-?off" "mermaid" "journey|CJM|Путь стажёра" "erDiagram|Схема БД|индекс" "Пример .*запрос|request example|итогов(ый|ого) запрос" "гарантирует|guarantee"; do
  grep -qiE "$s" README.md docs/*.md 2>/dev/null && ok "README/docs cover: $s" || warn "README/docs lack: $s (decisions, C4, CJM, DB schema, LLM request example, guarantees)"
done
grep -rqiE 'pending' backend --include='*.go' 2>/dev/null && ok "runs/requests have a pending status" || warn "persist every LLM run as pending → success|failed"

echo "== Test data, video, access"
if [ -n "$(find . -path ./node_modules -prune -o \( -iname '*seed*' -o -name testdata \) -print 2>/dev/null | head -1)" ] || grep -rqiE 'seed|синтетич|synthetic' README.md docs 2>/dev/null; then ok "test data / seeds present or documented"; else warn "where is the synthetic test data? document it in README"; fi
grep -qiE '(youtu|rutube|disk\.yandex|drive\.google|\.mp4|\.webm)' README.md 2>/dev/null && ok "README links the demo video" || bad "README has no demo video link/file"
grep -rqiE 'тестов(ые|ая) инструкц|test data|синтетич' README.md docs 2>/dev/null && ok "test instructions labeled as test data" || warn "label knowledge-base documents as test data (not Avito policy)"

echo
echo "Manual (can't be checked here):"
echo "  - fresh clone (no .env) → docker compose up --build → main scenario passes with LLM_PROVIDER=fake"
echo "  - private repo: access granted to GitHub account hack-and-change-2026"
echo "  - video link on the first slide of the deck; deck exported to PDF"
echo "  - submission uploaded before the deadline (the version is frozen after it)"
exit $fail
