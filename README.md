在上一节我们完成了 select 语句的解释执行，本节我们看看 Update 和 Delete 对应的语句如何解释执行，当然他们的实现原理跟我们前面实现的 select 语句执行大同小异。无论是 update还是 delete 都是对数据表的修改，因此他们的实现方法基本相同。

假设我们要执行如下 sql 语句：

update STUDENT set MajorId=20 where MajorId=30 and GradYear=2020

delete from STUDENT where MajorId=30 and GradYear=2020
要完成上面的代码，我们需要 scan底层的文件块，找到所有满足 where 条件的记录，如果语句是 update，那么把找到的记录修改掉，如果是 delete，那么把找到的记录给删除掉。我们看看具体的代码实现，首先我们添加 UpdatePlanner 的接口定义，在 planner 的 interface.go 文件中增加代码如下：

package planner

import (
    "parser"
    "record_manager"
    "tx"
)

type Plan interface {
    Open() interface{}
    BlocksAccessed() int               //对应 B(s)
    RecordsOutput() int                //对应 R(s)
    DistinctValues(fldName string) int //对应 V(s,F)
    Schema() record_manager.SchemaInterface
}

type QueryPlanner interface {
    CreatePlan(data *parser.QueryData, tx *tx.Transation) Plan
}

type UpdatePlanner interface {
    /*
        解释执行 insert 语句，返回被修改的记录条数
    */
    ExecuteInsert(data *parser.InsertData, tx *tx.Transation) int

    /*
        解释执行 delete 语句，返回被删除的记录数
    */
    ExecuteDelete(data *parser.DeleteData, tx *tx.Transation) int

    /*
        解释执行 create table 语句，返回新建表中的记录数
    */
    ExecuteCreateTable(data *parser.CreateTableData, tx *tx.Transation) int

    /*
        解释执行 create index 语句，返回当前建立了索引的记录数
        TODO
    */
    //ExecuteCreateIndex(data *parser.CreateIndexData, tx *tx.Transation) int
}
在 planner 目录下新建文件 update_planner.go，输入代码如下：

package planner

import (
    metadata_manager "metadata_management"
    "parser"
    "query"
    "tx"
)

type BasicUpdatePlanner struct {
    mdm *metadata_manager.MetaDataManager
}

func NewBasicUpdatePlanner(mdm *metadata_manager.MetaDataManager) *BasicUpdatePlanner {
    return &BasicUpdatePlanner{
        mdm: mdm,
    }
}

func (b *BasicUpdatePlanner) ExecuteDelete(data *parser.DeleteData, tx *tx.Transation) int {
    /*
        先 scan 出要处理的记录，然后执行删除
    */
    tablePlan := NewTablePlan(tx, data.TableName(), b.mdm)
    selectPlan := NewSelectPlan(tablePlan, data.Pred())
    scan := selectPlan.Open()
    updateScan := scan.(*query.SelectionScan)
    count := 0
    for updateScan.Next() {
        updateScan.Delete()
        count += 1
    }
    updateScan.Close()
    return count
}

func (b *BasicUpdatePlanner) ExecuteModify(data *parser.ModifyData, tx *tx.Transation) int {
    /*
        先 scan 出选中的记录，然后修改记录的信息
    */
    tablePlan := NewTablePlan(tx, data.TableName(), b.mdm)
    selectPlan := NewSelectPlan(tablePlan, data.Pred())
    scan := selectPlan.Open()
    updateScan := scan.(*query.SelectionScan)
    count := 0
    for updateScan.Next() {
        val := data.NewValue().Evaluate(scan.(query.Scan))
        updateScan.SetVal(data.TargetField(), val)
        count += 1
    }
    return count
}

func (b *BasicUpdatePlanner) ExecuteInsert(data *parser.InsertData, tx *tx.Transation) int {
    tablePlan := NewTablePlan(tx, data.TableName(), b.mdm)
    updateScan := tablePlan.Open().(*query.TableScan)
    updateScan.Insert()
    insertFields := data.Fields()
    insertedVals := data.Vals()

    for i := 0; i < len(insertFields); i++ {
        updateScan.SetVal(insertFields[i], insertedVals[i])
    }

    updateScan.Close()
    return 1
}

func (b *BasicUpdatePlanner) ExecuteCreateTable(data *parser.CreateTableData, tx *tx.Transation) int {
    b.mdm.CreateTable(data.TableName(), data.NewSchema(), tx)
    return 0
}

func (b *BasicUpdatePlanner) ExecuteView(data *parser.ViewData, tx *tx.Transation) int {
    b.mdm.CreateView(data.ViewName(), data.ViewDef(), tx)
    return 0
}

func (b *BasicUpdatePlanner) ExecuteIndex(data *parser.IndexData, tx *tx.Transation) int {
    //b.mdm.CreateIndex
    //TODO

    return 0
}
在上面代码中BasicUpdatePlanner用于负责实现与数据库表内容修改相关的操作，例如插入，修改和删除，它导出的接口 ExecuteDelete, ExecuteModify, ExecuteInsert 分别负责表的删除，修改和插入，删除和修改的逻辑类似，首先都是通过 TablePlan 和 SelectPlan 找出要修改的记录，然后进行相应的操作，在上面代码实现中我们留有与索引相关的操作没有实现，因为索引是我们后续章节的一个重要内容。

完成了上面代码后，我们看看如何调用他们然后检验一下实现效果，在 main.go 中增加代码如下：

package main

import (
    bmg "buffer_manager"
    fm "file_manager"
    "fmt"
    lm "log_manager"
    metadata_manager "metadata_management"
    "parser"
    "planner"
    "query"
    "tx"
)

func PrintStudentTable(tx *tx.Transation, mdm *metadata_manager.MetaDataManager) {
    queryStr := "select name, majorid, gradyear from STUDENT"
    p := parser.NewSQLParser(queryStr)
    queryData := p.Query()
    test_planner := planner.CreateBasicQueryPlanner(mdm)
    test_plan := test_planner.CreatePlan(queryData, tx)
    test_interface := (test_plan.Open())
    test_scan, _ := test_interface.(query.Scan)
    for test_scan.Next() {
        fmt.Printf("name: %s, majorid: %d, gradyear: %d\n",
            test_scan.GetString("name"), test_scan.GetInt("majorid"),
            test_scan.GetInt("gradyear"))
    }
}

func CreateInsertUpdateByUpdatePlanner() {
    file_manager, _ := fm.NewFileManager("student", 2048)
    log_manager, _ := lm.NewLogManager(file_manager, "logfile.log")
    buffer_manager := bmg.NewBufferManager(file_manager, log_manager, 3)
    tx := tx.NewTransation(file_manager, log_manager, buffer_manager)
    mdm := metadata_manager.NewMetaDataManager(false, tx)

    updatePlanner := planner.NewBasicUpdatePlanner(mdm)
    createTableSql := "create table STUDENT (name varchar(16), majorid int, gradyear int)"
    p := parser.NewSQLParser(createTableSql)
    tableData := p.UpdateCmd().(*parser.CreateTableData)
    updatePlanner.ExecuteCreateTable(tableData, tx)

    insertSQL := "insert into STUDENT (name, majorid, gradyear) values(\"tylor\", 30, 2020)"
    p = parser.NewSQLParser(insertSQL)
    insertData := p.UpdateCmd().(*parser.InsertData)
    updatePlanner.ExecuteInsert(insertData, tx)
    insertSQL = "insert into STUDENT (name, majorid, gradyear) values(\"tom\", 35, 2023)"
    p = parser.NewSQLParser(insertSQL)
    insertData = p.UpdateCmd().(*parser.InsertData)
    updatePlanner.ExecuteInsert(insertData, tx)

    fmt.Println("table after insert:")
    PrintStudentTable(tx, mdm)

    updateSQL := "update STUDENT set majorid=20 where majorid=30 and gradyear=2020"
    p = parser.NewSQLParser(updateSQL)
    updateData := p.UpdateCmd().(*parser.ModifyData)
    updatePlanner.ExecuteModify(updateData, tx)

    fmt.Println("table after update:")
    PrintStudentTable(tx, mdm)

    deleteSQL := "delete from STUDENT where majorid=35"
    p = parser.NewSQLParser(deleteSQL)
    deleteData := p.UpdateCmd().(*parser.DeleteData)
    updatePlanner.ExecuteDelete(deleteData, tx)

    fmt.Println("table after delete")
    PrintStudentTable(tx, mdm)
}

func main() {
    CreateInsertUpdateByUpdatePlanner()
}
在上面代码实现中，CreateInsertUpdateByUpdatePlanner函数先创建BasicUpdatePlanner对象，然后调用其 ExecuteCreateTable接口创建 STUDENT 表，接着使用 sql 解释器解析 insert 语句后创建 InsertData 对象，然后调用ExecuteInsert接口将记录插入数据库表，接下来以同样的方式调用ExecuteModify， ExecuteDelete接口来实现对数据库表中有关记录的修改和删除，完成上面代码后 运行go run main.go,执行起来效果如下：

able after insert:
name: tylor, majorid: 30, gradyear: 2020
name: tom, majorid: 35, gradyear: 2023
table after update:
name: tylor, majorid: 20, gradyear: 2020
name: tom, majorid: 35, gradyear: 2023
table after delete
name: tylor, majorid: 20, gradyear: 2020
从运行结果可以看到，我们对数据库表的建立，插入，修改和删除等操作的基本结果是正确的。更多内容和调试演示视频请在 b 站搜索:Coding 迪斯尼。

