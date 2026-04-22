# Contribuindo com o Ward

Obrigado por contribuir com o Ward.

O Ward e um projeto open source, aberto para contribuicoes de qualquer pessoa, em qualquer linguagem de programacao, para qualquer parte do projeto (core, docs, tooling, testes, UI, automacao etc.).

## Regras basicas

1. Mantenha respeito e colaboracao.
2. Faça mudancas focadas e bem explicadas.
3. Nao introduza codigo ou conteudo com objetivo de causar dano a pessoas, sistemas ou comunidades.
4. Se usar IA, a responsabilidade final por qualidade e seguranca e do autor do PR.

## Fluxo de contribuicao

1. Fork este repositorio.
2. Crie uma branch: `feat/descricao-curta`, `fix/descricao-curta` etc.
3. Implemente as mudancas seguindo os padroes do projeto.
4. Abra um Pull Request com descricao clara do que mudou e por que.

## Padrao de commits

Use Conventional Commits:

- `feat: adiciona timeout de requisicao no proxy`
- `fix: valida tipo do refresh token no middleware`
- `docs: atualiza documentacao do fluxo de auth`
- `refactor: simplifica criacao no service de aplicacoes`
- `test: adiciona cenarios invalidos de login`
- `chore: atualiza dependencias de CI`

Tipos permitidos: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `build`, `ci`, `revert`.

## Contribuicoes com IA

Contribuicoes com IA sao bem-vindas, mas devem seguir o mesmo nivel tecnico:

1. Revise manualmente todo codigo gerado.
2. Garanta que nao introduz vulnerabilidades conhecidas (injecao, bypass de auth, token/cookie inseguro, vazamento de dados sensiveis, defaults inseguros etc.).
3. Mantenha consistencia com arquitetura e padroes existentes.
4. O autor do PR continua responsavel pelo resultado final.

## Creditos e atribuicao

Este projeto e 100% open source. Pode ser usado, modificado e customizado livremente, inclusive por startups e empresas.

Mantenha os avisos de copyright e licenca em redistribuicoes do codigo.

## Versao em ingles

- EN: [`../CONTRIBUTING.md`](../CONTRIBUTING.md)
