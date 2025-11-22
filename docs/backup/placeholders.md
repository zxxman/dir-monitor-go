# 占位符清单

以下占位符可在 `command` 字符串中使用，由系统在执行前进行替换：

- `${FILE_PATH}`: 触发事件的完整文件路径（绝对路径）
- `${FILE_NAME}`: 文件名（含扩展名）
- `${FILE_DIR}`: 所在目录路径（绝对路径）
- `${FILE_EXT}`: 文件扩展名（含点）
- `${EVENT_TIME}`: 触发事件的时间戳（RFC3339格式）
- `${EVENT_TYPE}`: 事件类型（created/modified/deleted/renamed）

此外，还可以使用系统环境变量作为占位符，例如 `${HOME}`、`${USER}` 等。

示例：
```json
{
  "command": "python process.py ${FILE_PATH} --ts ${EVENT_TIME} --user ${USER}"
}
```

**注意**：占位符必须使用 `${VARIABLE}` 格式，不支持其他格式。
