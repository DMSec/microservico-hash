# Microserviços em Golang, Python, MariaDB e Docker

![alt text](https://cdn-images-1.medium.com/max/800/1*I5kMbgX4qZkycpOFlcEbdw.png)


Descrição com detalhes do desafio: https://medium.com/dmsec/microservi%C3%A7os-python-golang-mariadb-e-dockers-e119a7285ed4

Nesse exemplo temos 2 microserviços e uma base de dados compartilhada.


1) Teste poderá ser feito pelo Postman;
2) Teste poderá ser pelo CURL;
3) Serviço de listagem - Responsável por expor e retornar uma API Rest de produtos cadastrados na base de dados;
4) Base de dados MariaDB com as tabelas Clientes, Produtos e Campanhas;
5) Serviço de desconto - Responsável por oferecer informações para o serviço de listagem de produtos, mas aplicando os descontos das campanhas. Nesse momento temos a campanha de blackfriday e a de aniversário;


## Clone este repositório

```
git clone https://github.com/DMSec/microservico-hash.git
cd microservico-hash
```

### Localhost

Se optar por rodar em localhost. Requisitos:
* Alterar o valor do host no código para a conexão com o banco de dados;
* Copiar as chaves para rodar localmente.

Para copiar:
```
cd keys
cp localhost\cert.pem cert.pem
cp localhost\private.key private.key
```

### Docker

Se optar por rodar em Docker. Requisitos:
* Copiar as chaves ou gerar uma nova chave para o serviço de desconto.

Para copiar:
```
cd keys
cp desconto\cert.pem cert.pem
cp desconto\private.key private.key
```
 
Para gerar as chaves:
```
openssl req -x509 -newkey rsa:4096 -keyout private.key -out cert.pem -days 365 -nodes -subj '/CN=desconto'
```
## Execução do docker-compose build - Construção dos nossos containers
```
docker-compose build
```
## Execução do docker compose up - Execução dos nossos containers
```
docker-compose up -d
```

Portas utilizadas pelos serviços:

* 3306  - Banco de dados;
* 11443 - Serviço de desconto;
* 11080 - Serviço de listagem;



## Listagem de produtos (sem usuário)
```
curl http://localhost:11080/products
```

## Listagem de produtos (com usuário aniversariante)

Neste caso, precisamos utilizar um ID no qual o dia da execução, seja o aniversário do cliente.

Você pode incluir no script de create_tables.sql, caso não exista.

Criei alguns registros na tabela até o dia 08/11/2019. Será aplicado 5% de desconto nos produtos.

```
curl -H 'X-USER-ID: 7' http://localhost:11080/products
```

### Blackfriday
Apesar da blackfriday ocorrer em uma sexta-feira, optei por criar uma rota de ativação e desativação de blackfriday, sendo assim para os testes de blackfriday
devemos ativar.


## Ativação blackfriday com alteração de % da campanha
```
curl -H 'blackfriday: 1' -H 'pct: 10' http://localhost:11080/blackfriday
```

## Teste durante a blackfriday

Neste caso, mesmo sendo aniversário do nosso cliente, valerá a regra da blackfriday que oferece os 10% limites no preço dos produtos.

```
curl -H 'X-USER-ID: 7' http://localhost:11080/products
```

## Desativação blackfriday com alteração de % da campanha
```
curl -H 'blackfriday: 0' -H 'pct: 10' http://localhost:11080/blackfriday
```

## Teste com usuário cadastrado e aniversário do usuário

Neste caso é necessário utilizar um ID de usuário existente na base e que seja aniversário dele. Você pode incluir no script de create_tables.sql, caso não exista.

```
curl -H 'X-USER-ID: 7' http://localhost:11080/products
```

