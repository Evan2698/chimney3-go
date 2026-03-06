
#!/bin/bash

# 进程名
PROCESS_NAME="chimney"
PROCESS_FULL_PATH="/home/evan/chimney3-go/chimney"

# 检测间隔（秒）
INTERVAL=20

while true; do
    # 使用ps -efww检测进程
    if ! ps -efww | grep -v grep | grep "$PROCESS_FULL_PATH" > /dev/null; then
        # 如果进程不存在，则后台启动
        nohup $PROCESS_FULL_PATH > /dev/null 2>&1 &
        echo "$(date): 进程 $PROCESS_FULL_PATH 未运行，已启动"
    else
        # 如果进程存在，什么都不做
        echo "$(date): 进程 $PROCESS_FULL_PATH 正在运行"
    fi

    # 等待指定时间
    sleep $INTERVAL
done
