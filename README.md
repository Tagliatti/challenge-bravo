# Challenge Bravo

## Descrição
Solução para o desafio [Challenge Bravo da Hurb](./CHALLENGE.md), desenvolvido em Go com SQLite. O projeto conta com testes de integração, migrações aplicadas ao iniciar a aplicação e uso de Docker.

O projeto permite converter valores monetários de moedas reais e fictícias em tempo real. Por padrão, há suporte para USD, BRL, EUR, BTC e ETH, mas é possível cadastrar mais moedas. Para moedas fictícias, é necessário informar o valor da taxa (Rate) em relação ao USD, que é a moeda padrão do projeto.

## Requisitos
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Git](https://git-scm.com/)

## Setup
1. Clone o repositório
2. Entre na pasta do projeto
3. Execute o projeto
   ```bash
   docker-compose up -d
   ```
   > Se você estiver executando pela primeira vez o container pode demorar um pouco para subir. Execute `docker-compose logs -f` para acompanhar o processo.

## Testes
Para rodar os testes, execute o comando:
```bash
docker-compose run --rm api go test ./...
```

Para rodar os testes com cobertura de código, execute o comando:
```bash
docker-compose run --rm api go test -coverprofile=coverage.out ./...
```

## Endpoints
### `GET /health`
Endpoint de healthcheck
```bash
curl "http://localhost:8080/health
```

### `GET /from={symbol}&to={symbol}&amount={number}`
Converte o valor de uma moeda para outra

```bash
curl "http://localhost:8080?from=USD&to=EUR&amount=1.234"
```

### `POST /`
Cria uma nova moeda

```bash
curl -X POST -d '{"symbol":"GTA", "rate": 1}' "http://localhost:8080"
```

### `PUT /{symbol}`
Atualiza o rate de uma moeda

```bash
curl -X PUT -d '{"rate": 0}' "http://localhost:8080/GTA"
```

### `DELETE /{symbol}`
Excluir uma moeda

```bash
curl -X DELETE "localhost:8080/GTA"
```
