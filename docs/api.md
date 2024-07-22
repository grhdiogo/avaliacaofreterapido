# Rest API

### [POST] /quote

Create a quotation

+ Header
  Content-Type: application/json

+ Body

| Attribute | Type | Required | Opt |
|-----------|------|----------|------|
| recipient  |    Recipent  |  True   |  -    |
| volumes  |    Array(Volume)  |  True   |  -    |


| Recipent | Type | Required | Opt |
|-----------|------|----------|------|
| address  |    Address  |  True   |  -    |


| Address | Type | Required | Opt |
|-----------|------|----------|------|
| zipcode  |    String  |  True   |  Only numbers    |


| Volume | Type | Required | Opt |
|-----------|------|----------|------|
| category  |    Number  |  True   |  Must be > 0    |
| amount  |    Number  |  True   |  Must be > 0    |
| unitary_weight  |    Number  |  True   |  Must be > 0    |
| price  |    Number  |  True   |  Must be > 0    |
| sku  |    String  |  False   | Max 255 caracters    |
| height  |    Number  |  True   |  Must be > 0    |
| width  |    Number  |  True   |  Must be > 0    |
| length  |    Number  |  True   |  Must be > 0    |



+ Response 200 (application/json)

| Attribute | Type  |
|-----------|------|
| carrier  |    Carrier  |  

| Carrier | Type  |
|-----------|------|
| name  |    String  |  
| service  |    String  |  
| deadline  |    String  |  
| price  |    Number  |  


### [GET] /metrics

+ Header

+ Query

| Attribute | Type | Required | Opt |
|-----------|------|----------|------|
| last_quotes  |    Number  |  False   |  -    |

+ Response 200 (application/json)

| Attribute | Type |
|-----------|------|
| carriersMetric  |    CarriersMetric[]  |

| CarriersMetric | Type |
|-----------|------|
| name  |    Number  |
| resultQuantity  |    Number  |
| totalValue  |    Number  |
| mostExpensiveFreight  |    Number  |
| cheaperFreight  |    Number  |
| averagePrice  |    Number  |