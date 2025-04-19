# transportecheckin-services

aws dynamodb delete-table --table-name tbl_tripcheck_company

aws dynamodb create-table --table-name tbl_tripcheck_company --attribute-definitions AttributeName=pk,AttributeType=S AttributeName=sk,AttributeType=S --key-schema AttributeName=pk,KeyType=HASH AttributeName=sk,KeyType=RANGE 

aws dynamodb create-table \
    --table-name tbl_tripcheck_company \
    --attribute-definitions \
        AttributeName=pk,AttributeType=S \
        AttributeName=sk,AttributeType=S \
    --key-schema \
        AttributeName=pk,KeyType=HASH \
        AttributeName=sk,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST

aws dynamodb delete-item \
    --table-name tbl_tripcheck_company \
    --key '{"pk": {"S": "R"}, "sk": {"S": "COMPANY#67fcd7db-9bd7-4294-9049-ee0eae5669b6"}}'

aws dynamodb delete-item \
	--table-name tbl_tripcheck_contract \
	--key '{"pk": {"S": "67fcd7db-9bd7-4294-9049-ee0eae5669b6"}, "sk": {"S": "CONTRACT#d8828b3e-4140-429c-b0ec-1cb87d07ff5f#victorsjc@gmail.com"}}'

aws dynamodb put-item \
    --table-name tbl_tripcheck_company \
    --item '{
    "pk": {"S": "R"},
    "sk": {"S": "COMPANY#67fcd7db-9bd7-4294-9049-ee0eae5669b6"},
    "Id": {"S":"67fcd7db-9bd7-4294-9049-ee0eae5669b6"},
    "RazaoSocial": {"S":"RAPHA TUR TRANSPORTES E TURISMO LTDA"},
    "CNPJ":{"S":"44990330000147"},
    "Contato":{"S":"12-99154-7073"},
    "Email":{"S": ""}}'

aws dynamodb get-item \
    --table-name tbl_tripcheck_company \
    --key '{"pk": {"S": "R"}, "sk": {"S": "COMPANY#67fcd7db-9bd7-4294-9049-ee0eae5669b6"}}'


aws dynamodb create-table \
    --table-name tbl_tripcheck_contract \
    --attribute-definitions \
        AttributeName=pk,AttributeType=S \
        AttributeName=sk,AttributeType=S \
    --key-schema \
        AttributeName=pk,KeyType=HASH \
        AttributeName=sk,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST

aws dynamodb scan --table-name tbl_tripcheck_contract

aws dynamodb create-table \
    --table-name tbl_tripcheck_user \
    --attribute-definitions \
        AttributeName=pk,AttributeType=S \
        AttributeName=sk,AttributeType=S \
    --key-schema \
        AttributeName=pk,KeyType=HASH \
        AttributeName=sk,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST

aws dynamodb put-item \
    --table-name tbl_tripcheck_user \
    --item '{
    "pk": {"S": "V"},
    "sk": {"S": "#METADATA#sha-256(victorsjc@gmail.com)"},
    "Id": {"S":"{sha-256(CPF)}"},
    "Status": {"S":"ATIVO"},
    "DataCriacao":{"S": ""},
    "DataAtualizacao":{"S": ""}}'

aws dynamodb put-item \
    --table-name tbl_tripcheck_user \
    --item '{
    "pk": {"S": "V},
    "sk": {"S": "CONTRATO#{contrato-id}#{sha-256(CPF)}"},
    "Perfil": {"S":"PASSAGEIRO"},"Perfil":{"S": "PROPRIETARIO"}}'