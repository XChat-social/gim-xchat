import * as grpcWeb from 'grpc-web';
import { BusinessExtClient } from './business.ext_grpc_web_pb';
import { GetTaskStatusReq } from './business.ext_pb';

// 创建 gRPC 客户端实例
const client = new BusinessExtClient('https://api.xchat.social/business-web', null, null);

// 创建 GetTaskStatus 请求对象
const taskId = 1001;
const request = new GetTaskStatusReq();
request.setTaskId(taskId);

// 创建 Metadata 对象
const metadata = {
    'device-id': 0,
    'token': '6e060626fcab9faca0f60423ebaff91d',
    'user-id': 2
};

// 调用服务方法
client.getTaskStatus(request, metadata, (err, response) => {
    if (err) {
        console.error('Error:', err.message);
    } else {
        console.log('Response:', response.toObject());
        const status = response.getStatus();
        const message = response.getMessage();
        console.log(`Task ID: ${taskId}, Status: ${status}, Message: ${message}`);
    }
});