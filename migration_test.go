package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var migration *Migration

func TestMigrationMain(t *testing.T) {
	migration = db.Migration()
}

func TestMigrationBasic(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Varchar(32).Primary().Create("test_table")
	assert.Equal("CREATE TABLE `test_table` (`test` VARCHAR(32) NOT NULL PRIMARY KEY) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationDrop(t *testing.T) {
	assert := assert.New(t)
	migration.Drop("test_table")
	assert.Equal("DROP TABLE `test_table`", migration.LastQuery)
}

func TestMigrationInsert(t *testing.T) {

}

func TestMigrationDataTypes(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").TinyInt(1).
		Column("test2").SmallInt(1).
		Column("test3").MediumInt(1).
		Column("test4").Int(1).
		Column("test5").BigInt(1).
		Column("test6").Char(1).
		Column("test7").Varchar(1).
		Column("test8").Binary(1).
		Column("test9").VarBinary(1).
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
		Create("test_table1")
	assert.Equal("CREATE TABLE `test_table1` (`test` TINYINT(1) NOT NULL , `test2` SMALLINT(1) NOT NULL , `test3` MEDIUMINT(1) NOT NULL , `test4` INT(1) NOT NULL , `test5` BIGINT(1) NOT NULL , `test6` CHAR(1) NOT NULL , `test7` VARCHAR(1) NOT NULL , `test8` BINARY(1) NOT NULL , `test9` VARBINARY(1) NOT NULL , `test10` BIT(1) NOT NULL , `test11` TINYTEXT NOT NULL , `test12` TEXT NOT NULL , `test13` MEDIUMTEXT NOT NULL , `test14` LONGTEXT NOT NULL , `test15` TINYBLOB NOT NULL , `test16` BLOB NOT NULL , `test17` MEDIUMBLOB NOT NULL , `test18` LONGBLOB NOT NULL , `test19` DATE NOT NULL , `test20` DATETIME NOT NULL , `test21` TIME NOT NULL , `test22` TIMESTAMP NOT NULL , `test23` YEAR NOT NULL , `test24` DOUBLE(2, 1) NOT NULL , `test25` DECIMAL(2, 1) NOT NULL , `test26` FLOAT(2, 1) NOT NULL , `test27` FLOAT(1) NOT NULL , `test28` ENUM('1', '2', '3', 'A', 'B', 'C') NOT NULL , `test29` SET('1', '2', '3', 'A', 'B', 'C') NOT NULL) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationTableType(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Varchar(32).MyISAM().Create("test_myisam_table")
	assert.Equal("CREATE TABLE `test_myisam_table` (`test` VARCHAR(32) NOT NULL) ENGINE=MYISAM", migration.LastQuery)
	migration.Column("test").Varchar(32).InnoDB().Create("test_innodb_table")
	assert.Equal("CREATE TABLE `test_innodb_table` (`test` VARCHAR(32) NOT NULL) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationDefault(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Varchar(32).Nullable().Default(nil).Create("test_default_null_table")
	assert.Equal("CREATE TABLE `test_default_null_table` (`test` VARCHAR(32) DEFAULT NULL) ENGINE=INNODB", migration.LastQuery)
	migration.Column("test").Varchar(32).Default("string").Create("test_default_string_table")
	assert.Equal("CREATE TABLE `test_default_string_table` (`test` VARCHAR(32) NOT NULL DEFAULT 'string') ENGINE=INNODB", migration.LastQuery)
	migration.Column("test").Int(32).Default(12).Create("test_default_int_table")
	assert.Equal("CREATE TABLE `test_default_int_table` (`test` INT(32) NOT NULL DEFAULT 12) ENGINE=INNODB", migration.LastQuery)
	migration.Column("test").Timestamp().Default("CURRENT_TIMESTAMP").Create("test_default_timestamp_table")
	assert.Equal("CREATE TABLE `test_default_timestamp_table` (`test` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=INNODB", migration.LastQuery)
	migration.Column("test").Timestamp().Default("CURRENT_TIMESTAMP").Create("test_default_timestamp_table")
	assert.Equal("CREATE TABLE `test_default_timestamp_table` (`test` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP) ENGINE=INNODB", migration.LastQuery)
	migration.Column("test").DateTime().Default("NOW()").Create("test_default_datetime_table")
	assert.Equal("CREATE TABLE `test_default_datetime_table` (`test` DATETIME NOT NULL DEFAULT NOW()) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationNullable(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Varchar(32).Nullable().Create("test_nullable_table")
	assert.Equal("CREATE TABLE `test_nullable_table` (`test` VARCHAR(32) DEFAULT NULL) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationUnsigned(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Int(10).Unsigned().Create("test_unsigned_table")
	assert.Equal("CREATE TABLE `test_unsigned_table` (`test` INT(10) UNSIGNED NOT NULL) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationAutoIncrement(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Int(10).AutoIncrement().Create("test_auto_increment_table")
	assert.Equal("CREATE TABLE `test_auto_increment_table` (`test` INT(10) NOT NULL AUTO_INCREMENT) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationComment(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Int(10).Comment("月月，搭拉安！").Create("test_column_comment_table")
	assert.Equal("CREATE TABLE `test_column_comment_table` (`test` INT(10) NOT NULL COMMENT '月月，搭拉安！') ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationTableComment(t *testing.T) {
	assert := assert.New(t)
	migration.Column("test").Int(10).Create("test_comment_table", "月月，搭拉安！")
	assert.Equal("CREATE TABLE `test_comment_table` (`test` INT(10) NOT NULL) ENGINE=INNODB, COMMENT='月月，搭拉安！'", migration.LastQuery)
}

func TestMigrationPrimaryKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).Primary().
		Column("test2").Varchar(32).
		Create("test_table2")
	assert.Equal("CREATE TABLE `test_table2` (`test` VARCHAR(32) NOT NULL PRIMARY KEY , `test2` VARCHAR(32) NOT NULL) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationNamingPrimaryKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Primary("pk_test", []string{"test", "test2"}).
		Create("test_table3")
	assert.Equal("CREATE TABLE `test_table3` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, PRIMARY KEY `pk_test` (`test`,`test2`)) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationMultiPrimaryKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Primary([]string{"test", "test2"}).
		Create("test_table4")
	assert.Equal("CREATE TABLE `test_table4` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, PRIMARY KEY (`test`,`test2`)) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationUniqueKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).Unique().
		Column("test2").Varchar(32).
		Create("test_table5")
	assert.Equal("CREATE TABLE `test_table5` (`test` VARCHAR(32) NOT NULL UNIQUE , `test2` VARCHAR(32) NOT NULL) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationNamingUniqueKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Unique("uk_test", []string{"test", "test2"}).
		Unique("uk_test2", []string{"test3", "test4"}).
		Create("test_table6")
	assert.Equal("CREATE TABLE `test_table6` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, UNIQUE KEY `uk_test` (`test`,`test2`), UNIQUE KEY `uk_test2` (`test3`,`test4`)) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationMultiUniqueKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Unique([]string{"test", "test2"}).
		Unique([]string{"test3", "test4"}).
		Create("test_table7")
	assert.Equal("CREATE TABLE `test_table7` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, UNIQUE KEY (`test`,`test2`), UNIQUE KEY (`test3`,`test4`)) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationAnonymousIndexKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).Index().
		Column("test2").Varchar(32).Index().
		Create("test_table14")
	assert.Equal("CREATE TABLE `test_table14` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL, INDEX `test` (`test`), INDEX `test2` (`test2`)) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationNamingIndexKey(t *testing.T) {
	assert := assert.New(t)
	migration.
		Column("test").Varchar(32).
		Column("test2").Varchar(32).
		Column("test3").Varchar(32).
		Column("test4").Varchar(32).
		Index("ik_test", []string{"test", "test2"}).
		Index("ik_test2", []string{"test3", "test4"}).
		Create("test_table9")
	assert.Equal("CREATE TABLE `test_table9` (`test` VARCHAR(32) NOT NULL , `test2` VARCHAR(32) NOT NULL , `test3` VARCHAR(32) NOT NULL , `test4` VARCHAR(32) NOT NULL, INDEX `ik_test` (`test`,`test2`), INDEX `ik_test2` (`test3`,`test4`)) ENGINE=INNODB", migration.LastQuery)
}

func TestMigrationMixedKeys(t *testing.T) {

}

func TestMigrationForeignKey(t *testing.T) {

}

func TestMigrationMultipleForeignKey(t *testing.T) {

}

func TestMigrationNamingForeignKey(t *testing.T) {

}