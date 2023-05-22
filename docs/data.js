db.user_tags.insertMany([
    {
        "time": ISODate("2023-01-21T12:30:00.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "china",
        "product_info": {
            "product_id": "adidas",
            "brand_id": "adidas",
            "category_id": "adidas",
            "price": 2
        }
    },
    {
        "time": ISODate("2023-01-21T12:30:01.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "china",
        "product_info": {
            "product_id": "adidas",
            "brand_id": "adidas",
            "category_id": "adidas",
            "price": 4
        }
    },
    {
        "time": ISODate("2023-01-21T12:30:01.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 5
        }
    },
    {
        "time": ISODate("2023-01-21T12:31:01.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 1
        }
    },
    {
        "time": ISODate("2023-01-21T12:31:02.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 1
        }
    },
    {
        "time": ISODate("2023-01-21T12:32:00.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 1
        }
    },
    {
        "time": ISODate("2023-01-21T12:32:01.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 1
        }
    },
    {
        "time": ISODate("2023-01-21T12:33:01.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 1
        }
    },
    {
        "time": ISODate("2023-01-21T12:33:04.000Z"),
        "cookie": "cookie2",
        "country": "test",
        "device": "PC",
        "action": "BUY",
        "origin": "test",
        "product_info": {
            "product_id": "nike",
            "brand_id": "nike",
            "category_id": "nike",
            "price": 1
        }
    }
])

// working query
db.user_tags.aggregate([{ $bucket: { groupBy: "$time", boundaries: ["2023-01-21T12:30:00.000Z", "2023-01-21T12:31:00.000Z", "2023-01-21T12:32:00.000Z"], default: "Other", output: { "count": { $sum: 1 }, "tags": { $push: { "time": "$time" } } } } }])

// groupy column wise but not time buckets
db.user_tags.aggregate([{ $group: { "_id": null, "time": { $push: "$time" }, "action": { $push: "$action" } } }])

