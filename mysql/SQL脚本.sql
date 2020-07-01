/*
 SQL server 获取表字段信息
 */
SELECT ColumnsName = c.name, Description = ex.value,
        t.name+'('+cast(c.max_length as varchar) +')'
FROM sys.columns c
LEFT OUTER JOIN sys.extended_properties ex
    ON
                ex.major_id = c.object_id
            AND ex.minor_id = c.column_id
            AND ex.name = 'MS_Description'
        left outer join
    systypes t
    on c.system_type_id=t.xtype
WHERE OBJECTPROPERTY(c.object_id, 'IsMsShipped')=0
  AND t.name<>'sysname'
  AND OBJECT_NAME(c.object_id)
    IN ('t_MoveTypeSub' )
ORDER BY OBJECT_NAME(c.object_id),c.column_id ASC
;
/*
MySQL 获取表字段信息
 */
select COLUMN_NAME,COLUMN_TYPE,COLUMN_COMMENT
from information_schema.COLUMNS
where table_name = 'movetypesub'
  and table_schema = 'kdsx_a1_cleandev';