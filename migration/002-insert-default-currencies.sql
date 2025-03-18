INSERT INTO currencies (symbol, rate)
SELECT 'USD', 0
WHERE NOT EXISTS (SELECT 1 FROM currencies WHERE symbol = 'USD');

INSERT INTO currencies (symbol, rate)
SELECT 'BRL', 0
WHERE NOT EXISTS (SELECT 1 FROM currencies WHERE symbol = 'BRL');

INSERT INTO currencies (symbol, rate)
SELECT 'EUR', 0
WHERE NOT EXISTS (SELECT 1 FROM currencies WHERE symbol = 'EUR');

INSERT INTO currencies (symbol, rate)
SELECT 'BTC', 0
WHERE NOT EXISTS (SELECT 1 FROM currencies WHERE symbol = 'BTC');

INSERT INTO currencies (symbol, rate)
SELECT 'ETH', 0
WHERE NOT EXISTS (SELECT 1 FROM currencies WHERE symbol = 'ETH');
