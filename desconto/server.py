import datetime
import logging
import string
import sys
import time
import os
import grpc
import decimal
import mysql.connector
import dmsec_pb2
import dmsec_pb2_grpc
from concurrent import futures
from mysql.connector import Error
from datetime import date

def getConnection():
    """GetConnection()
        Deve ser passado os parametros para conexão no Mariadb;
    """
    hostname = os.environ['MYSQL_HOST']
    username = os.environ['MYSQL_USER']
    password = os.environ['MYSQL_PASSWORD']
    database = os.environ['MYSQL_DBNAME']
    connection = mysql.connector.connect(host=hostname, user=username, passwd=password, db=database)
    return connection

def getBlackFriday():
    try:
        cnx = getConnection()
        cursor = cnx.cursor()
        print("Database version : ")
        query = "SELECT * FROM campanhas where status = 1 and campanha ='Blackfriday'"
        cursor.execute(query)
        records = cursor.fetchall()

        print("Total number of rows is: ", cursor.rowcount)
        print("Records: %s " % records)

        if cursor.rowcount > 0:
            print("Blackfriday true")
            return True
        else:
            print("Blackfriday false")
            return False

    except Error as e:
        print("Error reading data from MySQL table", e)
    finally:
        if (cnx.is_connected()):
            cnx.close()
            cursor.close()
            print("MySQL connection is closed")


def birthday(birthday):
    date_time_str = birthday
    date_time_obj = datetime.datetime.strptime(date_time_str, '%d/%m/%Y')
    print('Date:', date_time_obj.date())
    birthday = date_time_obj.date()
    today = datetime.date.today()
    print(today)
    days = today.day - birthday.day
    months = today.month - birthday.month
    retorno = False

    if (days == 0) & (months == 0):
        retorno = True
        return retorno
    else:
        return retorno

def getCampanhaPCT(campanha):
    try:
        cnx = getConnection()
        cursor = cnx.cursor()
        print("Database version : ")
        query = "SELECT pct FROM campanhas where status = 1 and campanha ='%s'"
        cursor.execute(query, (campanha,))
        records = cursor.fetchall()

        print("Total number of rows is: ", cursor.rowcount)
        print("Records: %s " % records)

        if cursor.rowcount > 0:
            print("Blackfriday true")
            return records
        else:
            print("Blackfriday false")
            return 0

    except Error as e:
        print("Error reading data from MySQL table", e)
    finally:
        if (cnx.is_connected()):
            cnx.close()
            cursor.close()
            print("MySQL connection is closed")



def clienteExistsAndBirthday(cliente):
    try:
        isBirthday = False
        cnx = getConnection()
        print("conection")
        cursor = cnx.cursor()
        print("conection")
        sql_select_query = "select * from clientes where id ='%s'"
        cursor.execute(sql_select_query, (cliente.id,))
        record = cursor.fetchall()

        for row in record:
            print("Id = ", row[0], )
            print("first_name = ", row[1])
            print("last_name = ", row[2])
            print("birthday  = ", row[3], "\n")
            print("Isbirthday?", birthday(row[3]))
            isBirthday = birthday(row[3])

        return isBirthday

    except Error as e:
        print("Error reading data from MySQL table", e)
    finally:
        if (cnx.is_connected()):
            cnx.close()
            cursor.close()
            print("MySQL connection is closed")


class Dmsec(dmsec_pb2_grpc.DescontoServicer):

    def AplicarDesconto(self, request, content):
        cliente = request.cliente
        produto = request.produto
        desconto = dmsec_pb2.DiscountValue()

        print("Entrei no Aplicardesconto")
        print("cliente:" + str(cliente.id))

        # Pode ser melhorado criando outro microservico para verificar os descontos para todos os clientes / exemplo em campanhas promocionais
        if (getBlackFriday()) and produto.price_in_cents > 0:
            pct = getCampanhaPCT("Blackfriday")
            percentual = decimal.Decimal(pct) / 100  # 10%
            price = decimal.Decimal(produto.price_in_cents) / 100
            novo_price = price - (price * percentual)
            value_in_cents = int(novo_price * 100)
            desconto = dmsec_pb2.DiscountValue(pct=percentual, value_in_cents=value_in_cents)
            produto_com_discount = dmsec_pb2.Produto(id=produto.id,
                                                     title=produto.title,
                                                     description=produto.description,
                                                     price_in_cents=produto.price_in_cents,
                                                     discount_value=desconto)
            return dmsec_pb2.DescontoResposta(produto=produto_com_discount)

        elif (clienteExistsAndBirthday(cliente)) and produto.price_in_cents > 0:
            # Para melhorar podemos parametrizar a porcentagem de desconto em outro microservico ou que busque do BD
            print('Entrei por aqui')
            pct =  getCampanhaPCT("Aniversario")
            percentual = decimal.Decimal(pct) / 100  # 05%
            price = decimal.Decimal(produto.price_in_cents) / 100
            novo_price = price - (price * percentual)
            value_in_cents = int(novo_price * 100)
            desconto = dmsec_pb2.DiscountValue(pct=percentual, value_in_cents=value_in_cents)
            produto_com_discount = dmsec_pb2.Produto(id=produto.id,
                                                     title=produto.title,
                                                     description=produto.description,
                                                     price_in_cents=produto.price_in_cents,
                                                     discount_value=desconto)
            return dmsec_pb2.DescontoResposta(produto=produto_com_discount)

        else:
            print("Sem desconto aplicado")
            percentual = 0  # Sem %
            value_in_cents = produto.price_in_cents
            desconto = dmsec_pb2.DiscountValue(pct=percentual, value_in_cents=value_in_cents)
            produto_com_discount = dmsec_pb2.Produto(id=produto.id,
                                                     title=produto.title,
                                                     description=produto.description,
                                                     price_in_cents=produto.price_in_cents,
                                                     discount_value=desconto)
            return dmsec_pb2.DescontoResposta(produto=produto_com_discount)


if __name__ == '__main__':
    port = sys.argv[1] if len(sys.argv) > 1 else 443
    host = '[::]:%s' % port
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=5))
    keys_dir = os.path.abspath(os.path.join('.', os.pardir, 'keys'))
    with open('%s/private.key' % keys_dir, 'rb') as f:
        private_key = f.read()
    with open('%s/cert.pem' % keys_dir, 'rb') as f:
        certificate_chain = f.read()
    server_credentials = grpc.ssl_server_credentials(((private_key, certificate_chain),))
    server.add_secure_port(host, server_credentials)
    dmsec_pb2_grpc.add_DescontoServicer_to_server(Dmsec(), server)

    logging.basicConfig(filename="server.log", level=logging.INFO)
    try:
        server.start()
        print('Serviço de Desconto na porta %s em execução' % host)
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        server.stop(0)
