package sqlite

import "gorm.io/gorm"

type DataObject interface {
    TableName() string
    Get() interface{}
}

type ClientSQLite interface {
    CloseSession()
    
    // 关闭连接
    Close()
    // 初始化建表操作
    InitTable(dos ...DataObject) error
    // 添加数据
    CreateObject(do DataObject) error
    // 判断数据是否存在
    ExistObject(do DataObject) (bool, error)
    // 删除数据
    DeleteObject(do DataObject) error
    // 更新数据
    UpdateObject(do DataObject) error
    // 关联更新
    AssociationsUpdateObject(do DataObject) error
    // 关联删除
    AssociationsDeleteObject(do DataObject) error
    // 查询单个数据
    GetObject(doCond DataObject, do DataObject) error
    // 预加载查询单个数据
    GetPreloadObject(doCond DataObject, do DataObject) error
    // 查询列表数据
    GetObjects(doCond DataObject, do DataObject, likeParam *LikeQueryParam,
        order *OrderQueryParam, offset, limit int) ([]interface{}, int64, error)
    // 批量更新数据
    BatchUpdate(fields, conflictFields, updateFields []string, tableName string, records [][]interface{}) error
    
    //事物相关
    ExecSqlWithTransaction(funcs ...func(db *gorm.DB) error) error
}
