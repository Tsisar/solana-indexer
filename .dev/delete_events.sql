UPDATE core.transactions
SET parsed = false;

DELETE FROM core.events WHERE true;