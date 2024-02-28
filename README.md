# rinha-backend-q1-2024

Desafio da rinha de backend 2024/Q1

## Repositório original da rinha
[Link para repositório da rinha](https://github.com/zanfranceschi/rinha-de-backend-2024-q1)

## Objetivo principal
Controle de concorrência com o tema de crédito e débitos.

## Endpoints
### Transações
**Requisição**
`POST /clientes/[id]/transacoes`
```json
{
    "valor": 1000,
    "tipo" : "c",
    "descricao" : "descricao"
}
```
Onde
- `[id]` (na URL) deve ser um número inteiro representando a identificação do cliente.
- `valor` deve ser um número inteiro positivo que representa centavos (não vamos trabalhar com frações de centavos). Por exemplo, R$ 10 são 1000 centavos.
- `tipo` deve ser apenas `c` para crédito ou `d` para débito.
- `descricao` deve ser uma string de 1 a 10 caracteres.

Todos os campos são obrigatórios.

**Resposta**
`HTTP 200 OK`
```json
{
    "limite" : 100000,
    "saldo" : -9098
}
```
Onde
- `limite` deve ser o limite cadastrado do cliente.
- `saldo` deve ser o novo saldo após a conclusão da transação.

### Extrato
**Requisição**
`GET /clientes/[id]/extrato`

Onde
- `[id]` (na URL) deve ser um número inteiro representando a identificação do cliente.

**Resposta**
`HTTP 200 OK`
```json
{
  "saldo": {
    "total": -9098,
    "data_extrato": "2024-01-17T02:34:41.217753Z",
    "limite": 100000
  },
  "ultimas_transacoes": [
    {
      "valor": 10,
      "tipo": "c",
      "descricao": "descricao",
      "realizada_em": "2024-01-17T02:34:38.543030Z"
    },
    {
      "valor": 90000,
      "tipo": "d",
      "descricao": "descricao",
      "realizada_em": "2024-01-17T02:34:38.543030Z"
    }
  ]
}
```
Onde
- `saldo`
    - `total` deve ser o saldo total atual do cliente (não apenas das últimas transações seguintes exibidas).
    - `data_extrato` deve ser a data/hora da consulta do extrato.
    - `limite` deve ser o limite cadastrado do cliente.
- `ultimas_transacoes` é uma lista ordenada por data/hora das transações de forma decrescente contendo até as 10 últimas transações com o seguinte:
    - `valor` deve ser o valor da transação.
    - `tipo` deve ser `c` para crédito e `d` para débito.
    - `descricao` deve ser a descrição informada durante a transação.
    - `realizada_em` deve ser a data/hora da realização da transação.

## Tecnologias Utilizadas:
- Linguagem: Go
- Framework Web: Fiber
- Banco de Dados: MongoDB
- Load balancer: NGINX

[![Stack](https://skillicons.dev/icons?i=go,mongodb,nginx)](https://skillicons.dev)

## Contato:
- [victormatheusmx132@gmail.com](mailto:victormatheusmx132@gmail.com)
- [Instagram](https://www.instagram.com/victorkf132/)

## Como Executar:
- Clone o repositório: git clone https://github.com/Victor-132/rinha-backend-2024-q1
- Navegue até o diretório: cd rinha-backend-2024-q1
- Execute o docker-compose: `docker-compose up`
- Aguarde os serviços subirem completamente.
- Acesse a API através de: http://localhost:9999
