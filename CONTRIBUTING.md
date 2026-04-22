# Contributing to Ward

Thanks for your interest in contributing to Ward.

Ward is community-driven and open to contributions from anyone, in any programming language, for any part of the project (core, docs, tooling, tests, UI, automation, etc.).

## Ground rules

1. Be respectful and collaborative.
2. Keep changes focused and well explained.
3. Do not introduce code or content intended to harm people, systems, or communities.
4. If you use AI tools, you are responsible for the final output quality and security.

## Contribution flow

1. Fork the repository.
2. Create a branch: `feat/short-description`, `fix/short-description`, etc.
3. Make your changes following project patterns.
4. Open a Pull Request with a clear description of what changed and why.

## Commit standard

Use Conventional Commits:

- `feat: add request timeout for proxy service`
- `fix: validate refresh token type in middleware`
- `docs: update auth flow documentation`
- `refactor: simplify application service creation path`
- `test: add login invalid payload cases`
- `chore: update CI dependencies`

Allowed types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `build`, `ci`, `revert`.

## AI-assisted contributions

AI-generated contributions are welcome, but must follow the same engineering bar:

1. Review all generated code manually.
2. Verify it does not add known vulnerabilities (injection, auth bypass, insecure token/cookie handling, sensitive data leakage, unsafe defaults, etc.).
3. Ensure the change follows existing architecture and coding patterns.
4. Keep accountability: PR author remains responsible for the final code.

## Credits and attribution

This project is 100% open source. You can use, modify, and adapt it freely, including for startups and enterprise environments.

Please keep original copyright and license notices in redistributed code.

## Portuguese version

- PT-BR: [`./Docs/CONTRIBUTING-PT-BR.md`](./Docs/CONTRIBUTING-PT-BR.md)
