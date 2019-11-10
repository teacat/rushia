package rushia

import (
	"fmt"
	"strings"
)

// CollationType 是資料表的字符校對格式。
type CollationType string

// CharsetType 是資料表的字符集格式。
type CharsetType string

// EngineType 是資料表的引擎格式。
type EngineType string

const (
	// EngineInnoDB 是用於資料庫表格的 `InnoDB` 引擎格式。
	EngineInnoDB EngineType = "innodb"
	// EngineMyISAM 是用於資料庫表格的 `MyISAM` 引擎格式。
	EngineMyISAM EngineType = "myisam"
)

const (
	// CollationBig5 是用於資料庫表格的 `big5_chinese_ci` 字符校對格式。
	CollationBig5 CollationType = "big5_chinese_ci"
	// CollationDEC8 是用於資料庫表格的 `dec8_swedish_ci` 字符校對格式。
	CollationDEC8 CollationType = "dec8_swedish_ci"
	// CollationCP850 是用於資料庫表格的 `cp850_general_ci` 字符校對格式。
	CollationCP850 CollationType = "cp850_general_ci"
	// CollationHP8 是用於資料庫表格的 `hp8_english_ci` 字符校對格式。
	CollationHP8 CollationType = "hp8_english_ci"
	// CollationKOI8R 是用於資料庫表格的 `koi8r_general_ci` 字符校對格式。
	CollationKOI8R CollationType = "koi8r_general_ci"
	// CollationLatin1 是用於資料庫表格的 `latin1_swedish_ci` 字符校對格式。
	CollationLatin1 CollationType = "latin1_swedish_ci"
	// CollationLatin2 是用於資料庫表格的 `latin2_general_ci` 字符校對格式。
	CollationLatin2 CollationType = "latin2_general_ci"
	// CollationSwe7 是用於資料庫表格的 `swe7_swedish_ci` 字符校對格式。
	CollationSwe7 CollationType = "swe7_swedish_ci"
	// CollationASCII 是用於資料庫表格的 `ascii_general_ci` 字符校對格式。
	CollationASCII CollationType = "ascii_general_ci"
	// CollationUJIS 是用於資料庫表格的 `ujis_japanese_ci` 字符校對格式。
	CollationUJIS CollationType = "ujis_japanese_ci"
	// CollationSJIS 是用於資料庫表格的 `sjis_japanese_ci` 字符校對格式。
	CollationSJIS CollationType = "sjis_japanese_ci"
	// CollationHebrew 是用於資料庫表格的 `hebrew_general_ci` 字符校對格式。
	CollationHebrew CollationType = "hebrew_general_ci"
	// CollationTIS620 是用於資料庫表格的 `tis620_thai_ci` 字符校對格式。
	CollationTIS620 CollationType = "tis620_thai_ci"
	// CollationEUCKR 是用於資料庫表格的 `euckr_korean_ci` 字符校對格式。
	CollationEUCKR CollationType = "euckr_korean_ci"
	// CollationKOI8U 是用於資料庫表格的 `koi8u_general_ci` 字符校對格式。
	CollationKOI8U CollationType = "koi8u_general_ci"
	// CollationGB2312 是用於資料庫表格的 `gb2312_chinese_ci` 字符校對格式。
	CollationGB2312 CollationType = "gb2312_chinese_ci"
	// CollationGreek 是用於資料庫表格的 `greek_general_ci` 字符校對格式。
	CollationGreek CollationType = "greek_general_ci"
	// CollationCP1250 是用於資料庫表格的 `cp1250_general_ci` 字符校對格式。
	CollationCP1250 CollationType = "cp1250_general_ci"
	// CollationGBK 是用於資料庫表格的 `gbk_chinese_ci` 字符校對格式。
	CollationGBK CollationType = "gbk_chinese_ci"
	// CollationLatin5 是用於資料庫表格的 `latin5_turkish_ci` 字符校對格式。
	CollationLatin5 CollationType = "latin5_turkish_ci"
	// CollationARMSCII8 是用於資料庫表格的 `armscii8_general_ci` 字符校對格式。
	CollationARMSCII8 CollationType = "armscii8_general_ci"
	// CollationUTF8 是用於資料庫表格的 `utf8_general_ci` 字符校對格式。
	CollationUTF8 CollationType = "utf8_general_ci"
	// CollationUCS2 是用於資料庫表格的 `ucs2_general_ci` 字符校對格式。
	CollationUCS2 CollationType = "ucs2_general_ci"
	// CollationCP866 是用於資料庫表格的 `cp866_general_ci` 字符校對格式。
	CollationCP866 CollationType = "cp866_general_ci"
	// CollationKeybcs2 是用於資料庫表格的 `keybcs2_general_ci` 字符校對格式。
	CollationKeybcs2 CollationType = "keybcs2_general_ci"
	// CollationMacCE 是用於資料庫表格的 `macce_general_ci` 字符校對格式。
	CollationMacCE CollationType = "macce_general_ci"
	// CollationMacRoman 是用於資料庫表格的 `macroman_general_ci` 字符校對格式。
	CollationMacRoman CollationType = "macroman_general_ci"
	// CollationCP852 是用於資料庫表格的 `cp852_general_ci` 字符校對格式。
	CollationCP852 CollationType = "cp852_general_ci"
	// CollationLatin7 是用於資料庫表格的 `latin7_general_ci` 字符校對格式。
	CollationLatin7 CollationType = "latin7_general_ci"
	// CollationUTF8MB4 是用於資料庫表格的 `utf8mb4_general_ci` 字符校對格式。
	CollationUTF8MB4 CollationType = "utf8mb4_general_ci"
	// CollationCP1251 是用於資料庫表格的 `cp1251_general_ci` 字符校對格式。
	CollationCP1251 CollationType = "cp1251_general_ci"
	// CollationUTF16 是用於資料庫表格的 `utf16_general_ci` 字符校對格式。
	CollationUTF16 CollationType = "utf16_general_ci"
	// CollationCP1256 是用於資料庫表格的 `cp1256_general_ci` 字符校對格式。
	CollationCP1256 CollationType = "cp1256_general_ci"
	// CollationCP1257 是用於資料庫表格的 `cp1257_general_ci` 字符校對格式。
	CollationCP1257 CollationType = "cp1257_general_ci"
	// CollationUTF32 是用於資料庫表格的 `utf32_general_ci` 字符校對格式。
	CollationUTF32 CollationType = "utf32_general_ci"
	// CollationBinary 是用於資料庫表格的 `binary` 字符校對格式。
	CollationBinary CollationType = "binary"
	// CollationGEOSTD8 是用於資料庫表格的 `geostd8_general_ci` 字符校對格式。
	CollationGEOSTD8 CollationType = "geostd8_general_ci"
	// CollationCP932 是用於資料庫表格的 `cp932_japanese_ci` 字符校對格式。
	CollationCP932 CollationType = "cp932_japanese_ci"
	// CollationEUCJPMS 是用於資料庫表格的 `eucjpms_japanese_ci` 字符校對格式。
	CollationEUCJPMS CollationType = "eucjpms_japanese_ci"
)

const (
	// CharsetBig5 是用於資料庫表格的 `big5` 字符集。
	CharsetBig5 CharsetType = "big5"
	// CharsetDEC8 是用於資料庫表格的 `dec8` 字符集。
	CharsetDEC8 CharsetType = "dec8"
	// CharsetCP850 是用於資料庫表格的 `cp850` 字符集。
	CharsetCP850 CharsetType = "cp850"
	// CharsetHP8 是用於資料庫表格的 `hp8` 字符集。
	CharsetHP8 CharsetType = "hp8"
	// CharsetKOI8R 是用於資料庫表格的 `koi8r` 字符集。
	CharsetKOI8R CharsetType = "koi8r"
	// CharsetLatin1 是用於資料庫表格的 `latin1` 字符集。
	CharsetLatin1 CharsetType = "latin1"
	// CharsetLatin2 是用於資料庫表格的 `latin2` 字符集。
	CharsetLatin2 CharsetType = "latin2"
	// CharsetSwe7 是用於資料庫表格的 `swe7` 字符集。
	CharsetSwe7 CharsetType = "swe7"
	// CharsetASCII 是用於資料庫表格的 `ascii` 字符集。
	CharsetASCII CharsetType = "ascii"
	// CharsetUJIS 是用於資料庫表格的 `ujis` 字符集。
	CharsetUJIS CharsetType = "ujis"
	// CharsetSJIS 是用於資料庫表格的 `sjis` 字符集。
	CharsetSJIS CharsetType = "sjis"
	// CharsetHebrew 是用於資料庫表格的 `hebrew` 字符集。
	CharsetHebrew CharsetType = "hebrew"
	// CharsetTIS620 是用於資料庫表格的 `tis620` 字符集。
	CharsetTIS620 CharsetType = "tis620"
	// CharsetEUCKR 是用於資料庫表格的 `euckr` 字符集。
	CharsetEUCKR CharsetType = "euckr"
	// CharsetKOI8U 是用於資料庫表格的 `koi8u` 字符集。
	CharsetKOI8U CharsetType = "koi8u"
	// CharsetGB2312 是用於資料庫表格的 `gb2312` 字符集。
	CharsetGB2312 CharsetType = "gb2312"
	// CharsetGreek 是用於資料庫表格的 `greek` 字符集。
	CharsetGreek CharsetType = "greek"
	// CharsetCP1250 是用於資料庫表格的 `cp1250` 字符集。
	CharsetCP1250 CharsetType = "cp1250"
	// CharsetGBK 是用於資料庫表格的 `gbk` 字符集。
	CharsetGBK CharsetType = "gbk"
	// CharsetLatin5 是用於資料庫表格的 `latin5` 字符集。
	CharsetLatin5 CharsetType = "latin5"
	// CharsetARMSCII8 是用於資料庫表格的 `armscii8` 字符集。
	CharsetARMSCII8 CharsetType = "armscii8"
	// CharsetUTF8 是用於資料庫表格的 `utf8` 字符集。
	CharsetUTF8 CharsetType = "utf8"
	// CharsetUCS2 是用於資料庫表格的 `ucs2` 字符集。
	CharsetUCS2 CharsetType = "ucs2"
	// CharsetCP866 是用於資料庫表格的 `cp866` 字符集。
	CharsetCP866 CharsetType = "cp866"
	// CharsetKeybcs2 是用於資料庫表格的 `keybcs2` 字符集。
	CharsetKeybcs2 CharsetType = "keybcs2"
	// CharsetMacCE 是用於資料庫表格的 `macce` 字符集。
	CharsetMacCE CharsetType = "macce"
	// CharsetMacRoman 是用於資料庫表格的 `macroman` 字符集。
	CharsetMacRoman CharsetType = "macroman"
	// CharsetCP852 是用於資料庫表格的 `cp852` 字符集。
	CharsetCP852 CharsetType = "cp852"
	// CharsetLatin7 是用於資料庫表格的 `latin7` 字符集。
	CharsetLatin7 CharsetType = "latin7"
	// CharsetUTF8MB4 是用於資料庫表格的 `utf8mb4` 字符集。
	CharsetUTF8MB4 CharsetType = "utf8mb4"
	// CharsetCP1251 是用於資料庫表格的 `cp1251` 字符集。
	CharsetCP1251 CharsetType = "cp1251"
	// CharsetUTF16 是用於資料庫表格的 `utf16` 字符集。
	CharsetUTF16 CharsetType = "utf16"
	// CharsetCP1256 是用於資料庫表格的 `cp1256` 字符集。
	CharsetCP1256 CharsetType = "cp1256"
	// CharsetCP1257 是用於資料庫表格的 `cp1257` 字符集。
	CharsetCP1257 CharsetType = "cp1257"
	// CharsetUTF32 是用於資料庫表格的 `utf32` 字符集。
	CharsetUTF32 CharsetType = "utf32"
	// CharsetBinary 是用於資料庫表格的 `binary` 字符集。
	CharsetBinary CharsetType = "binary"
	// CharsetGEOSTD8 是用於資料庫表格的 `geostd8` 字符集。
	CharsetGEOSTD8 CharsetType = "geostd8"
	// CharsetCP932 是用於資料庫表格的 `cp932` 字符集。
	CharsetCP932 CharsetType = "cp932"
	// CharsetEUCJPMS 是用於資料庫表格的 `eucjpms` 字符集。
	CharsetEUCJPMS CharsetType = "eucjpms"
)

// Migration 是一個資料庫表格的遷移系統。
type Migration struct {
	table   table
	columns []column

	// LasyQuery 是最後一次所執行的 SQL 指令。
	LastQuery string
}

// table 是一個資料表格與其資訊。
type table struct {
	name        string
	comment     string
	charset     CharsetType
	collation   CollationType
	primaryKeys []key
	indexKeys   []key
	uniqueKeys  []key
	foreignKeys []key
	engineType  EngineType
}

// column 是單個欄位與其資訊。
type column struct {
	name          string
	dataType      string
	length        interface{}
	comment       string
	unsigned      bool
	primary       bool
	unique        bool
	foreign       string
	autoIncrement bool
	defaultValue  interface{}
	nullable      bool
	extras        bool
}

// key 是一個索引的資料。
type key struct {
	name          string
	columns       []string
	targetColumns []string
	onUpdate      string
	onDelete      string
}

// TinyInt 會將最後一個欲建立的欄位資料型態設置為 `tinyint`。
func (m Migration) TinyInt(length int) Migration {
	return m.setColumnType("tinyint", length)
}

// SmallInt 會將最後一個欲建立的欄位資料型態設置為 `smallint`。
func (m Migration) SmallInt(length int) Migration {
	return m.setColumnType("smallint", length)
}

// MediumInt 會將最後一個欲建立的欄位資料型態設置為 `mediumint`。
func (m Migration) MediumInt(length int) Migration {
	return m.setColumnType("mediumint", length)
}

// Int 會將最後一個欲建立的欄位資料型態設置為 `int`。
func (m Migration) Int(length int) Migration {
	return m.setColumnType("int", length)
}

// BigInt 會將最後一個欲建立的欄位資料型態設置為 `bigint`。
func (m Migration) BigInt(length int) Migration {
	return m.setColumnType("bigint", length)
}

// Char 會將最後一個欲建立的欄位資料型態設置為 `char`。
func (m Migration) Char(length int) Migration {
	return m.setColumnType("char", length)
}

// Varchar 會將最後一個欲建立的欄位資料型態設置為 `varchar`。
func (m Migration) Varchar(length int) Migration {
	return m.setColumnType("varchar", length)
}

// TinyText 會將最後一個欲建立的欄位資料型態設置為 `tinytext`。
func (m Migration) TinyText() Migration {
	return m.setColumnType("tinytext")
}

// Text 會將最後一個欲建立的欄位資料型態設置為 `text`。
func (m Migration) Text() Migration {
	return m.setColumnType("text")
}

// MediumText 會將最後一個欲建立的欄位資料型態設置為 `mediumtext`。
func (m Migration) MediumText() Migration {
	return m.setColumnType("mediumtext")
}

// LongText 會將最後一個欲建立的欄位資料型態設置為 `longtext`。
func (m Migration) LongText() Migration {
	return m.setColumnType("longtext")
}

// Binary 會將最後一個欲建立的欄位資料型態設置為 `binary`。
func (m Migration) Binary(length int) Migration {
	return m.setColumnType("binary", length)
}

// Varbinary 會將最後一個欲建立的欄位資料型態設置為 `varbinary`。
func (m Migration) Varbinary(length int) Migration {
	return m.setColumnType("varbinary", length)
}

// Bit 會將最後一個欲建立的欄位資料型態設置為 `bit`。
func (m Migration) Bit(length int) Migration {
	return m.setColumnType("bit", length)
}

// TinyBlob 會將最後一個欲建立的欄位資料型態設置為 `tinyblob`。
func (m Migration) TinyBlob() Migration {
	return m.setColumnType("tinyblob")
}

// Blob 會將最後一個欲建立的欄位資料型態設置為 `blob`。
func (m Migration) Blob() Migration {
	return m.setColumnType("blob")
}

// MediumBlob 會將最後一個欲建立的欄位資料型態設置為 `mediumblob`。
func (m Migration) MediumBlob() Migration {
	return m.setColumnType("mediumblob")
}

// LongBlob 會將最後一個欲建立的欄位資料型態設置為 `longblob`。
func (m Migration) LongBlob() Migration {
	return m.setColumnType("longblob")
}

// Date 會將最後一個欲建立的欄位資料型態設置為 `date`。
func (m Migration) Date() Migration {
	return m.setColumnType("date")
}

// DateTime 會將最後一個欲建立的欄位資料型態設置為 `datetime`。
func (m Migration) DateTime() Migration {
	return m.setColumnType("dateTime")
}

// Time 會將最後一個欲建立的欄位資料型態設置為 `time`。
func (m Migration) Time() Migration {
	return m.setColumnType("time")
}

// Timestamp 會將最後一個欲建立的欄位資料型態設置為 `timestamp`。
func (m Migration) Timestamp() Migration {
	return m.setColumnType("timestamp")
}

// Year 會將最後一個欲建立的欄位資料型態設置為 `year`。
func (m Migration) Year() Migration {
	return m.setColumnType("year")
}

// Double 會將最後一個欲建立的欄位資料型態設置為 `double`。
func (m Migration) Double(length ...int) Migration {
	return m.setColumnType("double", length)
}

// Decimal 會將最後一個欲建立的欄位資料型態設置為 `decimal`。
//     .Decimal(2, 1)
func (m Migration) Decimal(length ...int) Migration {
	return m.setColumnType("decimal", length)
}

// Float 會將最後一個欲建立的欄位資料型態設置為 `float`。
//     .Float(2, 1)
//     .Float(1)
func (m Migration) Float(length ...int) Migration {
	return m.setColumnType("float", length)
}

// Enum 會將最後一個欲建立的欄位資料型態設置為 `enum`。
//     .Enum(1, 2, "A", "B")
func (m Migration) Enum(types ...interface{}) Migration {
	return m.setColumnType("enum", types)
}

// Set 會將最後一個欲建立的欄位資料型態設置為 `set`。
//     .Set(1, 2, "A", "B")
func (m Migration) Set(types ...interface{}) Migration {
	return m.setColumnType("set", types)
}

// Column 會建立一個新的欄位。
func (m Migration) Column(name string) Migration {
	m.columns = append(m.columns, column{name: name, defaultValue: false})
	return m
}

// Charset 會設置資料表格的字符集。
func (m Migration) Charset(charset CharsetType) Migration {
	m.table.charset = charset
	return m
}

// Collation 會設置資料表格的字符校對格式。
func (m Migration) Collation(collation CollationType) Migration {
	m.table.collation = collation
	return m
}

// Engine 能夠設置資料表的引擎種類。
func (m Migration) Engine(engine EngineType) Migration {
	m.table.engineType = engine
	return m
}

// Primary 會在沒有參數的情況下將某個欄位設定為主鍵。
//     .Column("id").Primary()
// 當傳入的參數是一個字串切片時，會將這些欄位名稱作為主鍵群組。
//     .Primary([]string{"id", "username"})
// 當第一個參數是字串，第二個則是字串切片時則會建立一個命名的主鍵群組。
//     .Primary("pk_group", []string{"id", "username"})
func (m Migration) Primary(args ...interface{}) Migration {
	switch len(args) {
	// Primary()
	case 0:
		m.columns[len(m.columns)-1].primary = true

	// Primary([]string{"column1", "column2"})
	case 1:
		m.table.primaryKeys = append(m.table.primaryKeys, key{
			columns: args[0].([]string),
		})

	// Primary("primary_keys", []string{"column1", "column2"})
	case 2:
		m.table.primaryKeys = append(m.table.primaryKeys, key{
			name:    args[0].(string),
			columns: args[1].([]string),
		})
	}
	return m
}

// Unique 會在沒有參數的情況下將某個欄位設定為不重覆鍵。
//     .Column("id").Unique()
// 當傳入的參數是一個字串切片時，會將這些欄位名稱作為不重覆鍵群組。
//     .Unique([]string{"id", "username"})
// 當第一個參數是字串，第二個則是字串切片時則會建立一個命名的不重覆鍵群組。
//     .Unique("uk_group", []string{"id", "username"})
func (m Migration) Unique(args ...interface{}) Migration {
	switch len(args) {
	// Unique()
	case 0:
		m.columns[len(m.columns)-1].unique = true

	// Unique([]string{"column1", "column2"})
	case 1:
		m.table.uniqueKeys = append(m.table.uniqueKeys, key{
			columns: args[0].([]string),
		})

	// Unique("unique_keys", []string{"column1", "column2"})
	case 2:
		m.table.uniqueKeys = append(m.table.uniqueKeys, key{
			name:    args[0].(string),
			columns: args[1].([]string),
		})
	}
	return m
}

// Index 會在沒有參數的情況下將某個欄位設定為索引。
//     .Column("id").Index()
// 當傳入的參數是一個字串切片時，會將這些欄位名稱作為索引群組。
//     .Index([]string{"id", "username"})
// 當第一個參數是字串，第二個則是字串切片時則會建立一個命名的索引群組。
//     .Index("ik_group", []string{"id", "username"})
func (m Migration) Index(args ...interface{}) Migration {
	switch len(args) {
	// Index()
	case 0:
		// 建立一個匿名的索引群組並將最後一個新增的欄位放入群組裡。
		m.table.indexKeys = append(m.table.indexKeys, key{
			columns: []string{m.columns[len(m.columns)-1].name},
		})

	// Index([]string{"column1", "column2"})
	case 1:
		// 替每個欄位各建立一個以該欄位為名的群組。
		for _, v := range args[0].([]string) {
			m.table.indexKeys = append(m.table.indexKeys, key{
				name:    v,
				columns: []string{v},
			})
		}

	// Index("index_keys", []string{"column1", "column2"})
	case 2:
		m.table.indexKeys = append(m.table.indexKeys, key{
			name:    args[0].(string),
			columns: args[1].([]string),
		})
	}
	return m
}

// OnUpdate 能夠決定外鍵資料變更時，相關欄位該做什麼處置（例如：`NO ACTION`、`SET NULL`、等）。
func (m Migration) OnUpdate(action string) Migration {
	m.table.foreignKeys[len(m.table.foreignKeys)-1].onUpdate = action
	return m
}

// OnDelete 能夠決定外鍵資料被刪除時，相關欄位該做什麼處置（例如：`NO ACTION`、`SET NULL`、等）。
func (m Migration) OnDelete(action string) Migration {
	m.table.foreignKeys[len(m.table.foreignKeys)-1].onDelete = action
	return m
}

// Foreign 會將最後一個欄位設置為外鍵，或者是建立一個匿名／命名群組。
// 欲將最後一個欄位設置為外鍵欄位可以這麼做：
//     Columns("id").Foreign("users.id")
// 欲建立一個匿名的外鍵群組則是像這樣：
//     .Foreign([]string{"id", "username"}, []string{"users.id", "users.username"})
// 而一個命名的外鍵群組，第一個參數則是群組的名稱，其他參數與匿名群組無異。
//     .Foreign("fk_group", []string{"id", "username"}, []string{"users.id", "users.username"})
func (m Migration) Foreign(args ...interface{}) Migration {
	switch len(args) {
	// Foreign("users.id")
	case 1:
		// 透過分割字串中的 `.` 點號來分析 `資料表.欄位` 格式並取得資料表與欄位名稱。
		targetTable := strings.Split(args[0].(string), ".")[0]
		// 掃描目前資料表格的所有外鍵，並查看是否有已經存在的目標資料表格。
		index := -1
		for k, v := range m.table.foreignKeys {
			// 分割 `資料表.欄位` 字串格式並取得外鍵的目標資料表與欄位名稱。
			tableName := strings.Split(v.targetColumns[0], ".")[0]
			if tableName == targetTable {
				index = k
				break
			}
		}
		// 如果已經有了相同的目標資料表，我們就將新的外鍵合併到該群組中。
		if index != -1 {
			m.table.foreignKeys[index].columns = append(m.table.foreignKeys[index].columns, m.columns[len(m.columns)-1].name)
			m.table.foreignKeys[index].targetColumns = append(m.table.foreignKeys[index].targetColumns, args[0].(string))
			// 不然就建立新的群組。
		} else {
			m.table.foreignKeys = append(m.table.foreignKeys, key{
				columns:       []string{m.columns[len(m.columns)-1].name},
				targetColumns: []string{args[0].(string)},
			})
		}

	// Foreign([]string{"id", "password"}, []string{"users.id", "users.password"})
	case 2:
		m.table.foreignKeys = append(m.table.foreignKeys, key{
			columns:       args[0].([]string),
			targetColumns: args[1].([]string),
		})

	// Foreign("foreign_keys", []string{"id", "password"}, []string{"users.id", "users.password"})
	case 3:
		m.table.foreignKeys = append(m.table.foreignKeys, key{
			name:          args[0].(string),
			columns:       args[1].([]string),
			targetColumns: args[2].([]string),
		})
	}
	return m
}

// Nullable 會將最後一個欄位設置為允許空值。
func (m Migration) Nullable() Migration {
	m.columns[len(m.columns)-1].defaultValue = nil
	m.columns[len(m.columns)-1].nullable = true
	return m
}

// Unsigned 會將最後一個欄位設置為非負數欄位。
func (m Migration) Unsigned() Migration {
	m.columns[len(m.columns)-1].unsigned = true
	return m
}

// Comment 能夠替最後一個欄位設置說明。
func (m Migration) Comment(text string) Migration {
	m.columns[len(m.columns)-1].comment = text
	return m
}

// Default 能夠替最後一個欄位設置預設值，可以是 nil、字串、正整數或者 `CURRENT_TIMESTAMP` 和 `NOW()`。
func (m Migration) Default(value interface{}) Migration {
	m.columns[len(m.columns)-1].defaultValue = value
	return m
}

// AutoIncrement 會將最後一個欄位設置為自動遞增，這僅能用於正整數欄位上。
func (m Migration) AutoIncrement() Migration {
	m.columns[len(m.columns)-1].autoIncrement = true
	return m
}

// Table 會準備一個資料表格供後續建立。
func (m Migration) Table(tableName string, comment ...string) Migration {
	// 設置表格名稱。
	m.table.name = tableName
	// 如果有指定表格備註的話就將其保存。
	if len(comment) != 0 {
		m.table.comment = comment[0]
	}
	return m
}

// Create 會執行先前的所有設置並且建立出一個相對應資料表格與其中的所有欄位。
func (m Migration) Create() (query string) {
	// 建置出主要的 SQL 執行指令。
	query = m.tableBuilder()
	// 保存最後一次所執行的 SQL 指令。
	m.LastQuery = query
	// 清除資料、欄位來重新開始一個資料表格遷移系統。
	m.clean()
	return
}

// Drop 會移除指定的資料表格，並在找不到資料表格時回傳錯誤。
func (m Migration) Drop(tableNames ...string) (query string) {
	query = m.dropBuilder(false, tableNames...)
	return
}

// DropIfExists 會移除指定的資料表格，若指定的資料表格存在的話。
func (m Migration) DropIfExists(tableNames ...string) (query string) {
	query = m.dropBuilder(true, tableNames...)
	return
}

// setColumnType 會替最後一個欄位設置其資料型態與長度。
func (m Migration) setColumnType(dataType string, arg ...interface{}) Migration {
	m.columns[len(m.columns)-1].dataType = dataType
	// 如果有指定長度的話就將其保存。
	if len(arg) == 1 {
		m.columns[len(m.columns)-1].length = arg[0]
	}
	return m
}

// tableBuilder 會建置出主要的資料表格 SQL 執行指令。
func (m Migration) tableBuilder() (query string) {
	var contentQuery string
	// 建立出不同的 SQL 執行指令。
	columnQuery := m.columnBuilder()
	foreignQuery := m.indexBuilder("FOREIGN KEY")
	primaryQuery := m.indexBuilder("PRIMARY KEY")
	uniqueQuery := m.indexBuilder("UNIQUE KEY")
	indexQuery := m.indexBuilder("INDEX")
	// 主要的開頭 SQL 執行指令。
	query = fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` ", m.table.name)

	// 湊合欄位和主鍵、索引的 SQL 執行指令。
	if columnQuery != "" {
		contentQuery += fmt.Sprintf("%s, ", columnQuery)
	}
	if foreignQuery != "" {
		contentQuery += fmt.Sprintf("%s, ", foreignQuery)
	}
	if primaryQuery != "" {
		contentQuery += fmt.Sprintf("%s, ", primaryQuery)
	}
	if uniqueQuery != "" {
		contentQuery += fmt.Sprintf("%s, ", uniqueQuery)
	}
	if indexQuery != "" {
		contentQuery += fmt.Sprintf("%s, ", indexQuery)
	}
	if contentQuery != "" {
		query += fmt.Sprintf("(%s) ", trim(contentQuery))
	}

	// 引擎種類。
	engineType := m.table.engineType
	if engineType == "" {
		engineType = EngineInnoDB
	}
	query += fmt.Sprintf("ENGINE=%s, ", strings.ToUpper(string(engineType)))
	// 字符集。
	if m.table.charset != "" {
		query += fmt.Sprintf("DEFAULT CHARSET=%s, ", m.table.charset)
	}
	// 校對字符格式。
	if m.table.collation != "" {
		query += fmt.Sprintf("COLLATE=%s, ", m.table.collation)
	}
	// 備註。
	if m.table.comment != "" {
		query += fmt.Sprintf("COMMENT='%s', ", m.table.comment)
	}
	// 移除結尾多餘的逗號與空白。
	query = trim(query)
	return
}

// clean 會清空上個遷移資料。
func (m Migration) clean() {
	m.table.comment = ""
	m.table.engineType = ""
	m.table.foreignKeys = []key{}
	m.table.indexKeys = []key{}
	m.table.primaryKeys = []key{}
	m.table.uniqueKeys = []key{}
	m.table.name = ""
	m.columns = []column{}
}

// indexBuilder 會建置移除資料表格的 SQL 執行指令。。
func (m Migration) dropBuilder(check bool, tableNames ...string) (query string) {
	var tables string
	// 遍歷資料表名稱切片來移除指定的資料表格。
	for _, name := range tableNames {
		tables = fmt.Sprintf("`%s`, ", name)
	}
	tables = strings.TrimRight(tables, ", ")

	// 建立 SQL 執行指令來準備移除指定資料表格。
	query = fmt.Sprintf("DROP TABLE %s", tables)
	if check {
		query = fmt.Sprintf("DROP TABLE IF EXISTS %s", tables)
	}
	// 保存最後一次執行的 SQL 指令。
	m.LastQuery = query
	// 清除資料、欄位來重新開始一個資料表格遷移系統。
	m.clean()
	return
}

// indexBuilder 會替資料表格的索引建置相關的 SQL 執行指令。
func (m Migration) indexBuilder(indexName string) (query string) {
	var keys []key
	var targetTable, targetColumns string
	// 透過指定的型態決定要從哪裡取得索引滋藥。
	switch indexName {
	case "PRIMARY KEY":
		keys = m.table.primaryKeys
	case "UNIQUE KEY":
		keys = m.table.uniqueKeys
	case "INDEX":
		keys = m.table.indexKeys
	case "FOREIGN KEY":
		keys = m.table.foreignKeys
	}

	// 每個索引群組。
	for _, v := range keys {
		// 建立欄位執行指令（`column_1`, `column_2`）。
		columns := fmt.Sprintf("`%s`", strings.Join(v.columns, "`,`"))
		// 替外鍵的目標欄位建立相關的 SQL 執行指令。
		if len(v.targetColumns) != 0 {
			targetColumns = ""
			for _, c := range v.targetColumns {
				// 從字串的 `目標表格.目標欄位` 格式中取得表格和欄位的名稱。
				splitedStr := strings.Split(c, ".")
				targetTable = splitedStr[0]
				// 將取得到的欄位名稱追加到 SQL 執行指令中。
				targetColumns += fmt.Sprintf("`%s`, ", splitedStr[1])
			}
			// 移除結尾多餘的逗號與空白。
			targetColumns = trim(targetColumns)
		}
		// 外鍵更新、移除時的相關動作執行指令。
		var onUpdate, onDelete string
		if v.onUpdate != "" {
			onUpdate = fmt.Sprintf(" ON UPDATE %s", v.onUpdate)
		}
		if v.onDelete != "" {
			onDelete = fmt.Sprintf(" ON DELETE %s", v.onDelete)
		}
		// 沒有群組名稱的索引。
		if v.name == "" && len(v.targetColumns) == 0 {
			query += fmt.Sprintf("%s (%s), ", indexName, columns)
			// 命名索引。
		} else if v.name != "" && len(v.targetColumns) == 0 {
			query += fmt.Sprintf("%s `%s` (%s), ", indexName, v.name, columns)
			// 沒有群組名稱的外鍵。
		} else if v.name == "" && len(v.targetColumns) != 0 {
			query += fmt.Sprintf("%s (%s) REFERENCES `%s` (%s)%s%s, ", indexName, columns, targetTable, targetColumns, onUpdate, onDelete)
			// 多個外鍵。
		} else if v.name != "" && len(v.targetColumns) != 0 {
			query += fmt.Sprintf("%s %s (%s) REFERENCES `%s` (%s)%s%s, ", indexName, v.name, columns, targetTable, targetColumns, onUpdate, onDelete)
		}
	}
	// 移除結尾多餘的逗號與空白。
	query = trim(query)
	return
}

// columnBuilder 會替資料表格的欄位建置相關的 SQL 執行指令。
func (m Migration) columnBuilder() (query string) {
	for _, v := range m.columns {
		// 欄位名稱。
		query += fmt.Sprintf("`%s` ", v.name)

		// 資料型態。
		dataType := strings.ToUpper(v.dataType)
		switch t := v.length.(type) {
		// VARCHAR(30)
		case int:
			query += fmt.Sprintf("%s(%d) ", dataType, t)
		// ENUM(1, 2, 3, 4)
		case []int:
			// 將 `[]int`` 轉換成字串並且透過逗點分隔。
			query += fmt.Sprintf("%s(%s) ", dataType, strings.Trim(strings.Join(strings.Split(fmt.Sprint(t), " "), ", "), "[]"))
		// ENUM("A", "B", "C")
		case []string:
			query += fmt.Sprintf("%s(%s) ", dataType, strings.Join(t, ", "))
		// FLOAT(1, 2) or ENUM(1, 2, "A", "B")
		case []interface{}:
			// 將選項從長度資料中展開。
			options := ""
			for _, o := range t {
				switch tt := o.(type) {
				case int:
					options += fmt.Sprintf("%d, ", tt)
				case string:
					options += fmt.Sprintf("'%s', ", tt)
				}
			}
			// 移除結尾多餘的逗號與空白。
			query += fmt.Sprintf("%s(%s) ", dataType, trim(options))
		// DATETIME
		case nil:
			query += fmt.Sprintf("%s ", dataType)
		}

		// 非負數。
		if v.unsigned {
			query += "UNSIGNED "
		}
		// 允許空值。
		if !v.nullable {
			query += "NOT NULL "
		}
		// 自動遞增。
		if v.autoIncrement {
			query += "AUTO_INCREMENT "
		}
		// 預設值。
		switch t := v.defaultValue.(type) {
		case int:
			query += fmt.Sprintf("DEFAULT %d ", t)
		case nil:
			query += fmt.Sprintf("DEFAULT NULL ")
		case string:
			if t == "CURRENT_TIMESTAMP" || t == "NOW()" || strings.Contains(t, "ON UPDATE ") {
				query += fmt.Sprintf("DEFAULT %s ", t)
			} else {
				query += fmt.Sprintf("DEFAULT '%s' ", t)
			}
		}
		// 主鍵。
		if v.primary {
			query += "PRIMARY KEY "
		}
		if v.unique {
			query += "UNIQUE "
		}
		if v.foreign != "" {
			m.table.foreignKeys = append(m.table.foreignKeys, key{
				columns:       []string{v.name},
				targetColumns: []string{v.foreign},
			})
		}
		// 備註。
		if v.comment != "" {
			query += fmt.Sprintf("COMMENT '%s'", v.comment)
		}
		// 最終結尾。
		query += ", "
	}
	// 移除結尾多餘的逗號與空白。
	query = trim(query)
	return
}
