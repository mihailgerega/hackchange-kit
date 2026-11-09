# Third-party notices

hackkit vendors parts of the following projects. Modified files are marked in the table; each project's license text follows.

| Project | Commit/version | License | Files in this repo | Modified |
|---|---|---|---|---|
| [obra/superpowers](https://github.com/obra/superpowers) | v7.0.2 (22f0b46) | MIT | skills: test-driven-development, writing-plans, subagent-driven-development, executing-plans, verification-before-completion, requesting-code-review, systematic-debugging, dispatching-parallel-agents, using-git-worktrees, finishing-a-development-branch; hooks/session-start.sh (adapted) | yes: renamed references, OpenSpec integration, hackkit sections appended |
| [mattpocock/skills](https://github.com/mattpocock/skills) | 4ae2317 | MIT | skills: grilling, prototype | yes: description, hackathon section, issue → design.md |
| [samber/cc-skills-golang](https://github.com/samber/cc-skills-golang) | d68d80a | MIT | skills: golang-testing, golang-error-handling, golang-database, golang-security, golang-concurrency, golang-context, golang-design-patterns, golang-project-layout | yes: shortened descriptions, cross-references, evals removed |
| [affaan-m/ECC](https://github.com/affaan-m/ECC) | v2.2.3 (d2ef268) | MIT | skills: react-patterns, react-testing, api-design, postgres-patterns, redis-patterns; agents: go-reviewer, go-build-resolver, typescript-reviewer, react-reviewer, security-reviewer, database-reviewer | agents: yes (boilerplate removed, references, security addendum) |
| [anthropics/skills](https://github.com/anthropics/skills) | 2c5958e | Apache-2.0 | skills/frontend-design (with LICENSE.txt); skills/go-mcp-server/mcp_best_practices.md (with LICENSE-mcp_best_practices.txt) | no |

## obra/superpowers (MIT)

```
MIT License

Copyright (c) 2025 Jesse Vincent

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## mattpocock/skills (MIT)

```
MIT License

Copyright (c) 2026 Matt Pocock

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## samber/cc-skills-golang (MIT)

```
MIT License

Copyright (c) 2026 Samuel Berthe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## affaan-m/ECC (MIT)

```
MIT License

Copyright (c) 2026 Affaan Mustafa

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## anthropics/skills (Apache-2.0)

The full Apache License 2.0 text ships next to the vendored files: `plugins/hackkit/skills/frontend-design/LICENSE.txt` and `plugins/hackkit/skills/go-mcp-server/LICENSE-mcp_best_practices.txt`.

## Материалы Авито

`skills/avito-context/` и `skills/kickoff/templates/seed/onboarding-checklist.yaml` пересказывают открытые материалы, перечисленные в кейсе Hack&Change 2026 (чек-лист стажёра, AvitoTech Playbook, avito-tech/tech-internship, статьи AvitoTech на Хабре, avito.tech, start.avito.ru, доклад о дизайн-системе). Тексты не копируются, у каждого факта есть ссылка на источник (`skills/avito-context/references/sources.md`). Права на материалы принадлежат Авито.
