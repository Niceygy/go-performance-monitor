DROP TABLE IF EXISTS stats;
CREATE TABLE stats (
    MID BIGINT,
    MNAME TINYTEXT,
    CPU BIGINT,
    RAM_TOTAL BIGINT,
    RAM_USED BIGINT,
    DISK BIGINT,
    SEC BIGINT
);