// Copyright 2015 The Xorm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dialects

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"xorm.io/xorm/core"
	"xorm.io/xorm/schemas"
)

var (
	oracleReservedWords = map[string]bool{
		"ACCESS":                    true,
		"ACCOUNT":                   true,
		"ACTIVATE":                  true,
		"ADD":                       true,
		"ADMIN":                     true,
		"ADVISE":                    true,
		"AFTER":                     true,
		"ALL":                       true,
		"ALL_ROWS":                  true,
		"ALLOCATE":                  true,
		"ALTER":                     true,
		"ANALYZE":                   true,
		"AND":                       true,
		"ANY":                       true,
		"ARCHIVE":                   true,
		"ARCHIVELOG":                true,
		"ARRAY":                     true,
		"AS":                        true,
		"ASC":                       true,
		"AT":                        true,
		"AUDIT":                     true,
		"AUTHENTICATED":             true,
		"AUTHORIZATION":             true,
		"AUTOEXTEND":                true,
		"AUTOMATIC":                 true,
		"BACKUP":                    true,
		"BECOME":                    true,
		"BEFORE":                    true,
		"BEGIN":                     true,
		"BETWEEN":                   true,
		"BFILE":                     true,
		"BITMAP":                    true,
		"BLOB":                      true,
		"BLOCK":                     true,
		"BODY":                      true,
		"BY":                        true,
		"CACHE":                     true,
		"CACHE_INSTANCES":           true,
		"CANCEL":                    true,
		"CASCADE":                   true,
		"CAST":                      true,
		"CFILE":                     true,
		"CHAINED":                   true,
		"CHANGE":                    true,
		"CHAR":                      true,
		"CHAR_CS":                   true,
		"CHARACTER":                 true,
		"CHECK":                     true,
		"CHECKPOINT":                true,
		"CHOOSE":                    true,
		"CHUNK":                     true,
		"CLEAR":                     true,
		"CLOB":                      true,
		"CLONE":                     true,
		"CLOSE":                     true,
		"CLOSE_CACHED_OPEN_CURSORS": true,
		"CLUSTER":                   true,
		"COALESCE":                  true,
		"COLUMN":                    true,
		"COLUMNS":                   true,
		"COMMENT":                   true,
		"COMMIT":                    true,
		"COMMITTED":                 true,
		"COMPATIBILITY":             true,
		"COMPILE":                   true,
		"COMPLETE":                  true,
		"COMPOSITE_LIMIT":           true,
		"COMPRESS":                  true,
		"COMPUTE":                   true,
		"CONNECT":                   true,
		"CONNECT_TIME":              true,
		"CONSTRAINT":                true,
		"CONSTRAINTS":               true,
		"CONTENTS":                  true,
		"CONTINUE":                  true,
		"CONTROLFILE":               true,
		"CONVERT":                   true,
		"COST":                      true,
		"CPU_PER_CALL":              true,
		"CPU_PER_SESSION":           true,
		"CREATE":                    true,
		"CURRENT":                   true,
		"CURRENT_SCHEMA":            true,
		"CURREN_USER":               true,
		"CURSOR":                    true,
		"CYCLE":                     true,
		"DANGLING":                  true,
		"DATABASE":                  true,
		"DATAFILE":                  true,
		"DATAFILES":                 true,
		"DATAOBJNO":                 true,
		"DATE":                      true,
		"DBA":                       true,
		"DBHIGH":                    true,
		"DBLOW":                     true,
		"DBMAC":                     true,
		"DEALLOCATE":                true,
		"DEBUG":                     true,
		"DEC":                       true,
		"DECIMAL":                   true,
		"DECLARE":                   true,
		"DEFAULT":                   true,
		"DEFERRABLE":                true,
		"DEFERRED":                  true,
		"DEGREE":                    true,
		"DELETE":                    true,
		"DEREF":                     true,
		"DESC":                      true,
		"DIRECTORY":                 true,
		"DISABLE":                   true,
		"DISCONNECT":                true,
		"DISMOUNT":                  true,
		"DISTINCT":                  true,
		"DISTRIBUTED":               true,
		"DML":                       true,
		"DOUBLE":                    true,
		"DROP":                      true,
		"DUMP":                      true,
		"EACH":                      true,
		"ELSE":                      true,
		"ENABLE":                    true,
		"END":                       true,
		"ENFORCE":                   true,
		"ENTRY":                     true,
		"ESCAPE":                    true,
		"EXCEPT":                    true,
		"EXCEPTIONS":                true,
		"EXCHANGE":                  true,
		"EXCLUDING":                 true,
		"EXCLUSIVE":                 true,
		"EXECUTE":                   true,
		"EXISTS":                    true,
		"EXPIRE":                    true,
		"EXPLAIN":                   true,
		"EXTENT":                    true,
		"EXTENTS":                   true,
		"EXTERNALLY":                true,
		"FAILED_LOGIN_ATTEMPTS":     true,
		"FALSE":                     true,
		"FAST":                      true,
		"FILE":                      true,
		"FIRST_ROWS":                true,
		"FLAGGER":                   true,
		"FLOAT":                     true,
		"FLOB":                      true,
		"FLUSH":                     true,
		"FOR":                       true,
		"FORCE":                     true,
		"FOREIGN":                   true,
		"FREELIST":                  true,
		"FREELISTS":                 true,
		"FROM":                      true,
		"FULL":                      true,
		"FUNCTION":                  true,
		"GLOBAL":                    true,
		"GLOBALLY":                  true,
		"GLOBAL_NAME":               true,
		"GRANT":                     true,
		"GROUP":                     true,
		"GROUPS":                    true,
		"HASH":                      true,
		"HASHKEYS":                  true,
		"HAVING":                    true,
		"HEADER":                    true,
		"HEAP":                      true,
		"IDENTIFIED":                true,
		"IDGENERATORS":              true,
		"IDLE_TIME":                 true,
		"IF":                        true,
		"IMMEDIATE":                 true,
		"IN":                        true,
		"INCLUDING":                 true,
		"INCREMENT":                 true,
		"INDEX":                     true,
		"INDEXED":                   true,
		"INDEXES":                   true,
		"INDICATOR":                 true,
		"IND_PARTITION":             true,
		"INITIAL":                   true,
		"INITIALLY":                 true,
		"INITRANS":                  true,
		"INSERT":                    true,
		"INSTANCE":                  true,
		"INSTANCES":                 true,
		"INSTEAD":                   true,
		"INT":                       true,
		"INTEGER":                   true,
		"INTERMEDIATE":              true,
		"INTERSECT":                 true,
		"INTO":                      true,
		"IS":                        true,
		"ISOLATION":                 true,
		"ISOLATION_LEVEL":           true,
		"KEEP":                      true,
		"KEY":                       true,
		"KILL":                      true,
		"LABEL":                     true,
		"LAYER":                     true,
		"LESS":                      true,
		"LEVEL":                     true,
		"LIBRARY":                   true,
		"LIKE":                      true,
		"LIMIT":                     true,
		"LINK":                      true,
		"LIST":                      true,
		"LOB":                       true,
		"LOCAL":                     true,
		"LOCK":                      true,
		"LOCKED":                    true,
		"LOG":                       true,
		"LOGFILE":                   true,
		"LOGGING":                   true,
		"LOGICAL_READS_PER_CALL":    true,
		"LOGICAL_READS_PER_SESSION": true,
		"LONG":                      true,
		"MANAGE":                    true,
		"MASTER":                    true,
		"MAX":                       true,
		"MAXARCHLOGS":               true,
		"MAXDATAFILES":              true,
		"MAXEXTENTS":                true,
		"MAXINSTANCES":              true,
		"MAXLOGFILES":               true,
		"MAXLOGHISTORY":             true,
		"MAXLOGMEMBERS":             true,
		"MAXSIZE":                   true,
		"MAXTRANS":                  true,
		"MAXVALUE":                  true,
		"MIN":                       true,
		"MEMBER":                    true,
		"MINIMUM":                   true,
		"MINEXTENTS":                true,
		"MINUS":                     true,
		"MINVALUE":                  true,
		"MLSLABEL":                  true,
		"MLS_LABEL_FORMAT":          true,
		"MODE":                      true,
		"MODIFY":                    true,
		"MOUNT":                     true,
		"MOVE":                      true,
		"MTS_DISPATCHERS":           true,
		"MULTISET":                  true,
		"NATIONAL":                  true,
		"NCHAR":                     true,
		"NCHAR_CS":                  true,
		"NCLOB":                     true,
		"NEEDED":                    true,
		"NESTED":                    true,
		"NETWORK":                   true,
		"NEW":                       true,
		"NEXT":                      true,
		"NOARCHIVELOG":              true,
		"NOAUDIT":                   true,
		"NOCACHE":                   true,
		"NOCOMPRESS":                true,
		"NOCYCLE":                   true,
		"NOFORCE":                   true,
		"NOLOGGING":                 true,
		"NOMAXVALUE":                true,
		"NOMINVALUE":                true,
		"NONE":                      true,
		"NOORDER":                   true,
		"NOOVERRIDE":                true,
		"NOPARALLEL":                true,
		"NOREVERSE":                 true,
		"NORMAL":                    true,
		"NOSORT":                    true,
		"NOT":                       true,
		"NOTHING":                   true,
		"NOWAIT":                    true,
		"NULL":                      true,
		"NUMBER":                    true,
		"NUMERIC":                   true,
		"NVARCHAR2":                 true,
		"OBJECT":                    true,
		"OBJNO":                     true,
		"OBJNO_REUSE":               true,
		"OF":                        true,
		"OFF":                       true,
		"OFFLINE":                   true,
		"OID":                       true,
		"OIDINDEX":                  true,
		"OLD":                       true,
		"ON":                        true,
		"ONLINE":                    true,
		"ONLY":                      true,
		"OPCODE":                    true,
		"OPEN":                      true,
		"OPTIMAL":                   true,
		"OPTIMIZER_GOAL":            true,
		"OPTION":                    true,
		"OR":                        true,
		"ORDER":                     true,
		"ORGANIZATION":              true,
		"OSLABEL":                   true,
		"OVERFLOW":                  true,
		"OWN":                       true,
		"PACKAGE":                   true,
		"PARALLEL":                  true,
		"PARTITION":                 true,
		"PASSWORD":                  true,
		"PASSWORD_GRACE_TIME":       true,
		"PASSWORD_LIFE_TIME":        true,
		"PASSWORD_LOCK_TIME":        true,
		"PASSWORD_REUSE_MAX":        true,
		"PASSWORD_REUSE_TIME":       true,
		"PASSWORD_VERIFY_FUNCTION":  true,
		"PCTFREE":                   true,
		"PCTINCREASE":               true,
		"PCTTHRESHOLD":              true,
		"PCTUSED":                   true,
		"PCTVERSION":                true,
		"PERCENT":                   true,
		"PERMANENT":                 true,
		"PLAN":                      true,
		"PLSQL_DEBUG":               true,
		"POST_TRANSACTION":          true,
		"PRECISION":                 true,
		"PRESERVE":                  true,
		"PRIMARY":                   true,
		"PRIOR":                     true,
		"PRIVATE":                   true,
		"PRIVATE_SGA":               true,
		"PRIVILEGE":                 true,
		"PRIVILEGES":                true,
		"PROCEDURE":                 true,
		"PROFILE":                   true,
		"PUBLIC":                    true,
		"PURGE":                     true,
		"QUEUE":                     true,
		"QUOTA":                     true,
		"RANGE":                     true,
		"RAW":                       true,
		"RBA":                       true,
		"READ":                      true,
		"READUP":                    true,
		"REAL":                      true,
		"REBUILD":                   true,
		"RECOVER":                   true,
		"RECOVERABLE":               true,
		"RECOVERY":                  true,
		"REF":                       true,
		"REFERENCES":                true,
		"REFERENCING":               true,
		"REFRESH":                   true,
		"RENAME":                    true,
		"REPLACE":                   true,
		"RESET":                     true,
		"RESETLOGS":                 true,
		"RESIZE":                    true,
		"RESOURCE":                  true,
		"RESTRICTED":                true,
		"RETURN":                    true,
		"RETURNING":                 true,
		"REUSE":                     true,
		"REVERSE":                   true,
		"REVOKE":                    true,
		"ROLE":                      true,
		"ROLES":                     true,
		"ROLLBACK":                  true,
		"ROW":                       true,
		"ROWID":                     true,
		"ROWNUM":                    true,
		"ROWS":                      true,
		"RULE":                      true,
		"SAMPLE":                    true,
		"SAVEPOINT":                 true,
		"SB4":                       true,
		"SCAN_INSTANCES":            true,
		"SCHEMA":                    true,
		"SCN":                       true,
		"SCOPE":                     true,
		"SD_ALL":                    true,
		"SD_INHIBIT":                true,
		"SD_SHOW":                   true,
		"SEGMENT":                   true,
		"SEG_BLOCK":                 true,
		"SEG_FILE":                  true,
		"SELECT":                    true,
		"SEQUENCE":                  true,
		"SERIALIZABLE":              true,
		"SESSION":                   true,
		"SESSION_CACHED_CURSORS":    true,
		"SESSIONS_PER_USER":         true,
		"SET":                       true,
		"SHARE":                     true,
		"SHARED":                    true,
		"SHARED_POOL":               true,
		"SHRINK":                    true,
		"SIZE":                      true,
		"SKIP":                      true,
		"SKIP_UNUSABLE_INDEXES":     true,
		"SMALLINT":                  true,
		"SNAPSHOT":                  true,
		"SOME":                      true,
		"SORT":                      true,
		"SPECIFICATION":             true,
		"SPLIT":                     true,
		"SQL_TRACE":                 true,
		"STANDBY":                   true,
		"START":                     true,
		"STATEMENT_ID":              true,
		"STATISTICS":                true,
		"STOP":                      true,
		"STORAGE":                   true,
		"STORE":                     true,
		"STRUCTURE":                 true,
		"SUCCESSFUL":                true,
		"SWITCH":                    true,
		"SYS_OP_ENFORCE_NOT_NULL$":  true,
		"SYS_OP_NTCIMG$":            true,
		"SYNONYM":                   true,
		"SYSDATE":                   true,
		"SYSDBA":                    true,
		"SYSOPER":                   true,
		"SYSTEM":                    true,
		"TABLE":                     true,
		"TABLES":                    true,
		"TABLESPACE":                true,
		"TABLESPACE_NO":             true,
		"TABNO":                     true,
		"TEMPORARY":                 true,
		"THAN":                      true,
		"THE":                       true,
		"THEN":                      true,
		"THREAD":                    true,
		"TIMESTAMP":                 true,
		"TIME":                      true,
		"TO":                        true,
		"TOPLEVEL":                  true,
		"TRACE":                     true,
		"TRACING":                   true,
		"TRANSACTION":               true,
		"TRANSITIONAL":              true,
		"TRIGGER":                   true,
		"TRIGGERS":                  true,
		"TRUE":                      true,
		"TRUNCATE":                  true,
		"TX":                        true,
		"TYPE":                      true,
		"UB2":                       true,
		"UBA":                       true,
		"UID":                       true,
		"UNARCHIVED":                true,
		"UNDO":                      true,
		"UNION":                     true,
		"UNIQUE":                    true,
		"UNLIMITED":                 true,
		"UNLOCK":                    true,
		"UNRECOVERABLE":             true,
		"UNTIL":                     true,
		"UNUSABLE":                  true,
		"UNUSED":                    true,
		"UPDATABLE":                 true,
		"UPDATE":                    true,
		"USAGE":                     true,
		"USE":                       true,
		"USER":                      true,
		"USING":                     true,
		"VALIDATE":                  true,
		"VALIDATION":                true,
		"VALUE":                     true,
		"VALUES":                    true,
		"VARCHAR":                   true,
		"VARCHAR2":                  true,
		"VARYING":                   true,
		"VIEW":                      true,
		"WHEN":                      true,
		"WHENEVER":                  true,
		"WHERE":                     true,
		"WITH":                      true,
		"WITHOUT":                   true,
		"WORK":                      true,
		"WRITE":                     true,
		"WRITEDOWN":                 true,
		"WRITEUP":                   true,
		"XID":                       true,
		"YEAR":                      true,
		"ZONE":                      true,
	}

	oracleQuoter = schemas.Quoter{'"', '"', schemas.AlwaysReserve}
)

type oracle struct {
	Base
}

func (db *oracle) Init(d *core.DB, uri *URI) error {
	db.quoter = oracleQuoter
	return db.Base.Init(d, db, uri)
}

func (db *oracle) SQLType(c *schemas.Column) string {
	var res string
	switch t := c.SQLType.Name; t {
	case schemas.Bit, schemas.TinyInt, schemas.SmallInt, schemas.MediumInt, schemas.Int, schemas.Integer, schemas.BigInt, schemas.Bool, schemas.Serial, schemas.BigSerial:
		res = "NUMBER"
	case schemas.Binary, schemas.VarBinary, schemas.Blob, schemas.TinyBlob, schemas.MediumBlob, schemas.LongBlob, schemas.Bytea:
		return schemas.Blob
	case schemas.Time, schemas.DateTime, schemas.TimeStamp:
		res = schemas.Date
	case schemas.TimeStampz:
		res = "TIMESTAMP"
	case schemas.Float, schemas.Double, schemas.Numeric, schemas.Decimal:
		res = "NUMBER"
	case schemas.Text, schemas.MediumText, schemas.LongText, schemas.Json:
		res = "CLOB"
	case schemas.Char, schemas.Varchar, schemas.TinyText:
		res = "VARCHAR2"
	default:
		res = t
	}

	hasLen1 := (c.Length > 0)
	hasLen2 := (c.Length2 > 0)

	if hasLen2 {
		res += "(" + strconv.Itoa(c.Length) + "," + strconv.Itoa(c.Length2) + ")"
	} else if hasLen1 {
		res += "(" + strconv.Itoa(c.Length) + ")"
	}
	return res
}

func (db *oracle) AutoIncrStr() string {
	return "AUTO_INCREMENT"
}

func (db *oracle) IsReserved(name string) bool {
	_, ok := oracleReservedWords[strings.ToUpper(name)]
	return ok
}

func (db *oracle) DropTableSQL(tableName, autoincrCol string) ([]string, bool) {
	var sqls = []string{
		fmt.Sprintf("DROP TABLE `%s`", tableName),
	}
	if autoincrCol != "" {
		sqls = append(sqls, fmt.Sprintf("DROP SEQUENCE %s", seqName(tableName)))
	}
	return sqls, false
}

func (db *oracle) CreateTableSQL(table *schemas.Table, tableName string) ([]string, bool) {
	var sql = "CREATE TABLE "
	if tableName == "" {
		tableName = table.Name
	}

	quoter := db.Quoter()
	sql += quoter.Quote(tableName) + " ("

	pkList := table.PrimaryKeys

	for _, colName := range table.ColumnsSeq() {
		col := table.GetColumn(colName)
		sql += db.StringNoPk(col)
		sql = strings.TrimSpace(sql)
		sql += ", "
	}

	if len(pkList) > 0 {
		sql += "PRIMARY KEY ( "
		sql += quoter.Join(pkList, ",")
		sql += " ), "
	}

	var sqls = []string{sql[:len(sql)-2] + ")"}

	if table.AutoIncrColumn() != nil {
		var sql2 = fmt.Sprintf(`CREATE sequence %s 
			minvalue 1
       		nomaxvalue
       		start with 1
       		increment by 1
       		nocycle
			   nocache`, seqName(tableName))
		sqls = append(sqls, sql2)
	}
	return sqls, false
}

func (db *oracle) SetQuotePolicy(quotePolicy QuotePolicy) {
	switch quotePolicy {
	case QuotePolicyNone:
		var q = oracleQuoter
		q.IsReserved = schemas.AlwaysNoReserve
		db.quoter = q
	case QuotePolicyReserved:
		var q = oracleQuoter
		q.IsReserved = db.IsReserved
		db.quoter = q
	case QuotePolicyAlways:
		fallthrough
	default:
		db.quoter = oracleQuoter
	}
}

func (db *oracle) IndexCheckSQL(tableName, idxName string) (string, []interface{}) {
	args := []interface{}{tableName, idxName}
	return `SELECT INDEX_NAME FROM USER_INDEXES ` +
		`WHERE TABLE_NAME = :1 AND INDEX_NAME = :2`, args
}

func (db *oracle) IsTableExist(ctx context.Context, tableName string) (bool, error) {
	return db.HasRecords(ctx, `SELECT table_name FROM user_tables WHERE table_name = :1`, tableName)
}

func (db *oracle) IsColumnExist(ctx context.Context, tableName, colName string) (bool, error) {
	args := []interface{}{tableName, colName}
	query := "SELECT column_name FROM USER_TAB_COLUMNS WHERE table_name = :1" +
		" AND column_name = :2"
	return db.HasRecords(ctx, query, args...)
}

func seqName(tableName string) string {
	return "SEQ_" + strings.ToUpper(tableName)
}

func (db *oracle) GetColumns(ctx context.Context, tableName string) ([]string, map[string]*schemas.Column, error) {
	//s := "SELECT column_name,data_default,data_type,data_length,data_precision,data_scale," +
	//	"nullable FROM USER_TAB_COLUMNS WHERE table_name = :1"

	s := `select   column_name   from   user_cons_columns   
  where   constraint_name   =   (select   constraint_name   from   user_constraints   
			  where   table_name   =   :1  and   constraint_type   ='P')`
	var pkName string
	err := db.DB().QueryRowContext(ctx, s, tableName).Scan(&pkName)
	if err != nil {
		return nil, nil, err
	}

	s = `SELECT USER_TAB_COLS.COLUMN_NAME, USER_TAB_COLS.DATA_DEFAULT, USER_TAB_COLS.DATA_TYPE, USER_TAB_COLS.DATA_LENGTH, 
		USER_TAB_COLS.data_precision, USER_TAB_COLS.data_scale, USER_TAB_COLS.NULLABLE,
		user_col_comments.comments
		FROM USER_TAB_COLS 
		LEFT JOIN user_col_comments on user_col_comments.TABLE_NAME=USER_TAB_COLS.TABLE_NAME 
		AND user_col_comments.COLUMN_NAME=USER_TAB_COLS.COLUMN_NAME
		WHERE USER_TAB_COLS.table_name = :1`

	rows, err := db.DB().QueryContext(ctx, s, tableName)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols := make(map[string]*schemas.Column)
	colSeq := make([]string, 0)
	for rows.Next() {
		col := new(schemas.Column)
		col.Indexes = make(map[string]int)

		var colName, colDefault, nullable, dataType, dataPrecision, dataScale, comment *string
		var dataLen int

		err = rows.Scan(&colName, &colDefault, &dataType, &dataLen, &dataPrecision,
			&dataScale, &nullable, &comment)
		if err != nil {
			return nil, nil, err
		}

		col.Name = strings.Trim(*colName, `" `)
		if colDefault != nil {
			col.Default = *colDefault
			col.DefaultIsEmpty = false
		}

		if *nullable == "Y" {
			col.Nullable = true
		} else {
			col.Nullable = false
		}

		if comment != nil {
			col.Comment = *comment
		}
		if pkName != "" && pkName == col.Name {
			col.IsPrimaryKey = true

			has, err := db.HasRecords(ctx, "SELECT * FROM USER_SEQUENCES WHERE SEQUENCE_NAME = :1", seqName(tableName))
			if err != nil {
				return nil, nil, err
			}
			if has {
				col.IsAutoIncrement = true
			}

			fmt.Println("-----", pkName, col.Name, col.IsPrimaryKey)
		}

		var (
			ignore     bool
			dt         string
			len1, len2 int
		)
		dts := strings.Split(*dataType, "(")
		dt = dts[0]
		if len(dts) > 1 {
			lens := strings.Split(dts[1][:len(dts[1])-1], ",")
			if len(lens) > 1 {
				len1, _ = strconv.Atoi(lens[0])
				len2, _ = strconv.Atoi(lens[1])
			} else {
				len1, _ = strconv.Atoi(lens[0])
			}
		}

		switch dt {
		case "VARCHAR2":
			col.SQLType = schemas.SQLType{Name: schemas.Varchar, DefaultLength: len1, DefaultLength2: len2}
		case "NVARCHAR2":
			col.SQLType = schemas.SQLType{Name: schemas.NVarchar, DefaultLength: len1, DefaultLength2: len2}
		case "TIMESTAMP WITH TIME ZONE":
			col.SQLType = schemas.SQLType{Name: schemas.TimeStampz, DefaultLength: 0, DefaultLength2: 0}
		case "NUMBER":
			col.SQLType = schemas.SQLType{Name: schemas.Double, DefaultLength: len1, DefaultLength2: len2}
		case "LONG", "LONG RAW", "NCLOB", "CLOB":
			col.SQLType = schemas.SQLType{Name: schemas.Text, DefaultLength: 0, DefaultLength2: 0}
		case "RAW":
			col.SQLType = schemas.SQLType{Name: schemas.Binary, DefaultLength: 0, DefaultLength2: 0}
		case "ROWID":
			col.SQLType = schemas.SQLType{Name: schemas.Varchar, DefaultLength: 18, DefaultLength2: 0}
		case "AQ$_SUBSCRIBERS":
			ignore = true
		default:
			col.SQLType = schemas.SQLType{Name: strings.ToUpper(dt), DefaultLength: len1, DefaultLength2: len2}
		}

		if ignore {
			continue
		}

		if _, ok := schemas.SqlTypes[col.SQLType.Name]; !ok {
			return nil, nil, fmt.Errorf("Unknown colType %v %v", *dataType, col.SQLType)
		}

		col.Length = dataLen

		if col.SQLType.IsText() || col.SQLType.IsTime() {
			if !col.DefaultIsEmpty {
				col.Default = "'" + col.Default + "'"
			}
		}
		cols[col.Name] = col
		colSeq = append(colSeq, col.Name)
	}

	/*select *
	from user_tab_comments
	where Table_Name='用户表' */

	return colSeq, cols, nil
}

func (db *oracle) GetTables(ctx context.Context) ([]*schemas.Table, error) {
	s := "SELECT table_name FROM user_tables WHERE TABLESPACE_NAME = :1 AND table_name NOT LIKE :2"
	args := []interface{}{strings.ToUpper(db.uri.User), "%$%"}

	rows, err := db.DB().QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := make([]*schemas.Table, 0)
	for rows.Next() {
		table := schemas.NewEmptyTable()
		err = rows.Scan(&table.Name)
		if err != nil {
			return nil, err
		}

		tables = append(tables, table)
	}
	return tables, nil
}

func (db *oracle) GetIndexes(ctx context.Context, tableName string) (map[string]*schemas.Index, error) {
	args := []interface{}{tableName}
	s := "SELECT t.column_name,i.uniqueness,i.index_name FROM user_ind_columns t,user_indexes i " +
		"WHERE t.index_name = i.index_name and t.table_name = i.table_name and t.table_name =:1"

	rows, err := db.DB().QueryContext(ctx, s, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indexes := make(map[string]*schemas.Index, 0)
	for rows.Next() {
		var indexType int
		var indexName, colName, uniqueness string

		err = rows.Scan(&colName, &uniqueness, &indexName)
		if err != nil {
			return nil, err
		}

		indexName = strings.Trim(indexName, `" `)

		var isRegular bool
		if strings.HasPrefix(indexName, "IDX_"+tableName) || strings.HasPrefix(indexName, "UQE_"+tableName) {
			indexName = indexName[5+len(tableName):]
			isRegular = true
		}

		if uniqueness == "UNIQUE" {
			indexType = schemas.UniqueType
		} else {
			indexType = schemas.IndexType
		}

		var index *schemas.Index
		var ok bool
		if index, ok = indexes[indexName]; !ok {
			index = new(schemas.Index)
			index.Type = indexType
			index.Name = indexName
			index.IsRegular = isRegular
			indexes[indexName] = index
		}
		index.AddColumn(colName)
	}
	return indexes, nil
}

func (db *oracle) Filters() []Filter {
	return []Filter{
		&SeqFilter{Prefix: ":", Start: 1},
	}
}

// https://github.com/godror/godror
type godrorDriver struct {
}

func parseOracle(driverName, dataSourceName string) (*URI, error) {
	var connStr = dataSourceName
	if !strings.HasPrefix(connStr, "oracle://") {
		connStr = fmt.Sprintf("oracle://%s", connStr)
	}

	u, err := url.Parse(connStr)
	if err != nil {
		return nil, err
	}

	db := &URI{
		DBType: schemas.ORACLE,
		Host:   u.Hostname(),
		Port:   u.Port(),
		DBName: strings.TrimLeft(u.RequestURI(), "/"),
	}

	if u.User != nil {
		db.User = u.User.Username()
		db.Passwd, _ = u.User.Password()
	}

	fmt.Printf("%#v\n", db)

	if db.DBName == "" {
		return nil, errors.New("dbname is empty")
	}
	return db, nil
}

func (cfg *godrorDriver) Parse(driverName, dataSourceName string) (*URI, error) {
	return parseOracle(driverName, dataSourceName)
}

type oci8Driver struct {
}

// dataSourceName=user/password@ipv4:port/dbname
// dataSourceName=user/password@[ipv6]:port/dbname
func (p *oci8Driver) Parse(driverName, dataSourceName string) (*URI, error) {
	return parseOracle(driverName, dataSourceName)
}
