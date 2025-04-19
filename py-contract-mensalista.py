import csv
import uuid
import boto3
import hashlib

# Configurações da tabela DynamoDB
TABLE_NAME = "tbl_tripcheck_contract"
COMPANY_ID = "67fcd7db-9bd7-4294-9049-ee0eae5669b6"
CONTRACT_TYPE = "MENSALISTA"

TABLE_NAME_USER = "tbl_tripcheck_user"

# Configuração do cliente DynamoDB
dynamodb = boto3.client("dynamodb", region_name="sa-east-1")

# Função para carregar e processar o arquivo CSV
def process_csv(file_path):
    try:
        with open(file_path, mode="r", encoding="utf-8") as file:
            reader = csv.DictReader(file)
            
            # Itera pelas linhas do CSV
            for row in reader:
                # Obtém os valores das colunas
                nome = row["Nome"]
                email = row["Email"]
                cpf = row["CPF"]
                contato = row["Contato"]
                local_embarque = row["Local"]
                dias_semana = row["DaysOfWeek"].split(",")  # Divide os dias em lista
                horario_retorno = row["HR"]

                # Gera o UUID para o contrato
                contract_id = str(uuid.uuid4())

                # Gera o SHA-256 para o usuario
                usuario_id = hashlib.sha256(cpf.encode()).hexdigest()

                # Gera o SHA-256 para o email
                usuario_sk = hashlib.sha256(email.encode()).hexdigest()

                # Prepara os dados para inserção no DynamoDB
                item = {
                    "pk": {"S": COMPANY_ID},  # E-mail como chave de partição
                    "sk": {"S": f"CONTRACT#{usuario_id}#{contract_id}"},  # SK com prefixo e UUID
                    "CompanyId": {"S": COMPANY_ID},  # Identificador fixo da empresa
                    "Tipo": {"S": CONTRACT_TYPE},  # Tipo do contrato
                    "Nome": {"S": nome},
                    "Contato": {"S": contato},
                    "Local": {"S": local_embarque},
                    "DaysOfWeek": {"SS": dias_semana},  # Conjunto de dias da semana
                    "HR": {"S": horario_retorno},
                    "Status": {"S": "ATIVO"}
                }

                # Insere o item no DynamoDB
                response = dynamodb.put_item(TableName=TABLE_NAME, Item=item)

                item_usuario = {
                    "pk": {"S": nome[0].upper()},  # E-mail como chave de partição
                    "sk": {"S": f"METADATA#{usuario_sk}"},  # SK com prefixo e UUID
                    "Id": {"S": usuario_id},
                    "Status": {"S": "ATIVO"}
                }

                response = dynamodb.put_item(TableName=TABLE_NAME_USER, Item=item_usuario)

                item_usuario = {
                    "pk": {"S": nome[0].upper()},  # E-mail como chave de partição
                    "sk": {"S": f"CONTRACT#{usuario_id}#{contract_id}"},  # SK com prefixo e UUID
                    "Perfil": {"SS": ["PASSAGEIRO"]}
                }

                response = dynamodb.put_item(TableName=TABLE_NAME_USER, Item=item_usuario)

                print(f"Contrato inserido: {response}")
    except Exception as e:
        print(f"Erro ao processar o arquivo CSV: {e}")

# Chamada da função para processar um arquivo
file_path = "passageiros.csv"  # Substitua pelo caminho do arquivo CSV
process_csv(file_path)
