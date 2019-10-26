# Microserviços em Golang, Python, MariaDB e Docker


<image>https://cdn-images-1.medium.com/max/800/1*I5kMbgX4qZkycpOFlcEbdw.png</image>


Nesse exemplo temos 2 microserviços e uma base de dados compartilhada.

1) Base de dados MariaDB com as tabelas Clientes, Produtos e Campanhas;
2) Serviço de listagem - Responsável por expor e retornar uma API Rest de produtos cadastrados na base de dados;
3) Serviço de desconto - Responsável por oferecer informações para o serviço de listagem de produtos, mas aplicando os descontos das campanhas. Nesse momento temos a campanha de blackfriday e a de aniversário;


Informações detalhadas podem ser encontradas no https://medium.com/dmsec

## Clone este repositório

```
git clone https://github.com/DMSec/microservico-hash.git
cd microservico-hash
```

## Execução do docker compose build - Construção dos nossos containers
```
docker-compose build
```
## Execução do docker compose up - Execução dos nossos containers
```
docker-compose up -d
```

Portas utilizadas pelos serviços:
11443 - Serviço de desconto;
11080 - Serviço de listagem;

## Ativação e desativação de blackfriday com alteração de % da campanha


## Teste com blackfriday ativa sem usuário no header do POST
```
curl http://localhost:11080/products
```

## Teste com usuário cadastrado e blackfriday desativada
```
curl -H 'X-USER-ID: 1' http://localhost:11080/products
```

## Teste com usuário cadastrado e aniversário do usuário

Neste caso é necessário utilizar um ID de usuário existente na base e que seja aniversário dele. Você pode incluir no script de create_tables.sql, caso não exista.

```
curl -H 'X-USER-ID: 7' http://localhost:11080/products
```

