CREATE INDEX type ON blockchain.txns USING btree (type);
CREATE INDEX tx_detail_name ON blockchain.txns USING btree ((fields->'txDetail'->>'name'));