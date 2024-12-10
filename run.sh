#!/bin/bash

# 定义一个函数，用于打包、停止旧进程并启动新进程
restart_service() {
    local service_dir=$1
    local service_name=$2
    local log_file=$3

    echo "进入 ${service_dir} 目录"
    cd "${service_dir}" || exit 1

    echo "删除旧的 ${service_name} 可执行文件"
    rm -f "${service_name}"

    echo "开始打包 ${service_name}"
    go build -o "${service_name}" main.go
    if [ $? -ne 0 ]; then
        echo "打包 ${service_name} 失败！"
        exit 1
    fi
    echo "打包 ${service_name} 成功"

    echo "停止旧的 ${service_name} 服务进程"
    pkill -9 "${service_name}" 2>/dev/null || true

    echo "启动新的 ${service_name} 服务"
    nohup "./${service_name}" > "../../${log_file}" 2>&1 &
    echo "${service_name} 服务已启动，日志输出到 ${log_file}"

    echo "返回上级目录"
    cd - || exit 1
}

# 逐个重启服务
restart_service "cmd/business" "business" "business.log"
restart_service "cmd/logic" "logic" "logic.log"
restart_service "cmd/connect" "connect" "connect.log"
restart_service "cmd/api" "api" "api.log"

# 如果需要启动 file 服务，取消注释以下代码
# restart_service "cmd/file" "file" "file.log"

echo "所有服务已成功启动！后台运行中。"