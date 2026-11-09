#!/usr/bin/env bash
# SessionStart hook: inject the using-hackkit skill so every session starts with the
# team workflow and the case rules. Adapted from obra/superpowers hooks/session-start (MIT).
# Uses only bash builtins: at startup PATH may not be set up yet.
# Set HACKKIT_DISABLE=1 to skip injection.

[ "${HACKKIT_DISABLE:-}" = "1" ] && exit 0

case "$0" in
  */*) here="${0%/*}" ;;
  *) here="." ;;
esac
skill_md="${here}/../skills/using-hackkit/SKILL.md"

content="hackkit: using-hackkit skill not found"
if [ -r "$skill_md" ]; then
  IFS= read -r -d '' content < "$skill_md" || true
  while [[ $content == *$'\n' ]]; do content="${content%$'\n'}"; done
fi

# Drop the YAML front matter: the model needs the body.
if [[ $content == ---$'\n'* ]]; then
  rest="${content#---$'\n'}"
  content="${rest#*$'\n'---$'\n'}"
fi

esc="$content"
esc="${esc//\\/\\\\}"
esc="${esc//\"/\\\"}"
esc="${esc//$'\n'/\\n}"
esc="${esc//$'\r'/\\r}"
esc="${esc//$'\t'/\\t}"

ctx="<EXTREMELY_IMPORTANT>\nThe hackkit plugin is active. Below is your 'hackkit:using-hackkit' skill. Invoke other hackkit skills with the Skill tool.\n\n${esc}\n</EXTREMELY_IMPORTANT>"

trap '' PIPE
printf '{\n  "hookSpecificOutput": {\n    "hookEventName": "SessionStart",\n    "additionalContext": "%s"\n  }\n}\n' "$ctx" 2>/dev/null || true
exit 0
