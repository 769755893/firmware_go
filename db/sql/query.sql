-- name: QueryAllFirmware
SELECT * FROM $1 WHERE %s = "%s" AND %s = 1 ORDER BY %s DESC