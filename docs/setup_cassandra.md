# Setup cassandra cluster

pip install ccm cqlsh

Using `ccm`
https://github.com/riptano/ccm

Create keyspace
```
$ echo "CREATE KEYSPACE mimuwapi WITH \
replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};" | cqlsh
```

Create type 
```
echo "
use mimuwapi;
CREATE TYPE mimuwapi.productInfo (
    product_id text
    brand_id text
    category_id text
    price int
);" | cqlsh
```

Create schemas
```
echo "
use mimuwapi;
create table userTagEvents (
id UUID,
time text,
cookie text,
country text,
device text,
action text,
origin text,
PRIMARY KEY(id)
);" | cqlsh
```
