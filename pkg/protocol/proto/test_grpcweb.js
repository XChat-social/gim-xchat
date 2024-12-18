const { BusinessExtClient } = require('./business.ext_grpc_web_pb'); // 服务类
const { Empty } = require('google-protobuf/google/protobuf/empty_pb'); // 使用 google.protobuf.Empty

// const client = new BusinessExtClient('http://13.61.35.52:8020', null, null);
const client = new BusinessExtClient('http://localhost:8081', null, null);

// 创建空请求对象
const request = new Empty();

// 调用服务方法
client.getTwitterAuthorizeURL(request, {}, (err, response) => {
    if (err) {
        console.error('Error:', err.message);
    } else {
        console.log('Response:', response.toObject());
    }
});