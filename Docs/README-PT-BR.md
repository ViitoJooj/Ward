# Ward

Ward e um API Gateway modular escrito em Go, com foco em seguranca, confiabilidade e visibilidade operacional.
Ele fica entre cliente e servicos backend para centralizar autenticacao, controle de requisicoes e registro de eventos.

## O que o Ward resolve

Em arquiteturas com varios servicos, e comum a autenticacao ficar espalhada, as regras de borda divergem e o monitoramento perde consistencia.
O Ward atua como ponto unico de entrada para padronizar esse comportamento.

Com isso, os servicos internos ficam livres para focar na regra de negocio.

## Como funciona por baixo dos panos

O fluxo de requisicao segue uma esteira previsivel:

1. Middleware de CORS valida origem e preflight.
2. Rotas protegidas passam pelo middleware de autenticacao via `access_token` em cookie.
3. Handlers cuidam apenas da camada HTTP (entrada/saida).
4. Services executam as regras de negocio.
5. Repositories fazem somente persistencia e consulta.
6. O log da requisicao e gravado de forma assincrona para nao bloquear a resposta.

Essa separacao reduz acoplamento e facilita evolucao segura.

## Arquitetura

- **Handlers**: interface HTTP.
- **Services**: orquestracao e regra de negocio.
- **Repositories**: acesso a dados.
- **DTOs**: contratos de entrada e saida da API.

## Seguranca

O projeto foi desenhado para inspirar confianca em ambiente real:

- Senhas sao armazenadas com hash bcrypt.
- Access token e refresh token sao separados, assinados e validados por tipo.
- Tokens de acesso possuem janela curta de validade.
- Endpoints protegidos validam identidade antes da execucao da logica de negocio.
- Metadados de requisicoes sao registrados para auditoria e resposta a incidentes.

## Confiabilidade e desempenho

- Baseado em Go com `fasthttp` para baixo overhead.
- Logging assincrono para preservar latencia.
- Limites de responsabilidade claros entre camadas.
- Controles de seguranca centralizados no gateway.

## Navegacao de documentos

- README em ingles: [`../README.md`](../README.md)
- OpenAPI/Swagger: [`./swagger/core-api.yml`](./swagger/core-api.yml)
- Guia de contribuicao: [`./CONTRIBUTING-PT-BR.md`](./CONTRIBUTING-PT-BR.md)
- Politica de seguranca: [`./SECURITY-PT-BR.md`](./SECURITY-PT-BR.md)
- Licenca (resumo PT-BR): [`./LICENSE-PT-BR.md`](./LICENSE-PT-BR.md)

## License

MIT
