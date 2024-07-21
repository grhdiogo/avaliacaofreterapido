# avaliacaofreterapido

# ● Observação: Para consumir a API da Frete Rápido, você vai precisar dos dados obrigatórios:
# CNPJ Remetente: 25.438.296/0001-58 (apenas números) 
# Usar o mesmo CNPJ para "shipper.registered_number" e "dispatchers.registered_number"
# Token de autenticação: 1d52a9b6b78cf07b08586152459a5c90
# Código Plataforma: 5AKVkHqCn
# Cep: 29161-376 (dispatchers[*].zipcode
# "unitary_price" deve ser informado. 
# variaveis
FRETERAPIDO_HOST = "https://sp.freterapido.com/api/v3"
FRETERAPIDO_TOKEN="1d52a9b6b78cf07b08586152459a5c90"
FRETERAPIDO_PLATFORM_CODE="5AKVkHqCn"
FRETERAPIDO_DISPATHER_CEP="29161376"
FRETERAPIDO_CNPJ="25438296000158"

# TODO:
carrier.Deadline é string ou int??
Verificar qual o tipo de simulação será enviado
É um dispather com vários volumes ou vários dispathers com 1 volume
UnitaryPrice vai ser o valor que é passado, ou o valor dividido pela quantidade
Tratar sempre como dispather[0] já que só mandarei 1?