package rushia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var migration Migration

func TestMigrationMain(t *testing.T) {
	migration = NewMigration()
}

func TestMigrationBasic(t *testing.T) {
	assert := assert.New(t)
	query := migration.Table("test_table").Column("test").Varchar(32).Primary().Create()
	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table` (`test` VARCHAR(32) NOT NULL PRIMARY KEY) ENGINE=INNODB", query)
}

func TestMigrationDrop(t *testing.T) {
	assert := assert.New(t)
	query := migration.Drop("test_table")
	assert.Equal("DROP TABLE `test_table`", query)
}

func TestMigrationDropIfExists(t *testing.T) {
	assert := assert.New(t)
	query := migration.DropIfExists("test_table")

	assert.Equal("DROP TABLE IF EXISTS `test_table`", query)
}

func TestMigrationDataTypes(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table1").
		Column("test").TinyInt(1).
		Column("test2").SmallInt(1).
		Column("test3").MediumInt(1).
		Column("test4").Int(1).
		Column("test5").BigInt(1).
		Column("test6").Char(1).
		Column("test7").Varchar(1).
		Column("test8").Binary(1).
		Column("test9").Varbinary(1).
		Column("test10").Bit(1).
		Column("test11").TinyText().
		Column("test12").Text().
		Column("test13").MediumText().
		Column("test14").LongText().
		Column("test15").TinyBlob().
		Column("test16").Blob().
		Column("test17").MediumBlob().
		Column("test18").LongBlob().
		Column("test19").Date().
		Column("test20").DateTime().
		Column("test21").Time().
		Column("test22").Timestamp().
		Column("test23").Year().
		Column("test24").Double(2, 1).
		Column("test25").Decimal(2, 1).
		Column("test26").Float(2, 1).
		Column("test27").Float(1).
		Column("test28").Enum("1", "2", "3", "A", "B", "C").
		Column("test29").Set("1", "2", "3", "A", "B", "C").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table1` (`test` TINYINT(1) NOT NULL , `test2` SMALLINT(1) NOT NULL , `test3` MEDIUMINT(1) NOT NULL , `test4` INT(1) NOT NULL , `test5` BIGINT(1) NOT NULL , `test6` CHAR(1) NOT NULL , `test7` VARCHAR(1) NOT NULL , `test8` BINARY(1) NOT NULL , `test9` VARBINARY(1) NOT NULL , `test10` BIT(1) NOT NULL , `test11` TINYTEXT NOT NULL , `test12` TEXT NOT NULL , `test13` MEDIUMTEXT NOT NULL , `test14` LONGTEXT NOT NULL , `test15` TINYBLOB NOT NULL , `test16` BLOB NOT NULL , `test17` MEDIUMBLOB NOT NULL , `test18` LONGBLOB NOT NULL , `test19` DATE NOT NULL , `test20` DATETIME NOT NULL , `test21` TIME NOT NULL , `test22` TIMESTAMP NOT NULL , `test23` YEAR NOT NULL , `test24` DOUBLE(2, 1) NOT NULL , `test25` DECIMAL(2, 1) NOT NULL , `test26` FLOAT(2, 1) NOT NULL , `test27` FLOAT(1) NOT NULL , `test28` ENUM('1', '2', '3', 'A', 'B', 'C') NOT NULL , `test29` SET('1', '2', '3', 'A', 'B', 'C') NOT NULL) ENGINE=INNODB", query)
}

func TestMigrationTableType(t *testing.T) {
	assert := assert.New(t)
	query := migration.Table("test_myisam_table").Column("test").Varchar(32).Engine(EngineMyISAM).Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_myisam_table` (`test` VARCHAR(32) NOT NULL) ENGINE=MYISAM", query)
	query = migration.Table("test_innodb_table").Column("test").Varchar(32).Engine(EngineInnoDB).Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_innodb_table` (`test` VARCHAR(32) NOT NULL) ENGINE=INNODB", query)
}

func TestMigrationDefault(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_default_null_table").
		Column("test").
		Varchar(32).
		Nullable().
		Default(nil).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_default_null_table` (`test` VARCHAR(32) DEFAULT NULL) ENGINE=INNODB", query)
	query = migration.
		Table("test_default_string_table").
		Column("test").
		Varchar(32).
		Default("string").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_default_string_table` (`test` VARCHAR(32) NOT NULL DEFAULT 'string') ENGINE=INNODB", query)
	query = migration.
		Table("test_default_int_table").
		Column("test").
		Int(32).
		Default(12).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_default_int_table` (`test` INT(32) NOT NULL DEFAULT 12) ENGINE=INNODB", query)
	query = migration.
		Table("test_default_timestamp_table").
		Column("test").
		Timestamp().
		Default("CURRENT_TIMESTAMP").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_default_timestamp_table` (`test` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=INNODB", query)
	query = migration.
		Table("test_on_update_timestamp_table").
		Column("test").
		Timestamp().
		Default("CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_on_update_timestamp_table` (`test` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP) ENGINE=INNODB", query)
	query = migration.
		Table("test_default_timestamp_table").
		Column("test").
		Timestamp().
		Default("CURRENT_TIMESTAMP").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_default_timestamp_table` (`test` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=INNODB", query)
	query = migration.
		Table("test_default_datetime_table").
		Column("test").
		DateTime().
		Default("NOW()").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_default_datetime_table` (`test` DATETIME NOT NULL DEFAULT NOW()) ENGINE=INNODB", query)
}

func TestMigrationNullable(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_nullable_table").
		Column("test").
		Varchar(32).
		Nullable().
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_nullable_table` (`test` VARCHAR(32) DEFAULT NULL) ENGINE=INNODB", query)
}

func TestMigrationUnsigned(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_unsigned_table").
		Column("test").
		Int(10).
		Unsigned().
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_unsigned_table` (`test` INT(10) UNSIGNED NOT NULL) ENGINE=INNODB", query)
}

func TestMigrationAutoIncrement(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_auto_increment_table").
		Column("test").
		Int(10).
		AutoIncrement().
		Primary().
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_auto_increment_table` (`test` INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY) ENGINE=INNODB", query)
}

func TestMigrationComment(t *testing.T) {
	assert := assert.New(t)
	query := migration.Table("test_column_comment_table").Column("test").Int(10).Comment("月月，搭拉安！").Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_column_comment_table` (`test` INT(10) NOT NULL COMMENT '月月，搭拉安！') ENGINE=INNODB", query)
}

func TestMigrationTableComment(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_comment_table", "月月，搭拉安！").
		Column("test").
		Int(10).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_comment_table` (`test` INT(10) NOT NULL) ENGINE=INNODB, COMMENT='月月，搭拉安！'", query)
}

func TestMigrationPrimaryKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table2").
		Column("test").Varchar(32).Primary().
		Column("test2").Varchar(32).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table2` (`test` VARCHAR(32) NOT NULL PRIMARY KEY , `test2` VARCHAR(32) NOT NULL) ENGINE=INNODB", query)
}

func TestMigrationNamingPrimaryKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table3").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Primary("pk_test", []string{"test", "test2"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table3` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, PRIMARY KEY `pk_test` (`test`,`test2`)) ENGINE=INNODB", query)
}

func TestMigrationMultiPrimaryKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table4").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Primary([]string{"test", "test2"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table4` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, PRIMARY KEY (`test`,`test2`)) ENGINE=INNODB", query)
}

func TestMigrationUniqueKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table5").
		Column("test").Varchar(32).Unique().
		Column("test2").Varchar(32).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table5` (`test` VARCHAR(32) NOT NULL UNIQUE , `test2` VARCHAR(32) NOT NULL) ENGINE=INNODB", query)
}

func TestMigrationNamingUniqueKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table6").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Unique("uk_test", []string{"test", "test2"}).
		Unique("uk_test2", []string{"test3", "test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table6` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, UNIQUE KEY `uk_test` (`test`,`test2`), UNIQUE KEY `uk_test2` (`test3`,`test4`)) ENGINE=INNODB", query)
}

func TestMigrationMultiUniqueKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table7").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Unique([]string{"test", "test2"}).
		Unique([]string{"test3", "test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table7` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, UNIQUE KEY (`test`,`test2`), UNIQUE KEY (`test3`,`test4`)) ENGINE=INNODB", query)
}

func TestMigrationAnonymousIndexKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table14").
		Column("test").Varchar(32).Index().
		Column("test2").Varchar(32).Index().
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table14` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, INDEX (`test`), INDEX (`test2`)) ENGINE=INNODB", query)
}

func TestMigrationNamingIndexKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table9").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Index("ik_test", []string{"test", "test2"}).
		Index("ik_test2", []string{"test3", "test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table9` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, INDEX `ik_test` (`test`,`test2`), INDEX `ik_test2` (`test3`,`test4`)) ENGINE=INNODB", query)
}

func TestMigrationMixedKeys(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table10").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Primary([]string{"test", "test2"}).
		Unique([]string{"test3", "test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table10` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, PRIMARY KEY (`test`,`test2`), UNIQUE KEY (`test3`,`test4`)) ENGINE=INNODB", query)
	query = migration.
		Table("test_table11").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Index("ik_test", []string{"test", "test2"}).
		Unique([]string{"test3", "test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table11` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, UNIQUE KEY (`test3`,`test4`), INDEX `ik_test` (`test`,`test2`)) ENGINE=INNODB", query)
	query = migration.
		Table("test_table12").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Primary([]string{"test", "test2"}).
		Index("ik_test", []string{"test3", "test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table12` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, PRIMARY KEY (`test`,`test2`), INDEX `ik_test` (`test3`,`test4`)) ENGINE=INNODB", query)
	query = migration.
		Table("test_table13").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Column("test5").Varchar(32).
		Column("test6").Varchar(32).
		Index("ik_test", []string{"test", "test2"}).
		Unique([]string{"test3", "test4"}).
		Primary([]string{"test5", "test6"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table13` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL , `test5` VARCHAR(32) NOT NULL , `test6` VARCHAR(32) NOT NULL, PRIMARY KEY (`test5`,`test6`), UNIQUE KEY (`test3`,`test4`), INDEX `ik_test` (`test`,`test2`)) ENGINE=INNODB", query)
}

func TestMigrationForeignKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table15").
		Column("test").Varchar(32).Foreign("test_table13.test5").
		Column("test2").Varchar(32).Foreign("test_table13.test6").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table15` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, FOREIGN KEY (`test`,`test2`) REFERENCES `test_table13` (`test5`, `test6`)) ENGINE=INNODB", query)
}

func TestMigrationMultipleForeignKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table16").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Foreign([]string{"test"}, []string{"test_table13.test5"}).
		Foreign([]string{"test2"}, []string{"test_table12.test3"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table16` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, FOREIGN KEY (`test`) REFERENCES `test_table13` (`test5`), FOREIGN KEY (`test2`) REFERENCES `test_table12` (`test3`)) ENGINE=INNODB", query)
}

func TestMigrationNamingForeignKey(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table17").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Foreign("fk_test", []string{"test", "test2"}, []string{"test_table13.test5", "test_table13.test6"}).
		Foreign("fk_test2", []string{"test", "test2"}, []string{"test_table12.test3", "test_table12.test4"}).
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table17` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, FOREIGN KEY fk_test (`test`,`test2`) REFERENCES `test_table13` (`test5`, `test6`), FOREIGN KEY fk_test2 (`test`,`test2`) REFERENCES `test_table12` (`test3`, `test4`)) ENGINE=INNODB", query)
}

func TestMigrationForeignKeysOnActions(t *testing.T) {
	assert := assert.New(t)
	query := migration.
		Table("test_table18").
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Foreign([]string{"test"}, []string{"test_table13.test5"}).OnUpdate("NO ACTION").OnDelete("NO ACTION").
		Foreign([]string{"test2"}, []string{"test_table12.test3"}).OnUpdate("CASCADE").OnDelete("RESTRICT").
		Create()

	assert.Equal("CREATE TABLE IF NOT EXISTS `test_table18` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, FOREIGN KEY (`test`) REFERENCES `test_table13` (`test5`) ON UPDATE NO ACTION ON DELETE NO ACTION, FOREIGN KEY (`test2`) REFERENCES `test_table12` (`test3`) ON UPDATE CASCADE ON DELETE RESTRICT) ENGINE=INNODB", query)
}
